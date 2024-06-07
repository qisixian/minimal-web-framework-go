package web

import (
	"log"
	"net"
	"net/http"
)

// 定义函数类型？
type HandleFunc func(*Context)

type Server interface {
	http.Handler
	Start(addr string) error
	// AddRoute 注册路由的核心抽象
	AddRoute(method, path string, handler HandleFunc)
}

type HTTPServer struct {
	//不加名字吗？
	router
	mlds []Middleware
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{router: newRouter()}
}

func (s *HTTPServer) Use(mlds ...Middleware) {
	// 参数这样写的话，ms自动是一个切片是吗？
	if s.mlds == nil {
		s.mlds = mlds
		return
	}
	s.mlds = append(s.mlds, mlds...)
}

func (s *HTTPServer) AddRoute(method, path string, handler HandleFunc) {
	//直接穿透到HTTPServer里到router到addRoute
	//这样的话让组合/聚合好像继承啊
	//但是有一点是，这里在HTTPServer这个结构体里又重复定义了一遍AddRoute这个方法，而不是直接用的router里的addRoute
	//这么重复定义是为了公开函数的访问权限，go语言中函数小写开头是private访问权限，大写开头是public访问权限
	s.addRoute(method, path, handler)
}

func (s *HTTPServer) Start(addr string) error {
	// 端口启动前
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	// 端口启动后
	// 注册本服务器到你的管理平台
	// 比如说你注册到 etcd，然后你打开管理界面，你就能看到这个实例
	println("成功监听端口 " + addr)

	return http.Serve(listener, s)
	//return http.ListenAndServe(addr, m)
	//ListenAndServe是阻塞的，如果想启动后做点事情要用Serve, 并且把listener创建出来
	//启动后可以做点事情
}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 这个ServeHTTP方法在ListenAndServe之后，是自动被调用的是吗？
	//writer.Write([]byte("hello world"))
	ctx := &Context{
		Req:            request,
		Resp:           writer,
		RespStatusCode: 200,
	}
	// 构造/组装了一个调用链？
	// 为什么不能从前往后？因为你是用next构造你的previous？
	// 调到最后一个肯定是serve？
	// 最后一个就是查找路由，执行你的业务逻辑，然后返回
	root := s.serve //这个是一个函数类型
	for i := len(s.mlds) - 1; i >= 0; i-- {
		m := s.mlds[i] // 最后一个Middleware类型的函数
		root = m(root) // 把root函数传进最后一个Middleware类型的函数的next参数里 //这样传参的时候调用没调用函数？还是只是声明了参数？
	}
	// 这里的 root(ctx) 和前面的 m(root) 有什么区别？
	// root(ctx)肯定是执行了函数对吧，那 m(root) 呢？
	// root(ctx) ctx传入了一个Context类型的参数值
	// m(root) 传入了一个函数类型的参数值
	// 哦！因为Middleware类型的返回值是一个函数，所以第一次调用会返回一个函数，第二次调用才是真正执行？
	root(ctx)
	s.writeResp(ctx)
}

// serve也可以看做是一个handleFunc类型的函数
func (s *HTTPServer) serve(ctx *Context) {
	n, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.RespStatusCode = http.StatusNotFound
		ctx.RespData = []byte("404 NotFound")
		return
	}
	ctx.MatchedRoute = n.route
	n.handler(ctx)
}

// 因为在context把response缓存住了，所以在调用流程中需要加一步把resp刷回去
func (s *HTTPServer) writeResp(ctx *Context) {
	ctx.Resp.WriteHeader(ctx.RespStatusCode)
	ctx.Resp.Write(ctx.RespData)
}
