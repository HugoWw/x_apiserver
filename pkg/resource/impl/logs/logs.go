package logs

import (
	"fmt"
	"github.com/HugoWw/x_apiserver/pkg/app"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	v1 "github.com/HugoWw/x_apiserver/pkg/resource/v1"
	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"
	"strings"
)

func SetDebugLog(request *restful.Request, response *restful.Response) {
	resp := app.NewResponse(response)

	debugLog := v1.DebugLogReq{}
	err := request.ReadEntity(&debugLog)
	if err != nil {
		clog.Logger.Errorf("request read entity error:%v", err)
		resp.Response(app.InvalidParams)
		return
	}

	if debugLog.DebugModule == "" {
		if err := clog.GlobalClogSets.SetLogAtomicLevel(zap.InfoLevel); err != nil {
			resp.Response(app.ServerErrors.WithErrMsg(fmt.Sprintf("change log level error:%v", err)))
			return
		}
		resp.Response(app.Success)
		return
	}

	isALl := false
	needModuleList := []string{}
	moduleList := strings.Split(debugLog.DebugModule, ",")

	for _, v := range moduleList {
		if v == "all" {
			isALl = true
			continue
		}
		needModuleList = append(needModuleList, v)
	}

	if isALl {
		if err := clog.GlobalClogSets.SetLogAtomicLevel(zap.DebugLevel); err != nil {
			resp.Response(app.ServerErrors.WithErrMsg(fmt.Sprintf("change log level error:%v", err)))
			return
		}
	} else {
		if err := clog.GlobalClogSets.SetLogAtomicLevel(zap.DebugLevel, needModuleList...); err != nil {
			resp.Response(app.ServerErrors.WithErrMsg(fmt.Sprintf("change log level error:%v", err)))
			return
		}
	}

	resp.Response(app.Success)
}
