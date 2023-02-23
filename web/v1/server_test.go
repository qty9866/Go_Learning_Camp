package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := &HTTPServer{}           // NewServer
	var _ Server = &HTTPServer{} // NewServer

	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		fmt.Println("处理第一件事")
		fmt.Println("处理第二件事")
	})

	handler1 := func(ctx Context) {
		fmt.Println("处理第一件事")
	}
	handler2 := func(ctx Context) {
		fmt.Println("处理第一件事")
	}

	// 用户自己去管这种
	h.AddRoute(http.MethodGet, "/user", func(ctx Context) {
		handler1(ctx)
		handler2(ctx)
	})

	h.Get("/user", func(ctx Context) {

	})

	h.AddRoute1(http.MethodGet, "/user", handler1, handler2)
	// 方法一：用户完全委托给http包
	err := http.ListenAndServe(":8080", h)
	if err != nil {
		// todo handle error
	}
	http.ListenAndServeTLS(":443", "", "", h)

	// 方法二：用户自己手动管理
	h.Start("8081")
}
