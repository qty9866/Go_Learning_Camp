package web

import (
	"Go_Learning_Camp/web/v1"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	// 第一个步骤是构造路由树
	// 第二个步骤是验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			http.MethodGet,
			"/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	var mockHandler web.HandleFunc = func(ctx *web.Context) {}
	r := newRouter()
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 接下来就在这里断言两者相等
	wantRouter := &router{
		map[string]*node{
			http.MethodGet: {
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": {
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {
								path:    "detail",
								handler: mockHandler,
							},
						},
						starChild: &node{
							path:    "*",
							handler: mockHandler,
						},
					},
				},
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {
						path: "order",
						children: map[string]*node{
							"create": {
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": {
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	fmt.Println(msg, ok)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	}, "web:路径必须以\"/\"开头")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	}, "web:路径不能以\"/\"开头")

	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a//b/c/", mockHandler)
	}, "web:路径不能以\"/\"开头")

	// 根节点重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	}, "路由冲突,重复注册[/]")

	// 普通节点重复注册
	r = newRouter()
	r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	assert.Panicsf(t, func() {
		r.addRoute(http.MethodGet, "/a/b/c", mockHandler)
	}, "路由冲突,重复注册[/a/b/c]")
	// 这个是不行的 因为HandleFunc是不可比的
	//assert.Equal(t,wantRouter,r)
}

// 返回一个错误信息，bool代表是否真的相等
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的HTTP method"), false
		}
		//v,dst要相等
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	// 比较handler
	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler不相等"), false
	}

	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点%s不存在", path), false
		}
		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

//func TestRouter_findRoute(t *testing.T) {
//	testRoute := []struct {
//		method string
//		path   string
//	}{
//		{
//			method: http.MethodGet,
//			path:   "/",
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user",
//		},
//		{
//			method: http.MethodGet,
//			path:   "/user/home",
//		},
//		{
//			method: http.MethodGet,
//			path:   "/order/detail",
//		},
//		{
//			method: http.MethodPost,
//			path:   "/order/create",
//		},
//		{
//			method: http.MethodPost,
//			path:   "/login",
//		},
//	}
//
//	r := newRouter()
//	var mockHandler HandleFunc = func(ctx Context) {
//	}
//
//	for _, route := range testRoute {
//		r.addRoute(route.method, route.path, mockHandler)
//	}
//
//	testCases := []struct {
//		name   string
//		method string
//		path   string
//
//		wantFound bool
//		wantNode  *node
//	}{
//		{
//			// 方法都不存在
//			name:   "method not existed",
//			method: http.MethodOptions,
//			path:   "/user/home",
//		},
//		{
//			// 完全命中
//			name:   "order detail",
//			method: http.MethodGet,
//			path:   "/user/home",
//
//			wantFound: true,
//			wantNode: &node{
//				handler: mockHandler,
//				path:    "detail",
//			},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			n, found := r.findRoute(tc.method, tc.path)
//			assert.Equal(t, tc.wantFound, found)
//			if !found {
//				return
//			}
//			assert.Equal(t, tc.wantNode.path, n.path)
//			assert.Equal(t, tc.wantNode.children, n.children)
//			nHandler := reflect.ValueOf(n.handler)
//			yHandler := reflect.ValueOf(tc.wantNode.handler)
//			assert.Equal(t, nHandler, yHandler)
//		})
//	}
//}

func TestRouter_findRoute(t *testing.T) {
	testRoute := []struct {
		method string
		path   string
	}{
		{
			http.MethodGet,
			"/",
		},
		{
			http.MethodDelete,
			"/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	r := newRouter()
	var mockHandler web.HandleFunc = func(ctx *web.Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name      string
		method    string
		path      string
		wantFound bool
		wantNode  *node
		info      *matchInfo
	}{
		{
			// 方法都不存在
			name:   "method not found",
			method: http.MethodOptions,
			path:   "/order/detail",
		},
		{
			// 完全命中
			name:      "order detail",
			method:    http.MethodGet,
			path:      "/order/detail",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "*",
				},
			},
		},
		{
			// 根节点
			name:      "root",
			method:    http.MethodDelete,
			path:      "/",
			wantFound: true,
			info: &matchInfo{
				n: &node{
					handler: mockHandler,
					path:    "*",
				},
			},
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
			},
		},
		{
			// 根节点
			name:   "path not found",
			method: http.MethodGet,
			path:   "/aabbcc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wantFound, found)
			if !found {
				return
			}
			assert.Equal(t, tc.info.pathParams, mi.pathParams)
			n := mi.n
			wantVal := reflect.ValueOf(tc.info.n.handler)
			nVal := reflect.ValueOf(n.handler)
			assert.Equal(t, wantVal, nVal)
		})
	}
}
