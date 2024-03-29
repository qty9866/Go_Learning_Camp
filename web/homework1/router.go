package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

func NewRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type nodeType int

type node struct {
	// node类型
	nodeType nodeType

	// path 注册的路由
	path string

	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string
	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp

	// handler 命中路由之后执行的逻辑
	handler HandleFunc
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	if path == "" {
		panic("web: 路由为空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须由\"/\"开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以\"/\"结尾")
	}

	root, ok := r.trees[method]
	if !ok {
		//这个方法的路由树为空 创建跟节点
		root = &node{path: "/"}
		r.trees[method] = root
	}

	if path == "/" {
		if root.handler != nil {
			panic("web: 根节点重复注册,路由冲突")
		}
		root.handler = handler
		return
	}
	// 下面对路由进行切分
	segs := strings.Split(path[1:], "/")
	// 从这里开始进行一段一段处理
	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOfCreate(seg)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由注册冲突[%s]", path))
	}
	root.handler = handler
}

// 路由查找 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 根节点
	if path == "/" {
		return &matchInfo{n: root}, true
	}

	// 接续拆分路径
	segs := strings.Split(path[1:], "/")
	mi := &matchInfo{}
	for _, seg := range segs {
		var child *node
		// findChild 返回子节点
		// ok==false代表没有找到这个路径的子节点
		child, ok = root.findChild(seg)
		if !ok {
			if root.nodeType == nodeTypeAny {
				mi.n = root
				return mi, true
			}
			return nil, false
		}
		if child.paramName != "" {
			mi.addValue(child.paramName, seg)
		}
		root = child
	}
	mi.n = root
	return mi, true
}

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

// findChild 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否命中
func (n *node) findChild(path string) (*node, bool) {
	if n.children == nil {
		return n.childOfNonStatic(path)
	}
	res, ok := n.children[path]
	if !ok {
		return n.childOfNonStatic(path)
	}
	return res, ok
}

// childOfNonStatic 从非静态匹配的子节点里面查找
func (n *node) childOfNonStatic(path string) (*node, bool) {
	if n.regChild != nil {
		if n.regChild.regExpr.Match([]byte(path)) {
			return n.regChild, true
		}
	}
	if n.paramChild != nil {
		return n.paramChild, true
	}
	return n.starChild, n.starChild != nil
}

// childOfCreate 查找子节点，
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOfCreate(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由,已有路径参数路由,不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.regChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", path))
		}
		if n.starChild == nil {
			n.starChild = &node{path: path, nodeType: nodeTypeAny}
		}
		return n.starChild
	}

	// 以 : 开头，需要进一步解析，判断是参数路由还是正则路由
	if path[0] == ':' {
		paramName, expr, isReg := n.parseParam(path)
		if isReg {
			return n.childOfCreateReg(path, expr, paramName)
		}
		return n.childOfCreateParam(path, paramName)
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{path: path, nodeType: nodeTypeStatic}
		n.children[path] = child
	}
	return child
}

// parseParam 用于解析判断是不是正则表达式
// 第一个返回值是参数名字
// 第二个返回值是正则表达式
// 第三个返回值为 true 则说明是正则路由
func (n *node) parseParam(path string) (string, string, bool) {
	// 首先去除":"
	path = path[1:]
	segs := strings.SplitN(path, "(", 2)
	if len(segs) == 2 {
		expr := segs[1]
		if strings.HasSuffix(expr, ")") {
			return segs[0], expr[:len(expr)-1], true
		}
	}
	return path, "", false
}

func (n *node) childOfCreateReg(path string, expr string, paramName string) *node {
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
	}
	if n.paramChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.regChild != nil {
		if n.regChild.regExpr.String() != expr || n.paramName != paramName {
			panic(fmt.Sprintf("web: 路由冲突，正则路由冲突，已有 %s，新注册 %s", n.regChild.path, path))
		}
	}
	return n.regChild
}

func (n *node) childOfCreateParam(path string, paramName string) *node {
	if n.regChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
	}
	if n.paramChild != nil {
		if n.paramChild.path != path {
			panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
		}
	} else {
		n.paramChild = &node{path: path, paramName: paramName, nodeType: nodeTypeParam}
	}
	return n.paramChild
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m matchInfo) addValue(key string, value string) {
	//todo
	if m.pathParams == nil {
		// 大多数情况下，参数路径只会有一端
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
