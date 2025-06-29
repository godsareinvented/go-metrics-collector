package metric

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
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

func (metricManager *MetricManager) CollectAndSend(ctx context.Context) error {
	data := metricManager.pool.Get().(*interimData)
	if metricManager.metricsToCollect == nil || len(metricManager.metricsToCollect) == 0 || metricManager.dataCollector == nil || data.strategyMap == nil {
		return ErrNotInitialised
	}

	go metricManager.collect(ctx)
	go metricManager.send(ctx)

	return nil
}

func (metricManager *MetricManager) collect(ctx context.Context) {
	var data *interimData
	var metricName string
	var metrics []generaldto.Metrics

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data = metricManager.pool.Get().(*interimData)

			metricManager.dataCollector.CollectMetricData(&data.collectedData)

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
}

func (metricManager *MetricManager) send(ctx context.Context) {
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
