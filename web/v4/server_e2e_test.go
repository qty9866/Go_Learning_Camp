package web

import (
	"fmt"
	"testing"
)

// 这里放着端到端测试的代码

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Post("/form", func(ctx *Context) {
		ctx.Req.ParseForm()
		ctx.Resp.Write([]byte("hello form"))
	})

	s.Get("/values/:id", func(ctx *Context) {
		id, err := ctx.PathValueV1("id").AsInt64()
		if err != nil {
			ctx.Resp.WriteHeader(400)
			ctx.Resp.Write([]byte("web: id輸入錯誤"))
			return
		}
		ctx.Resp.Write([]byte(fmt.Sprintf("hello,%d", id)))
	})

	s.Start(":8081")
}
