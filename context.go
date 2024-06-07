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

	// 缓存的查询参数
	cacheQueryValues url.Values

	// 查找到路由之后记录下找到的路径
	MatchedRoute string

	// 缓存住你的响应, 因为http库默认的响应你写进去就之间返回给前端了
	RespStatusCode int
	RespData       []byte
}

func (ctx *Context) BindJSON(val any) error {
	if ctx.Req.Body == nil {
		return errors.New("web: body 为空")
	}
	decoder := json.NewDecoder(ctx.Req.Body)
	return decoder.Decode(val)
}

// 查询参数（URL问号之后的部分）
func (c Context) QueryValue(key string) StringValue {
	// params := c.Req.URL.Query()
	// 它返回一个 Values类型，它是一个 map[string][]string
	// 包装能包几层呢，你就往上点，肯定能点到头
	//   所以说如果查询参数里面有两个相同的key，那么两个值都会拿到
	//   那这里要做一个决策，你定义的方法，查询一个key，要返回 string 还是 []string，返回一个还是多个
	//   大多数场景都是一个key只会有一个值
	// 它调用ParseQuery(), ParseQuery再调用parseQuery()
	// 这个parseQuery()是确确实实的每次都会解析一遍查询串，那个rawQuery，它每次都会切割那个字符串去做处理
	// （不像 ParseForm() 只会Parse一次之后会把值存到 Request.Form 和 Request.PostForm 属性里）
	// 那你会想不希望你每次都解析，只解析一次，那就引入查询参数缓存（和Gin一样），在Context里缓存住
	// 这个缓存不会存在失效或者不一致的问题，因为在这个请求内请求是不会再变的
	// 到这一步你就能理解Gin的设计，你自己做到这一步你发现，你为了避免重复解析你就会引入缓存
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok || len(vals) == 0 {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}
}
func (c Context) PathValue(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: val}
}

func (c Context) RespJSON(responseCode int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespStatusCode = responseCode
	c.RespData = bs
	//下面是不缓存的写法
	//c.Resp.WriteHeader(responseCode)
	//要先写WriteHeader因为在Write里有说明：
	//If WriteHeader has not yet been called, Write calls
	//WriteHeader(http.StatusOK) before writing the data.
	//_, err = c.Resp.Write(bs)

	return err
}

func (c Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
	// 就这么封装吗？就这么封装，只是为了新人可能找不到http.SetCookie这个方法。其实你根本可以不提供

}

// 实际中我是不建议大家使用表单的，我一般只在 Body里面用JSON 通信，或者用 protobuf通信
// 表单在Go 的http.Request 里面有两个
// •Form：基本上可以认为，所有的表单数据都能拿到
// •PostForm：在编码是 x-www-form-urlencoded 的时候才能拿到
// 用Form就行，不用考虑PostForm
func (c Context) FormValue(key string) StringValue {
	err := c.Req.ParseForm()
	// 从ParseForm()里可以看到，它只会parse RawQuery一次，之后就会把表单数据存到 Request.Form 和 Request.PostForm 这两个属性里了
	if err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: c.Req.FormValue(key)}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
