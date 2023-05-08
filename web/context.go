package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	queryValues url.Values
}

func (c *Context) RespJsonOk(val any) error {
	return c.RespJson(http.StatusOK, val)
}

func (c *Context) RespJson(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(status)
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))
	n, err := c.Resp.Write(data)
	if n != len(data) { // 说明没写完
		return errors.New("web: 未写入全部数据")
	}
	return err
}

func (c *Context) SetCookie(Cookie *http.Cookie) {
	http.SetCookie(c.Resp, Cookie)
}

type SafeContext struct {
	Context
	mutex sync.RWMutex
}

func (c *SafeContext) RespJsonOk(val any) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.RespJson(http.StatusOK, val)
}
