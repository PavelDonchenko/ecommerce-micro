package mid

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PavelDonchenko/ecommerce-micro/common/metric"
)

func Metrics(service metric.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		appMetric := metric.NewHttpMetrics(c.Request.URL.Path, c.Request.Method)
		appMetric.Started()
		c.Next()

		response := c.Writer
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(response.Status())
		service.SaveHttp(appMetric)
	}
}

func MetricsHandler() gin.HandlerFunc {
	handler := promhttp.Handler()

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
