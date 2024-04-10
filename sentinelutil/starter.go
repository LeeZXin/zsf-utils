package sentinelutil

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/logging"
)

func init() {
	sentinel.InitDefault()
	logging.ResetGlobalLogger(new(EmptyLogger))
}
