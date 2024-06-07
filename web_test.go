package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	//s.Use(accesslog.NewBuilder().Build())
	// 这个要注掉，不然就循环引用了。middleware是上层包，依赖底层的web包，测试middleware的测试需要写在上层的middleware包里
	// 详见 Go语言圣经 11.2.4 外部测试包
	s.AddRoute(http.MethodGet, "/", func(ctx *Context) {
		ctx.RespData = []byte("hello, word")
	})
	// 这里这个函数我没有执行，只是一个定义，把这个函数存储到了node.handler,
	// 这个node.handler的类型是一个HandleFunc，是一个函数类型，函数签名，它里面存储的东西是这个函数代码，等待什么时候被调用
	// 一个变量里面存储的是一段代码
	// 一个变量里面存储的是一段代码？
	// 那你编译的时候能找到所有函数签名类型的变量，然后找到他们所存储的值吗？好像不行吧？这段代码是不是要运行的时候才编译的啊？
	s.AddRoute(http.MethodGet, "/user/123/info", func(ctx *Context) {
		println(ctx.MatchedRoute)
		println(ctx.Req.URL.Path)
		ctx.RespData = []byte("hello, user")
	})
	s.AddRoute(http.MethodPost, "/user", func(ctx *Context) {
		// 创建接口的人要去定义你给我传一个什么对象
		u := &User{}
		err := ctx.BindJSON(u)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(u)
		stringu, _ := json.Marshal(u)
		ctx.RespData = []byte("POST user: " + string(stringu))
	})
	s.Start(":8080")
}

type User struct {
	Name string `json:"name"`
}
