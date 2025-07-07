package prepared_option

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"time"
)

var DefaultFixedDelayListOptions = dto.RetryOptions{
	Attempts:  3,
	DelayList: []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
}
