package metric

import (
	"context"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/threading_pattern"
	"sync"
	"time"
)

// MetricManager todo: Стоит переформатировать названия пакетов: убрать упоминание пакета вначале структур, убрать использование снейк_кейс, камелКейс и т.д.
type MetricManager struct {
	ctx                  context.Context
	cancel               context.CancelFunc
	metricNamesToCollect []string
	dataCollector        interfaces.MetricDataCollectorInterface
	client               interfaces.ClientInterface
}

var (
	ErrNoInitializedArguments = errors.New("all metric manager arguments must be initialized")
	ErrNoMetricsToCollect     = errors.New("empty metric name list. No metrics to collect")

	// Переменные вынесены в область видимости пакета для избежагания постоянной инициализации переменных в функциях и выделения памяти
	collectingErr       error
	strategies          map[string]interfaces.ParsingStrategyInterface
	metric              dto.Metrics
	metricList          []dto.Metrics
	collectedMetricData agentDto.CollectedMetricData
)

func (m *MetricManager) CollectAndSend() (*sync.WaitGroup, chan error) {
	errCh := make(chan error)

	err := m.validateArguments()
	if nil != err {
		errCh <- err
		close(errCh)
		return &sync.WaitGroup{}, errCh
	}

	ch := m.startCollecting(errCh)
	wg := m.send(ch, errCh)

	go func(wg *sync.WaitGroup) {
		defer close(errCh)
		defer m.cancel()

		wg.Wait()
	}(wg)

	return wg, errCh
}

func (m *MetricManager) startCollecting(errCh chan<- error) chan *[]dto.Metrics {
	ch := make(chan *[]dto.Metrics)

	go func(ch chan *[]dto.Metrics, errCh chan<- error) {
		defer close(ch)

		var err error
		for {
			select {
			case <-m.ctx.Done():
				return
			default:
				err = m.collectMetrics()
				if nil != err {
					errCh <- err
					return
				}

				ch <- &metricList

				if config.Configuration.PollInterval > 0 {
					time.Sleep(time.Duration(config.Configuration.PollInterval) * time.Second)
				}
			}
		}
	}(ch, errCh)

	return ch
}

func (m *MetricManager) send(inputCh <-chan *[]dto.Metrics, errCh chan<- error) *sync.WaitGroup {
	ch := make(chan *[]dto.Metrics)
	wg, _ := threading_pattern.InitWorkerPool(m.ctx, config.Configuration.RateLimit, ch, func(_ int, metricList *[]dto.Metrics) error {
		err := m.client.SendBatch(metricList)
		if nil != err {
			errCh <- err
			m.cancel()
		}
		return err
	})

	go func(ch chan *[]dto.Metrics, inputCh <-chan *[]dto.Metrics, wg *sync.WaitGroup) {
		defer close(ch)

		var metricNameList *[]dto.Metrics
		for metricNameList = range inputCh {
			select {
			case <-m.ctx.Done():
				return
			default:
				ch <- metricNameList

				if config.Configuration.ReportInterval > 0 {
					time.Sleep(time.Duration(config.Configuration.ReportInterval) * time.Second)
				}
			}

		}
	}(ch, inputCh, wg)

	return wg
}

func (m *MetricManager) collectMetrics() error {
	collectedMetricData = agentDto.CollectedMetricData{}
	collectingErr = m.dataCollector.CollectMetricData(m.ctx, &collectedMetricData)
	if nil != collectingErr {
		return collectingErr
	}

	metricList = metricList[:0]
	for _, metricName := range m.metricNamesToCollect {
		metric = dto.Metrics{}
		collectingErr = strategies[metricName].ParseMetric(&metric, &collectedMetricData)
		if nil != collectingErr {
			return collectingErr
		}
		metricList = append(metricList, metric)
	}

	return nil
}

func (m *MetricManager) initStrategies() error {
	strategies = make(map[string]interfaces.ParsingStrategyInterface, len(m.metricNamesToCollect))

	var err error
	for _, metricName := range m.metricNamesToCollect {
		strategies[metricName], err = parser.GetStrategy(metricName)
		if nil != err {
			return err
		}
	}

	return nil
}

func (m *MetricManager) validateArguments() error {
	if nil == m.dataCollector || nil == m.client || nil == m.metricNamesToCollect {
		return ErrNoInitializedArguments
	}
	if len(m.metricNamesToCollect) == 0 {
		return ErrNoMetricsToCollect
	}
	return nil
}

func NewMetricManager(
	ctx context.Context,
	metricNameList []string,
	dataCollector interfaces.MetricDataCollectorInterface,
	client interfaces.ClientInterface,
) (MetricManager, error) {
	wrappedCtx, cancel := context.WithCancel(ctx)

	metricManager := MetricManager{
		metricNamesToCollect: metricNameList,
		dataCollector:        dataCollector,
		client:               client,
		ctx:                  wrappedCtx,
		cancel:               cancel,
	}
	err := metricManager.validateArguments()
	if nil != err {
		return MetricManager{}, err
	}

	metricList = make([]dto.Metrics, 0, len(metricList))

	err = metricManager.initStrategies()
	if nil != err {
		return MetricManager{}, err
	}

	return metricManager, nil
}
