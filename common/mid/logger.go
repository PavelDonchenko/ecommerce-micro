package mid

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
)

// Logger writes information about the request to the logs.
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.SetTraceID(c.Request.Context(), trace.SpanContextFromContext(c.Request.Context()).TraceID().String())
		c.Request = c.Request.WithContext(ctx)
		now := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		remoteAddr := c.Request.RemoteAddr

		log.Info(c.Request.Context(), "[REQUEST STARTED]", "method", method, "path", path, "remote address", remoteAddr)

		c.Next()

		statusCode := c.Writer.Status()
		end := time.Since(now).String()

		if statusCode < http.StatusInternalServerError {
			log.Error(c.Request.Context(), "[REQUEST COMPLETED]", "method", method, "path", path, "remoteaddr", remoteAddr,
				"statuscode", statusCode, "duration", end)
			return
		}

		log.Info(c.Request.Context(), "[REQUEST INTERNAL ERROR]", "method", method, "path", path, "remoteaddr", remoteAddr,
			"statuscode", statusCode, "duration", end, "errors")

	}
}
