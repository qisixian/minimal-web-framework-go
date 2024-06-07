package accesslog

import (
	"encoding/json"
	"log"
	web "minimal-web-framework-go"
)

type MiddlewareBuilder struct {
	// 为什么要加这个？把函数作为一个变量，这样外部就可以注入这个函数自己的实现了是吗？
	logFunc func(accesslog string)
}

func (b *MiddlewareBuilder) InsertLogFunc(logFunc func(accessLog string)) *MiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

// 这个方法用来初始化一个 MiddlewareBuilder
func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

// 这个方法会返回一个middleware
func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {

			defer func() {
				l := accesslog{
					Host:           ctx.Req.Host,
					Route:          ctx.MatchedRoute,
					Path:           ctx.Req.URL.Path,
					HTTPMethod:     ctx.Req.Method,
					RespStatusCode: ctx.RespStatusCode,
				}
				val, _ := json.Marshal(l)
				m.logFunc(string(val))
			}()
			//archive// log.Println(io.ReadAll(ctx.Req.Body)) // 问题：只能读一遍，读完之后就访问不到了
			next(ctx)
			//archive//log.Println(ctx.Resp) // 问题：写完之后是读不出来的
		}
	}
}

type accesslog struct {
	Host           string
	Route          string
	HTTPMethod     string `json:"http_method"`
	Path           string
	RespStatusCode int
}

//用这种创建Middleware的方法比较好，有有扩展性
