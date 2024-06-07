package web

import (
	"fmt"
	"strings"
)

// 代表路由
// 这种有明确概念的东西最好把它剥离开来，最好把server做的尽量清爽，让它职责尽量单一
type router struct {
	// 为什么要是trees，因为我们要解决GET POST DELETE 的问题，
	// 有两种方法，一种是在树的节点里标记支持什么http method，一种是让 GET POST DELETE 分别作为根节点一个http method一棵树。
	// 我们用第二钟比较简单，到后面这两种都可以试一下，你不写代码，我跟你说第二钟方法简单你可能没什么体会。
	// trees 代表的是森林， HTTP method 代表树的根节点
	// 这个根节点的key是HTTP method, value是path="/"的根节点
	// 用不用额外建一个节点考虑存在这个HTTP method但是完全没有path，连"/"都没有的情况？不用。默认只要有这个HTTP method就存在path "/"
	trees map[string]*node
	// go的这种把类型放在后边的语法还是有点怪，因为有时候看代码，类型是更重要的让我区分它有什么作用，命名只是个名字不太重要
	// 但是我理智上还是知道类型在语法上是不必要可以省略的，放在后边是合理的
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {

	// 怎么报错的同时让测试通过呢？
	if path == "" {
		panic("web: 路由是空字符串")
	}

	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}

	root, ok := r.trees[method]
	if !ok {
		//如果key是method的这棵树不存在，则创建有根节点的这颗树
		//创建一个node然后把这个node的指针赋值回去
		root = &node{path: "/"}
		r.trees[method] = root
	}

	// 根节点在前面创建method的那颗树的时候创建过了，但是没有handleFunc，这里补一下handleFunc后边跳过
	if path == "/" {
		root.handler = handleFunc
		return
	}

	// 支持多级路由

	// Split Path
	// segs := strings.Split(path, "/")
	// 问题1：path="/" 时会拿到一个 "" 空字符串，空字符串不应该进到下一步
	// 问题1解决方法可以在前面特殊处理一下
	// 问题2：path="/user" 时会拿到两个结果，一个"", 一个"user"
	// 问题2解决方法可以去掉 Path 里前导的这个 "/"
	segs := strings.Split(path[1:], "/")

	current := root
	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, a//b 之类的路由"))
		}
		//如果这个node还没有children，把children的map创建出来
		if current.children == nil {
			current.children = make(map[string]*node)
		}
		//然后还得看一下底下的child是不是已经存在，不存在再建，要不然就给它覆盖了
		child, ok := current.children[seg]
		if !ok {
			// 创建path那一段的child
			child = &node{path: seg}
			// tree.children[seg] = newNode
			// 上边这样写的话，这个tree一直是根节点的那个指针，child一直创建在根节点上面。应该有一种什么方法让这个节点递进下去
			current.children[seg] = child
		}
		current = child
	}
	// 遍历到最后一个节点出来再加handleFunc
	current.handler = handleFunc
	root.route = path // todo: 没看懂这行什么意思

}

// 路由查找 参数：method，path，返回：node节点，标记有没有找到
// 为什么返回的是node节点而不是handleFunc？
func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return root, true
	}
	current := root
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if current.children == nil {
			return nil, false
		}
		child, ok := current.children[seg]
		if !ok {
			return nil, false
		}
		current = child
	}
	if current.handler == nil {
		return nil, false
	}
	return current, true
}

// 树的定义

type node struct {
	// 我这一段。比如 /a/b/c 中的 b 这一段
	path string
	// 为了查找到路由之后记录下找到的路径，选择了在node记下route //todo: 这个选择还不是太理解
	route string
	// 命中了我这个节点我就执行你这个方法
	handler HandleFunc
	// children是一种递归，但是tree的节点的定义就是一种递归
	// children也用map？不用slice？每个children的key是什么？path的那一段？不会重复吗？
	// 可以根据path的那一段直接找到哪个leaf是吧，slice用0123定位就没法直接根据path找leaf，得遍历一遍leaf的path才行
	children map[string]*node
}
