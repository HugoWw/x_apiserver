package app

import (
	"fmt"
	"net/http"
)

type resultCode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"result,omitempty"`
}

var codes = map[int]string{}

func newResultCode(code int, msg string) *resultCode {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在,请更换一个", code))
	}

	codes[code] = msg
	return &resultCode{
		Code:    code,
		Message: msg,
	}
}

func (r *resultCode) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息：%s", r.GetCode(), r.GetMsg())
}

func (r *resultCode) GetCode() int {
	return r.Code
}

func (r *resultCode) GetMsg() string {
	return r.Message
}

func (r *resultCode) WithErrMsg(msg string) *resultCode {
	newErrCode := *r
	newErrCode.Message = msg
	return &newErrCode
}

func (r *resultCode) Msgf(args []interface{}) string {
	return fmt.Sprintf(r.Message, args...)
}

func (r *resultCode) StatusCode() int {
	switch r.GetCode() {
	case Success.GetCode():
		return http.StatusOK
	case ServerErrors.GetCode():
		return http.StatusInternalServerError
	case InvalidParams.GetCode(), LackParams.GetCode():
		return http.StatusBadRequest
	case Unauthorized.GetCode():
		return http.StatusForbidden
	case TooManyRequests.GetCode():
		return http.StatusTooManyRequests
	case AuthenticationFailed.GetCode():
		return http.StatusUnauthorized
	case NotFound.GetCode():
		return http.StatusNotFound
	case RequestTimeOut.GetCode():
		return http.StatusRequestTimeout
	}

	return http.StatusInternalServerError
}
