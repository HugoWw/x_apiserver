package router

import (
	"github.com/HugoWw/x_apiserver/pkg/resource/impl/auth"
	"github.com/HugoWw/x_apiserver/pkg/resource/v1"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

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

	a.Ws.Route(a.Ws.DELETE("/auth").To(auth.LoginOut).
		Doc("登出").
		Param(a.Ws.HeaderParameter("token", "token info").DataType("string").Required(true)).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "ok", v1.APIResponse[string]{}),
	)

	a.Ws.Route(a.Ws.PATCH("/auth").To(auth.RefreshToken).
		AllowedMethodsWithoutContentType([]string{"PATCH"}).
		Doc("刷新token").
		Param(a.Ws.HeaderParameter("token", "token info").DataType("string").Required(true)).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "ok", v1.APIResponse[string]{}),
	)
}
