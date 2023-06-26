//go:build e2e

package accesslog

import (
	"Go_Learning_Camp/web"
	"fmt"
	"testing"
)

func TestMiddleWareBuilder1(t *testing.T) {
	builder := MiddleWareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("Hello,it's me"))
	})
	server.Start(":8081")
}
