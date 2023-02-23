package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error
	// AddRoute 路由注册功能
	// method是HTTP方法
	// path 是路由
	// HandleFunc 是你的业务逻辑
	addRoute(method string, path string, handleFunc HandleFunc)
	// AddRoute1 这种允许注册多个，没有必要提供
	// 让用户自己去管理
	AddRoute1(method string, path string, handles ...HandleFunc)
}

//type HTTPSServer struct {
//	HTTPServer
//}

type HTTPServer struct {
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

// ServeHTTP 处理请求的入口
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 你的框架代码就在这里
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 接下来就是查找路由，并且执行命中的业务逻辑
	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {

}

func (h *HTTPServer) AddRoute1(method string, path string, handles ...HandleFunc) {
	//TODO implement me
	panic("implement me")
}

//func (h *HTTPServer) addRoute(method string, path string, handleFunc HandleFunc) {
//	// 这里注册到路由树里面
//	panic("implement me")
//}

func (h *HTTPServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServer) Post(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)

}
func (h *HTTPServer) Options(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

func (h *HTTPServer) Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		// todo handle error
		return err
	}
	// 在这里，可以让用户注册所谓的 after start回调
	// 比如说往你的 admin 注册一下自己这个实例
	// 在这里执行一些业务所需的前置条件
	return http.Serve(listen, h)
}

func (h *HTTPServer) Start1(addr string) error {
	return http.ListenAndServe(addr, h)
}
