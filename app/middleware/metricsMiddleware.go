package middleware

import (
	"cobaMetrics/app/config"
	"cobaMetrics/metrics"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(config config.IConfig, metrics *metrics.MetricsApp) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := string(ctx.Request().URI().FullURI())
		timer := prometheus.NewTimer(metrics.DurationReq.WithLabelValues(path))
		defer timer.ObserveDuration()
		fmt.Println("masuk middleware :", path)

		metrics.CounterReq.WithLabelValues(path).Inc()

		ctx.Next()

		return nil
	}
}
