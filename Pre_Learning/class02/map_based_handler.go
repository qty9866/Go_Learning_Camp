package main

import "net/http"

type HandlerBasedOnMap struct {
	// key 应该是method+Url
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter, request http.Request) {
	key := h.key(request)
}


func (s sdkHttpServer)