package tracing

import (
	"cobaMetrics/app/config"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"io"
)

func NewJaegerTracing(config config.IConfig) (opentracing.Tracer, io.Closer) {
	jaegerCfg := jaegerConfig.Configuration{
		ServiceName: config.Config().Jaeger.ServiceName,
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%v:%v", config.Config().Jaeger.Host, config.Config().Jaeger.Port),
		},
	}

	tracer, closer, err := jaegerCfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("cant connect to jaeger : %v", err)
	}

	return tracer, closer
}
