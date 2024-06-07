# minimal-web-framework-go

core:
* web server 
* router 
* context 
* middleware


middleware

就是一个一个函数链？前置函数后置函数？包装函数？

但是怎么感觉根Spring里讲的AOP不太一样？Spring是怎么讲的？


在HTTPServer里增加一个middleware数组
添加一个Use方法为HTTPServer里增加middleware


todo：后边还能再扩展accesslog，比如不用log函数让用户传入log方法，使用defer最后打印日志（不知道defer到哪了），添加打印匹配路径

为了扩展accesslog在accesslog中打印MatchedRoute，在Context中，查找到路由之后记下MatchedRoute。// todo: MatchedRoute功能还没写完
查找到路由之后怎么记录下找到的路径呢？选择了在路由节点中加入route //todo: 这个选择还不是太理解

加入对responseData的缓存，因为否则http的resp就之间返回给前端了

todo：具体看第七章第九讲2小时30分钟之后的部分

todo：添加middleware tracing, Prometheus. 第十讲1小时52分之前
todo：基于路由的middleware，在第十讲最后（没讲太多，让作为作业自己写）


服务端支持的渲染错误页面，先后端分离是用不上这种的，只有后端页面渲染才需要
用户可以自己在路由的方法里面渲染页面，把渲染好的页面写回Response
设计成核心功能的服务端页面渲染：第十讲2小时24分开始，到2小时45分结束


todo：看文件处理， 第十讲2小时45分之后文件上传，第十一讲35分之后文件下载，59分之后静态资源处理
这里额外考虑一个问题，怎么缓存到内存，避免每次都从磁盘加载
静态资源和下载不一样的是，你需要告诉浏览器具体的content-type，不然浏览器不知道怎么处理

todo：第十一讲1小时36分之后开始讲session

