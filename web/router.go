package web

import (
	"Go_Learning_Camp/web/v1"
	"fmt"
	"strings"
)

// 用来支持对路由树的操作
// 代表路由树(森林)

type router struct {
	// Beego Gin HTTP Method 对应一棵树
	// GET有一棵树，POST也有一棵树

	// http method => 路由树根节点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// 加一些限制
// path必须以"/"开头,不能以"/"结尾，中间也不能有连续的//
func (r *router) addRoute(method string, path string, handleFunc web.HandleFunc) {
	if path == "" {
		panic("web:路径不能为空字符串")
	}

	// 首先找到树
	root, ok := r.trees[method]
	if !ok {
		// 说明还没有根节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	// 开头不能没有/
	if path[0] != '/' {
		panic("web:路径必须以\"/\"'开头")
	}

	// 结尾
	if path != "/" && path[len(path)-1] == '/' {
		panic("web:路径不能以\"/\"结尾")
	}

	// 中间连续"//",可以用strings.contains("//")
	// 根节点处理一下
	if path == "/" {
		if root.handler != nil {
			panic("路由冲突，重复注册[/]")
		}
		root.handler = handleFunc
		return
	}
	// /user/home 被割成三段
	// 切割这个path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web:不能有连续的/")
		}
		// 递归下去 找准位置
		// 如果中途有节点不存在，你就要创建出来
		child := root.childOfCreate(seg)
		root = child
	}

	if root.handler != nil {
		panic(fmt.Sprintf("web:路由冲突,重复注册[%s]", path))
	}
	root.handler = handleFunc
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	// 基本上也是沿着树深度遍历查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n:          root,
			pathParams: nil,
		}, true
	}
	//这里把前置和后置的"/"都去掉
	path = strings.Trim(path, "/")

	// 按照"/"切割
	segs := strings.Split(path, "/")
	pathParams := make(map[string]string)
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 命中了路径参数
		if paramChild {
			// 因为path是":id"这种形式，所以要去掉：
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	// 代表我确实有这个节点
	// 但是节点是不是用户注册的有Handler方法的我就不知道了
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

type node struct {
	path string
	// 静态匹配的节点
	// path 到子节点的映射
	children map[string]*node

	// 加一个通配符匹配
	starChild *node
	// 缺一个代表用户注册的业务逻辑
	handler    web.HandleFunc
	paramChild *node
}

func (n *node) childOfCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return res, false, ok
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
