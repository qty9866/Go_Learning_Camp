package web

import (
	"net/http"
)

type HandleFunc func(ctx *Context)

type Middleware func(next HandleFunc) HandleFunc

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
	mdl []Middleware
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP HTTPServer 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 最后一个是这个
	root := h.serve

	// 然后这里就是利用最后一个不断往前面回溯组装链条 从后往前
	// 把后一个作为前一个的 next 构造号链条
	for i := len(h.mdl) - 1; i >= 0; i-- {
		root = h.mdl[i](root)
	}
	root(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.n.handler == nil {
		// 路由没有命中 就是404
		ctx.Resp.WriteHeader(404)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}
}

// Start 启动服务器
func (h *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

func (h *HTTPServer) Post(path string, handler HandleFunc) {
	h.addRoute(http.MethodPost, path, handler)
}

func (h *HTTPServer) Get(path string, handler HandleFunc) {
	h.addRoute(http.MethodGet, path, handler)
}
