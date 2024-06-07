package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

//"github.com/stretchr/testify/assert"

func Test_router_addRoute(t *testing.T) {

	tests := []struct {
		// 输入
		method string
		path   string
	}{
		// 静态匹配
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		// 加底下这俩就要panic报错了
		//{
		//	method: http.MethodGet,
		//	path:   "//home",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "//home1///",
		//},
		{
			method: http.MethodGet,
			path:   "/user/detail/profile",
		},
		{
			method: http.MethodGet,
			path:   "/order/cancel",
		},

		{
			method: http.MethodPost,
			path:   "/order/cancel",
		},
		// 这个乱写方法可以去掉，因为后边会讲，不会把这个方法暴露出去
		//{
		//	method: "乱写方法",
		//	path:   "/",
		//},
	}

	var handleFunc HandleFunc = func(context *Context) {

	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: handleFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: handleFunc,
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								children: map[string]*node{
									"profile": &node{
										path:    "profile",
										handler: handleFunc,
									},
								},
							},
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: handleFunc,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"cancel": &node{
								path:    "cancel",
								handler: handleFunc,
							},
						},
					},
				},
			},
			// 这个乱写方法可以去掉，因为后边会讲，不会把这个方法暴露出去
			//"乱写方法": &node{
			//	path:    "/",
			//	handler: handleFunc,
			//},
		},
	}

	res := &router{
		trees: map[string]*node{},
	}

	for _, tc := range tests {
		res.addRoute(tc.method, tc.path, handleFunc)

	}

	// 断言两棵树相等
	//assert.Equal(t, wantRouter, res)
	errStr, ok := wantRouter.equal(res)
	assert.True(t, ok, errStr)

	//暂时把findRoute的测试写在这

	findCases := []struct {
		name     string
		method   string
		path     string
		found    bool
		wantPath string
	}{
		{
			name:     "/",
			method:   http.MethodGet,
			path:     "/",
			found:    true,
			wantPath: "/",
		},
		{
			name:     "/user",
			method:   http.MethodGet,
			path:     "/user",
			found:    true,
			wantPath: "user",
		},
		{
			name:   "/user/detail",
			method: http.MethodGet,
			path:   "/user/detail",
			found:  false,
		},
	}

	for _, tc := range findCases {
		t.Run(tc.name, func(t *testing.T) {
			n, ok := res.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.found, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantPath, n.path)
			//应该assert一下两个handleFunc一不一样，但是因为不能直接比，就只比是不是nil了
			assert.NotNil(t, n.handler)
		})
	}

}

func (r router) equal(y *router) (string, bool) {
	for key, value := range r.trees {
		yvalue, ok := y.trees[key]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有 %s 的路由树", key), false
		}
		str, ok := value.equal(yvalue)
		if !ok {
			return key + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}
	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type(), yhv.Type()), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不相等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}
	for k, v := range y.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

/*


右键 -> Generate... -> Test for function 来创建测试

*/

func Test_router_findRoute(t *testing.T) {
	type fields struct {
		trees map[string]*node
	}
	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *node
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router{
				trees: tt.fields.trees,
			}
			got, got1 := r.findRoute(tt.args.method, tt.args.path)
			assert.Equalf(t, tt.want, got, "findRoute(%v, %v)", tt.args.method, tt.args.path)
			assert.Equalf(t, tt.want1, got1, "findRoute(%v, %v)", tt.args.method, tt.args.path)
		})
	}
}
