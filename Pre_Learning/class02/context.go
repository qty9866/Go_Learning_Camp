package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	R      *http.Request
}

// ReadJson 可以接受任何类型的参数
func (c *Context) ReadJson(req interface{}) error {
	// 帮我读了body
	// 帮我反序列化
	r := c.R
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		fmt.Println("error when unmarshal data", err)
		return err
	}
	return nil
}

func (c *Context) WriteJson(code int, resp interface{}) error {
	c.Writer.WriteHeader(code)
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("marshal data error", err)
	}
	_, err = c.Writer.Write(data)
	if err != nil {
		fmt.Println("write data error:", err)
		return err
	}
	return nil
}
