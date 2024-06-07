package accesslog

import (
	web "minimal-web-framework-go"
	"net/http"
	"testing"
)

func TestAccesslog(t *testing.T) {
	s := web.NewHTTPServer()
	s.Use(NewBuilder().Build())
	s.AddRoute(http.MethodGet, "/", func(ctx *web.Context) {
		ctx.RespData = []byte("hello, world")
	})
	s.Start(":8080")
}
