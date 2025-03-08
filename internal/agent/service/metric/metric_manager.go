package metric

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/buisness_logic/parser"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/threading_pattern"
	"sync"
	"time"
)

type (
	interimData struct {
		err           error
		metric        dto.Metrics
		metricList    []dto.Metrics
		collectedData agentDto.CollectedMetricData
		strategyMap   map[string]interfaces.ParsingStrategyInterface
	}

	// MetricManager todo: Стоит переформатировать названия пакетов: убрать упоминание пакета вначале структур, убрать использование снейк_кейс, камелКейс и т.д.
	MetricManager struct {
		ctx                  context.Context                         `validate:"required"`
		cancel               context.CancelFunc                      `validate:"required"`
		metricNamesToCollect []string                                `validate:"required,unique,min=1,max=1031,dive,metric_name"`
		dataCollector        interfaces.MetricDataCollectorInterface `validate:"required"`
		client               interfaces.ClientInterface              `validate:"required"`
		interimDataPool      *sync.Pool                              `validate:"required"`
	}
)

func (m *MetricManager) CollectAndSend() (*sync.WaitGroup, chan error) {
	errCh := make(chan error)

	err := config.Configuration.Validate.Struct(*m)
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

		for {
			select {
			case <-m.ctx.Done():
				return
			default:
				metricList, err := m.collectMetrics()
				if nil != err {
					errCh <- err
					return
				}

				ch <- metricList

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
		err := m.client.SendBatch(m.ctx, metricList)
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

func (m *MetricManager) collectMetrics() (*[]dto.Metrics, error) {
	d := m.interimDataPool.Get().(*interimData)

	d.collectedData = agentDto.CollectedMetricData{}
	d.err = m.dataCollector.CollectMetricData(m.ctx, &d.collectedData)
	if nil != d.err {
		return &d.metricList, nil
	}

	d.metricList = d.metricList[:0]
	for _, metricName := range m.metricNamesToCollect {
		d.metric = dto.Metrics{}
		d.err = d.strategyMap[metricName].ParseMetric(&d.metric, &d.collectedData)
		if nil != d.err {
			return &d.metricList, nil
		}
		d.metricList = append(d.metricList, d.metric)
	}

	return &d.metricList, nil
}

func initStrategies(strategyMap *map[string]interfaces.ParsingStrategyInterface, metricNames []string) error {
	*strategyMap = make(map[string]interfaces.ParsingStrategyInterface, len(metricNames))

	var err error
	for _, metricName := range metricNames {
		(*strategyMap)[metricName], err = parser.GetStrategy(metricName)
		if nil != err {
			return err
		}
	}

	return nil
}

func NewMetricManager(
	ctx context.Context,
	metricNameList []string,
	dataCollector interfaces.MetricDataCollectorInterface,
	client interfaces.ClientInterface,
) (MetricManager, error) {
	var err error

	pool := &sync.Pool{
		New: func() interface{} {
			s := interimData{}
			s.metricList = make([]dto.Metrics, 0, len(metricNameList))
			err = initStrategies(&s.strategyMap, metricNameList)
			return &s
		},
	}
	pool.New()
	if nil != err {
		return MetricManager{}, err
	}

	wrappedCtx, cancel := context.WithCancel(ctx)

	metricManager := MetricManager{
		ctx:                  wrappedCtx,
		cancel:               cancel,
		metricNamesToCollect: metricNameList,
		dataCollector:        dataCollector,
		client:               client,
		interimDataPool:      pool,
	}
	err = config.Configuration.Validate.Struct(metricManager)
	if nil != err {
		return MetricManager{}, err
	}

	return metricManager, nil
}
