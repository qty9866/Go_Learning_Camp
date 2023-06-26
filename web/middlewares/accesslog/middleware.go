package accesslog

import (
	"Go_Learning_Camp/web"
	"encoding/json"
)

type MiddleWareBuilder struct {
	logFunc func(log string)
}

type Middleware func(next web.HandleFunc) web.HandleFunc

func (m *MiddleWareBuilder) LogFunc(fn func(log string)) *MiddleWareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddleWareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					Path:       ctx.Req.URL.Path,
					HttpMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				m.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}

// 记录日志
type accessLog struct {
	Host string `json:"host,omitempty"`
	// 命中的路由
	Route      string `json:"route,omitempty"`
	HttpMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
