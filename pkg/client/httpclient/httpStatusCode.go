package httpclient

import (
	"github.com/HugoWw/x_apiserver/pkg/client/httpclient/cerrors"
	"net/http"
)

func httpStatusCode2String(code int, method string, host string, body *[]byte) *cerrors.ResultStatus {

	reason := ""

	switch code {
	case http.StatusConflict:
		if method == "POST" {
			reason = "maybe resource already exists,the server reported a conflict."
		} else {
			reason = "the server reported a conflict."
		}
	case http.StatusNotFound:
		reason = "the server could not find the requested resource."
	case http.StatusBadRequest:
		reason = "the server rejected our request for an unknown reason because of BadRequest."
	case http.StatusUnauthorized:
		reason = "the server has asked for the client to provide credentials."
	case http.StatusForbidden:
		reason = "the access is forbidden because the access is not authorized."
	case http.StatusNotAcceptable:
		reason = "the server was unable to respond with a content type that the client supports."
	case http.StatusUnsupportedMediaType:
		reason = "the server unsupported media type."
	case http.StatusMethodNotAllowed:
		reason = "the server does not allow this method on the requested resource."
	case http.StatusUnprocessableEntity:
		reason = "the server rejected our request due to an error in our request."
	case http.StatusServiceUnavailable:
		reason = "the server is currently unable to handle the request."
	case http.StatusGatewayTimeout:
		reason = "the server was unable to return a response in the time allotted."
	case http.StatusTooManyRequests:
		reason = "the server has received too many requests and has asked us to try again later."
	case http.StatusRequestTimeout:
		reason = "client time out,password maybe expired or request body sending times out."
	default:
		if code >= 500 {
			reason = "an error on the server (unknown) has prevented the request from succeeding"
		}
	}

	return &cerrors.ResultStatus{ErrorStatus: cerrors.APIStatus{
		HttpCode:   code,
		HttpMethod: method,
		Body:       body,
		Host:       host,
		Reason:     reason,
	}}
}
