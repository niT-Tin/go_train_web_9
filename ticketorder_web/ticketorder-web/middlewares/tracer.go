package middlewares

import (
	"fmt"
	"gotrains/ticketorder_web/ticketorder-web/global"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:           true,
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
		}
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		startSpan := tracer.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()
		c.Set("tracer", tracer)
		c.Set("parentSpan", startSpan)
		c.Next()
	}
}
