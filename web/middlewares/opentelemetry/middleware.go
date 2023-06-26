package opentelemetry

import (
	"Go_Learning_Camp/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "Go_Learning_Camp/web/middlewares/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

/*//NewMiddlewareBuilder 一般来说Builder模式都是返回指针
func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		Tracer: tracer,
	}
}*/

func (m MiddlewareBuilder) Build() web.Middleware {
	// 如果用户没有调用tracer可以初始化一个tracer
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			// 尝试和客户端的 trace 结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			_, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
			// 直接调用下一步
			next(ctx)

			// 这个是只有执行完 next 才可能有值
			span.SetName(ctx.MatchedRoute)
		}
	}
}
