package web

import "testing"

// 进行端到端的测试

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello World from Hud"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("Hello user from Hud"))
	})

	err := s.Start(":8081")
	if err != nil {
		//todo handle error
	}
}
