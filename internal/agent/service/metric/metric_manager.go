package metric

import (
	"context"
	"errors"
	"github.com/hashicorp/go-multierror"
	factory "github.com/oldhanasong/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/threading_pattern"
	"math"
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
		sendingCounter   int
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
	ch := metricManager.collect(ctx, &wg, errCh)
	metricManager.send(ch, &wg, errCh)
}

func (metricManager *MetricManager) collect(ctx context.Context, wg *sync.WaitGroup, errCh chan<- error) chan []generaldto.Metrics {
	ch := collectingChan()

	go func(ctx context.Context, ch chan []generaldto.Metrics, wg *sync.WaitGroup, errCh chan<- error) {
		defer close(ch)
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
					data.strategyMap[metricName].FillMetric(&data.metric, &data.collectedData)
					metrics = append(metrics, data.metric)
				}

				ch <- metrics

				sleep(config.Configuration.PollInterval)
			}
		}
	}(ctx, ch, wg, errCh)

	return ch
}

func (metricManager *MetricManager) send(ch <-chan []generaldto.Metrics, wg *sync.WaitGroup, _ chan<- error) {
	workerCh := make(chan []generaldto.Metrics, config.Configuration.RateLimit)
	_ = threading_pattern.InitWorkerPool(config.Configuration.RateLimit, workerCh, func(_ int, metricBatch []generaldto.Metrics) {
		_ = metricManager.client.SendBatch(metricBatch)
	})

	go func(ch <-chan []generaldto.Metrics, workerCh chan []generaldto.Metrics, wg *sync.WaitGroup) {
		defer close(workerCh)
		defer wg.Done()

		sleep(config.Configuration.ReportInterval)

		var batch []generaldto.Metrics
		var ok bool

		for {
			select {
			case batch, ok = <-ch:
				if !ok {
					return
				}

				workerCh <- batch

				if metricManager.doesReportNeedSleep(ch) {
					sleep(config.Configuration.ReportInterval)
				}
			default:
				sleep(config.Configuration.ReportInterval)
			}
		}
	}(ch, workerCh, wg)
}

func (metricManager *MetricManager) doesReportNeedSleep(ch <-chan []generaldto.Metrics) bool {
	metricManager.sendingCounter++
	cond := len(ch) == 0 || metricManager.sendingCounter >= config.Configuration.RateLimit
	if cond {
		metricManager.sendingCounter = 0
	}
	return cond
}

func New(metricList []string, dataCollector interfaces.MetricDataCollector, client interfaces.Client) (MetricManager, error) {
	if dataCollector == nil || client == nil {
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
		if (*strategyMap)[metricName], err = factory.Strategy(metricName); err != nil {
			return err
		}
	}

	return nil
}

func collectingChan() chan []generaldto.Metrics {
	chLen := int(math.Max(
		1,
		math.Ceil(config.Configuration.ReportInterval.Seconds()/config.Configuration.PollInterval.Seconds()),
	))

	return make(chan []generaldto.Metrics, chLen)
}

func sleep(d time.Duration) {
	if d > 0 {
		time.Sleep(d * time.Second)
	}
}
