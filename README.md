

# 简介
这是一个根据Kube-APIServer简化的项目，可以基于该简化的框架快速的搭建自己的web api服务.

----

# 功能
- 运行时动态修改平台日志级别(info｜debug)
- 自定义不同模块日志分类输出和级别定义
- 简化swag定义生成和对应接口的维护
- apiserver服务自身健康检查
- restful格式的http client请求

# 整体目录结构
该项目的大致目录结构如下：
```
./
├── README.md
├── build
│   ├── Dockerfile.swag
│   ├── apiserver
│   ├── apiserver.yaml
│   ├── build.sh
│   └── swagger-gen.go
├── cmd
│   └── x_apiserver
├── go.mod
├── go.sum
└── pkg
    ├── apiserver #--x-apiserver启动的核心框架
    ├── app  #---响应体封装
    ├── client #--各种客户端(http、kube client等等)
    ├── clog #---zap日志的封装
    ├── constant #--常量
    ├── middleware
    ├── permissions #--权限相关
    ├── resource  #--所有的路由、实现、结构体定义
    ├── signals   #--进程信号管理
    └── util
```

# 快速上手
这里以平台登陆认证对接其它认证服务器转发接口为例
1. 定义好api结构体
x_apiserver/pkg/resource/v1
```go
// 请求体内容
type AuthUserLoginReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// 响应体内容
type AuthData struct {
    Roles  map[string]string `json:"roles"`
    Token  string            `json:"token"`
    Status AuthStatus        `json:"status"`
}
```

2. 实现对应的接口的handler
x_apiserver/pkg/resource/impl/auth/auth.go
```go
func Login(request *restful.Request, response *restful.Response) {
	resp := app.NewResponse(response)
	token := request.HeaderParameter("Token")

	params := v1.AuthUserLoginReq{}

	if err := request.ReadEntity(&params); err != nil {
		clog.Logger.Errorf("Failed to parse the request parameters %v", err)
		resp.Response(app.InvalidParams)
		return
	}

	restReq := v1.RestAuthData{
		ClientIP: request.Request.RemoteAddr,
		Password: &v1.RestAuthPassword{
			Username: params.Username,
			Password: params.Password,
		},
		Token: nil,
	}

	data, err := json.Marshal(restReq)
	if err != nil {
		clog.Logger.Errorf("user login json marshal error:%v", err)
		resp.Response(app.ServerErrors.WithErrMsg("user login json marshal error"))
		return
	}

	restRes := v1.AuthData{}
	err = impl.HttpC.Post().
		SetPath("/v1/auth").
		SetHeader(constant.X_AUTH_TOKEN, token).
		Body(data).
		Do(context.TODO()).
		Into(&restRes)

	if err != nil {
		clog.Logger.Errorf("Post Controller /v1/auth failed %v", err)
		resp.Response(app.AuthenticationFailed.WithErrMsg(cerrors.ResponseForErrorReason(err)))
		return
	}

	resp.Response(restRes)
}

```

3. 生成api swag文档和注册api
x_apiserver/pkg/resource/route/auth.go
生成api接口swag文档具体内容
```go
type Auth struct {
	Ws *restful.WebService
}

func (a *Auth) Install() {
	tags := []string{"认证"}

	a.Ws.Route(a.Ws.POST("/auth").To(auth.Login).
		Doc("登陆").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.AuthUserLoginReq{}).
		Returns(200, "ok", v1.APIResponse[v1.AuthData]{}),
	)
}
```
注册api接口
x-apiserver/pkg/resource/register.go
```go
(&router.Auth{ws}).Install()
```

-----

# 项目构建
使用build.sh脚本来构建项目和swag文档。同时可以基于该脚本接入ci/cd完成自动的构建
```shell
./build.sh --help      
Execute build.sh script build programs.

usage: ./build.sh [OPTIONS]

The following flags are required.

    --build-apiserver           Only build apiserver program
    --build-apiserver-img  tag  Build apiserver image(if 'tag' is empty,the default is latest)
```

-----