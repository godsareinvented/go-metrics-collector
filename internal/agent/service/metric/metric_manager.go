package metric

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/config"
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"time"
)

type MetricManager struct {
	MetricList          []string
	MetricDataCollector interfaces.MetricDataCollector
	Client              interfaces.Client
	strategies          map[string]interfaces.ParsingStrategy
}

var (
	metricList []generaldto.Metrics
)

func (metricManager *MetricManager) CollectAndSend(ctx context.Context) {
	if metricManager.MetricList == nil {
		panic("metric list is empty")
	}
	if metricManager.Client == nil {
		panic("client is required")
	}

	go metricManager.collect(ctx)
	go metricManager.send(ctx)
}

func (metricManager *MetricManager) Init() {
	metricManager.strategies = make(map[string]interfaces.ParsingStrategy)

	for _, metricName := range metricManager.MetricList {
		metricManager.strategies[metricName] = parser.GetStrategy(metricName)
	}
}

func (metricManager *MetricManager) collect(ctx context.Context) {
	var metricCollectedData agentdto.CollectedMetricData
	var metrics []generaldto.Metrics
	for {
		select {
		case <-ctx.Done():
			return
		default:
			metricManager.MetricDataCollector.CollectMetricData(&metricCollectedData)

			metrics = make([]generaldto.Metrics, 0, len(metricManager.MetricList))
			for _, metricName := range metricManager.MetricList {
				strategy := metricManager.strategies[metricName]
				m := strategy.GetMetric(metricName, metricCollectedData)
				metrics = append(metrics, m)
			}
			metricList = metrics

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
			_ = metricManager.Client.SendBatch(metricList)

			time.Sleep(config.Configuration.ReportInterval * time.Second)
		}
	}
}
