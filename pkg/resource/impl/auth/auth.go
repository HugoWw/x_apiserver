package auth

import (
	"context"
	"encoding/json"
	"github.com/HugoWw/x_apiserver/pkg/app"
	"github.com/HugoWw/x_apiserver/pkg/client/httpclient/cerrors"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	"github.com/HugoWw/x_apiserver/pkg/constant"
	"github.com/HugoWw/x_apiserver/pkg/resource/impl"
	v1 "github.com/HugoWw/x_apiserver/pkg/resource/v1"
	"github.com/emicklei/go-restful/v3"
)

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

func LoginOut(request *restful.Request, response *restful.Response) {
	resp := app.NewResponse(response)
	token := request.HeaderParameter("Token")

	if token == "" {
		clog.Logger.Errorf("The token in the request header is empty in %s", request.Request.URL.Path)
		resp.Response(app.LackParams.WithErrMsg("请求头中缺少token"))
		return
	}

	err := impl.HttpC.Delete().
		SetPath("/v1/auth").
		SetHeader(constant.X_AUTH_TOKEN, token).
		Body(nil).
		Do(context.TODO()).
		Error()

	if err != nil {
		clog.Logger.Errorf("Delete(login out)  Controller /v1/auth failed %v", err)
		resp.Response(app.UnKnowError.WithCodeAndMsg(cerrors.ErrorCodeForResponse(err), cerrors.ResponseForErrorReason(err)))
		return
	}

	resp.Response(app.Success)
}
