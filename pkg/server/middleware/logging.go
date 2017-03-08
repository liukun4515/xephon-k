package middleware

import (
	"time"

	dlog "github.com/dyweb/gommon/log"
	"github.com/xephonhq/xephon-k/pkg/server/service"
	"github.com/xephonhq/xephon-k/pkg/util"
)

type LoggingInfoServiceMiddleware struct {
	service.InfoService
	logger *dlog.Entry
}

func NewLoggingInfoServiceMiddleware(service service.InfoService) LoggingInfoServiceMiddleware {
	logger := util.Logger.NewEntryWithPkg("k.s.m")
	return LoggingInfoServiceMiddleware{InfoService: service, logger: logger}
}

// FIXME: the naming here is misleading, the info actually return all the info, more than just version
// and how to hand things like info/version in go-kit
func (mw LoggingInfoServiceMiddleware) Version() string {
	defer func(begin time.Time) {
		// TODO: human readable time format, what's the number, ms, ns?
		mw.logger.Infof("GET /info %d", time.Since(begin))
	}(time.Now())
	return mw.InfoService.Version()
}