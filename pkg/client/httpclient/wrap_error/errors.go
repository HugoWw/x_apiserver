package wrap_error

import (
	"errors"
	"fmt"
	"net/http"
)

type APIInterfaceError interface {
	Status() APIStatus
}

// ResultStatus is a rest api request error wrapper
// used to wrap normal and abnormal errors for converting http requests status
type ResultStatus struct {
	ErrorStatus APIStatus
}

type APIStatus struct {
	HttpCode   int
	HttpMethod string
	Body       *[]byte
	Host       string
	Reason     string
}

func (a *ResultStatus) Error() string {
	return fmt.Sprintf("http code:%d, method:%s, host:%s, respbody:%s, reson: %s",
		a.ErrorStatus.HttpCode, a.ErrorStatus.HttpMethod, a.ErrorStatus.Host, string(*(a.ErrorStatus.Body)), a.ErrorStatus.Reason)
}

// Status return rest api request status info include error reason
func (a *ResultStatus) Status() APIStatus {
	return a.ErrorStatus
}

// ResponseForErrorReason return a brief message about the api request error
func ResponseForErrorReason(err error) string {
	if status := APIInterfaceError(nil); errors.As(err, &status) {
		return status.Status().Reason
	}

	return err.Error()
}

func ErrorCodeForResponse(err error) int {
	if status := APIInterfaceError(nil); errors.As(err, &status) {
		return status.Status().HttpCode
	}

	return 0
}

func IsNotFound(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusNotFound
}

func IsConflict(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusConflict
}

func IsBadRequest(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusBadRequest
}

func IsUnauthorized(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusUnauthorized
}

func IsForbidden(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusForbidden
}

func IsNotAcceptable(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusNotAcceptable
}

func IsUnsupportedMediaType(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusUnsupportedMediaType
}

func IsMethodNotAllowed(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusMethodNotAllowed
}

func IsUnprocessableEntity(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusUnprocessableEntity
}

func IsServiceUnavailable(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusServiceUnavailable
}

func IsGatewayTimeout(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusGatewayTimeout
}

func IsTooManyRequests(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusTooManyRequests
}

func IsRequestTimeout(err error) bool {
	return ErrorCodeForResponse(err) == http.StatusRequestTimeout
}
