package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	QueryValues url.Values
}

type StringValue struct {
	val string
	err error
}

func (c *Context) BindJson(val any) error {
	if val == nil {
		return errors.New("web: Input can't be nil")
	}
	if c.Req.Body == nil {
		return errors.New("web: Body can't be nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	// useNumber => 数字就会用Number来表示
	// 否则默认是float64
	decoder.UseNumber()
	// 如果要是有一个未知的字段，就会报错
	// 比如说你User只有Name和Email两个字段
	// JSON里面额外多了一个Age字段，那么就会报错
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	/*
		val, ok := c.Req.Form[key]
		if !ok {
			return "", errors.New("web: Key Not Found")
		}
		//val的类型是	[]string
		return val[0], nil
	*/

	// 直接调用API 这里其实也就是帮助用户调用了ParseForm()
	return c.Req.FormValue(key), nil
}

// QueryValue Query和表单比起来，它没有缓存
// 所以要考虑在这中间缓存
func (c *Context) QueryValue(key string) (string, error) {
	if c.QueryValues == nil {
		c.QueryValues = c.Req.URL.Query()
	}
	// 用户区别不出来是没有值，还是恰好是空的字符串
	return c.QueryValues.Get(key), nil
}

/*func (c *Context) QueryValue1(key string) (string, error) {
	params := c.Req.URL.Query()
	if params == nil {
		return "", errors.New("web: 没有任何查询参数")
	}
	val, ok := params[key]
	if !ok {
		return "", errors.New("web: 找不到这个key")
	}
	return val[0], nil
}*/

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key不存在"),
		}
	}
	return StringValue{
		val: val,
		err: nil,
	}

}

func (c *Context) QueryValueV1(key string) StringValue {
	if c.QueryValues == nil {
		c.QueryValues = c.Req.URL.Query()
	}
	val, ok := c.QueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: key 不存在")}
	}
	return StringValue{val: val[0]}
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
