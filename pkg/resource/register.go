package resource

import (
	"github.com/HugoWw/x_apiserver/pkg/apiserver/resources"
	"github.com/HugoWw/x_apiserver/pkg/constant"
	"github.com/HugoWw/x_apiserver/pkg/resource/impl"
	"github.com/HugoWw/x_apiserver/pkg/resource/router"
	"github.com/emicklei/go-restful/v3"
)

func init() {
	registerAndInject()
}

func registerAndInject() {
	resources.Default.AddAPIRegisterResource(constant.BaseResource, resources.GenericWebSvcFunc(genericWebSvc))
	resources.Default.AddInjectResourceObj(constant.BaseResource, resources.InjectFunc(impl.InjectFunc))
}

func genericWebSvc() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(constant.BaseResourceAPI).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) // you can specify this per route as well

	// register api router
	(&router.Auth{ws}).Install()
	(&router.Logs{ws}).Install()

	return ws
}
