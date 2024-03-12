package middleware

import (
	"cobaMetrics/app/config"
	"cobaMetrics/metrics"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(config config.IConfig, metrics *metrics.MetricsApp) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := string(ctx.Request().URI().Path())
		method := ctx.Method()
		timer := prometheus.NewTimer(metrics.DurationReq.WithLabelValues(path, method))
		defer timer.ObserveDuration()
		log.Infof("masuk middleware : %v %v", method, path)

		metrics.CounterReq.WithLabelValues(path, method).Inc()

		ctx.Next()

		return nil
	}
}
