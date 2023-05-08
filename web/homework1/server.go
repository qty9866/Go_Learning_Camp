package web

import "net/http"

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr是监听地址 如果是只指定端口 可以使用

	Start(addr string) error

	addRoute(method string, path string, handler HandleFunc)
}

// 确保HTTPServer结构体实现了Server接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: NewRouter(),
	}
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("404 NotFound!"))
		return
	}
	ctx.PathParams = mi.pathParams
	mi.n.handler(ctx)
}
