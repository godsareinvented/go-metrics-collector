package prepared_option

import (
	"github.com/godsareinvented/go-metrics-collector/internal/server/dto"
	"time"
)

var DefaultFixedDelayListOptions = dto.RetryOptions{
	Attempts:  3,
	DelayList: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
}
