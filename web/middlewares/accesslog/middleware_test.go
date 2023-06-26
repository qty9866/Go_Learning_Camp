package accesslog

import (
	"Go_Learning_Camp/web"
	"fmt"
	"net/http"
	"testing"
)

func TestMiddleWareBuilder(t *testing.T) {
	builder := MiddleWareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerWithMiddleware(mdl))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("Hello,it's me")
	})
	req, _ := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	req.Host = "localhost"
	server.ServeHTTP(nil, req)
}
