package web

import (
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
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
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

func (r *router) findRoute(method string, path string) (*node, bool) {
	// 基本上也是沿着树深度遍历查找下去
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return root, true
	}
	//这里把前置和后置的"/"都去掉
	path = strings.Trim(path, "/")

	// 按照"/"切割
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	// 代表我确实有这个节点
	// 但是节点是不是用户注册的有Handler方法的我就不知道了
	return root, true
}

type node struct {
	path string
	// children 子节点
	// 子节点的path=>node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符*表达的节点，任意匹配
	starChild *node
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

// 通配符更改
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.starChild, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		return n.starChild, n.starChild != nil
	}
	return child, ok
}
