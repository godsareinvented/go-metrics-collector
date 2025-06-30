package metric

import (
	"context"
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
	"sync"
	"time"
)

type (
	interimData struct {
		err           error
		metric        generaldto.Metrics
		metricList    []generaldto.Metrics
		collectedData agentdto.CollectedMetricData
		strategyMap   map[string]interfaces.ParsingStrategy
	}

	MetricManager struct {
		metricsToCollect []string
		dataCollector    interfaces.MetricDataCollector
		client           interfaces.Client
		pool             *sync.Pool
	}
)

var (
	ErrNotInitialised   = errors.New("dependencies are not initialized")
	ErrInvalidArguments = errors.New("not initialized metric list")
)

// Pool For tests
func (metricManager *MetricManager) Pool() *sync.Pool {
	return metricManager.pool
}

func (metricManager *MetricManager) CollectAndSend(ctx context.Context, onDone func(err error)) {
	data := metricManager.pool.Get().(*interimData)
	if metricManager.metricsToCollect == nil || len(metricManager.metricsToCollect) == 0 || metricManager.dataCollector == nil || data.strategyMap == nil {
		onDone(ErrNotInitialised)
		return
	}

	currentCtx := context.Background()
	currentCtx, cancel := util.CombineContexts(ctx, currentCtx)

	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}

	go func(errCh chan error) {
		wg.Wait()
		close(errCh)

		var finalErr error
		for err := range errCh {
			finalErr = multierror.Append(finalErr, err)
		}

		onDone(finalErr)
	}(errCh)

	wg.Add(2)
	metricManager.collect(currentCtx, cancel, &wg, errCh)
	metricManager.send(currentCtx, &wg, errCh)
}

func (metricManager *MetricManager) collect(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, errCh chan<- error) {
	go func(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, errCh chan<- error) {
		defer cancel()
		defer wg.Done()

		var data *interimData
		var metricName string
		var metrics []generaldto.Metrics

		for {
			select {
			case <-ctx.Done():
				return
			default:
				data = metricManager.pool.Get().(*interimData)

				err := metricManager.dataCollector.Collect(&data.collectedData)
				if err != nil {
					errCh <- err
					return
				}

				metrics = make([]generaldto.Metrics, 0, len(metricManager.metricsToCollect))
				for _, metricName = range metricManager.metricsToCollect {
					data.strategyMap[metricName].GetMetric(&data.metric, &data.collectedData)
					metrics = append(metrics, data.metric)
				}
				data.metricList = metrics
				metricManager.pool.Put(interface{}(data))

				time.Sleep(config.Configuration.PollInterval * time.Second)
			}
		}
	}(ctx, cancel, wg, errCh)
}

func (metricManager *MetricManager) send(ctx context.Context, wg *sync.WaitGroup, _ chan<- error) {
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				data := metricManager.pool.Get().(*interimData)
				_ = metricManager.client.SendBatch(data.metricList)

				time.Sleep(config.Configuration.ReportInterval * time.Second)
			}
		}
	}(ctx, wg)
}

func New(metricList []string, dataCollector interfaces.MetricDataCollector, client interfaces.Client) (MetricManager, error) {
	if len(metricList) == 0 || dataCollector == nil || client == nil {
		return MetricManager{}, ErrInvalidArguments
	}

	var err error
	pool := &sync.Pool{
		New: func() interface{} {
			s := interimData{}
			s.metricList = make([]generaldto.Metrics, 0, len(metricList))
			err = initStrategies(&s.strategyMap, metricList)
			return &s
		},
	}

	pool.New()
	if err != nil {
		return MetricManager{}, err
	}

	metricManager := MetricManager{
		metricsToCollect: metricList,
		dataCollector:    dataCollector,
		client:           client,
		pool:             pool,
	}

	return metricManager, nil
}

func initStrategies(strategyMap *map[string]interfaces.ParsingStrategy, metricNames []string) error {
	*strategyMap = make(map[string]interfaces.ParsingStrategy, len(metricNames))

	var err error
	for _, metricName := range metricNames {
		(*strategyMap)[metricName], err = parser.GetStrategy(metricName)
		if err != nil {
			return err
		}
	}

	return nil
}
