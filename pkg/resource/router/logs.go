package router

import (
	"github.com/HugoWw/x_apiserver/pkg/resource/impl/logs"
	v1 "github.com/HugoWw/x_apiserver/pkg/resource/v1"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

type Logs struct {
	Ws *restful.WebService
}

func (l *Logs) Install() {
	tags := []string{"日志"}

	l.Ws.Route(l.Ws.POST("/debuglog").To(logs.SetDebugLog).
		Doc("开启DEBUG日志").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(v1.DebugLogReq{}).
		Returns(200, "ok", v1.APIResponse[string]{}),
	)
}
