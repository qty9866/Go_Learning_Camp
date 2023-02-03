package main

import (
	"log"
	"net/http"
)

type Server interface {
	Route(pattern string, handleFunc http.HandlerFunc)
	Start(address string) error
}

// sdkHttpServer Struct 基于http库实现
type sdkHttpServer struct {
	Name string
}

// Route 注册路由
func (s *sdkHttpServer) Route(
	method string,
	pattern string,
	handleFunc func(ctx *Context) ) {
	key := s.
}
func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, nil)
}

func NewHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	req := &signUpReq{}

	ctx := Context{
		w,
		r,
	}
	err := ctx.ReadJson(req)
	if err != nil {
		log.Fatal("read json error:", err)
	}

	resp := &commonResponse{
		Data: 123,
	}
	err = ctx.WriteJson(http.StatusOK, resp)
	if err != nil {
		log.Fatal("write json data error", err)
	}

}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
