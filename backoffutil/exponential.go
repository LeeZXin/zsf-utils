package backoffutil

import (
	"github.com/LeeZXin/zsf-utils/randutil"
	"time"
)

var defaultBackoffStrategy = newDefaultExponentialStrategy()

type Strategy interface {
	Backoff(retries int) time.Duration
}

func newDefaultExponentialStrategy() Strategy {
	return &Exponential{Config: DefaultConfig}
}

type Exponential struct {
	Config Config
}

func (bc *Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	backoff *= 1 + bc.Config.Jitter*(randutil.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}

func Backoff(retries int) time.Duration {
	return defaultBackoffStrategy.Backoff(retries)
}
