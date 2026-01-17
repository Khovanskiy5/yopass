package middleware

import (
	"io"
	"net/http"

	"github.com/Khovanskiy5/yopass/internal/utils"
	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

func NewLoggingHandler(logger *zap.Logger, trustedProxies []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.CustomLoggingHandler(nil, next, func(_ io.Writer, params handlers.LogFormatterParams) {
			req := params.Request
			if req == nil {
				logger.Error("unable to log request: no request object")
				return
			}

			host := utils.GetRealClientIP(req, trustedProxies)
			uri := req.RequestURI

			if req.ProtoMajor == 2 && req.Method == "CONNECT" {
				uri = req.Host
			}
			if uri == "" {
				uri = params.URL.RequestURI()
			}

			logger.Info(
				"Request handled",
				zap.String("host", host),
				zap.Time("timestamp", params.TimeStamp),
				zap.String("method", req.Method),
				zap.String("uri", uri),
				zap.String("protocol", req.Proto),
				zap.Int("status", params.StatusCode),
				zap.Int("size", params.Size),
			)
		})
	}
}
