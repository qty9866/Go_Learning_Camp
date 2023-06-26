package web

import "Go_Learning_Camp/web/v5"

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

//type MiddlewareV1 interface {
//	Invoke()
//}
//
//// Interceptor 拦截器
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}

type Chain []web.HandleFunc

type HandleFuncV1 func(ctx *web.Context) (next bool)

type ChainV1 struct {
	handler []HandleFuncV1
}

func (c ChainV1) Run(ctx *web.Context) {
	for _, handleFunc := range c.handler {
		next := handleFunc(ctx)
		// 这种是中断执行
		if !next {
			return
		}
	}
}

type Net struct {
	handler []HandleFuncV1
}
