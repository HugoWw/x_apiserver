package app

import (
	"fmt"
	"github.com/HugoWw/x_apiserver/pkg/util"
	"github.com/emicklei/go-restful/v3"
)

type WrapResponse struct {
	respWrite *restful.Response
}

func NewResponse(resp *restful.Response) *WrapResponse {
	return &WrapResponse{resp}
}

func (r *WrapResponse) Response(data interface{}) {

	if v, ok := data.(*resultCode); ok {
		r.respWrite.WriteHeaderAndJson(v.StatusCode(), v, restful.MIME_JSON)
	} else {
		var successCode resultCode
		successCode = *Success
		successCode.Data = data

		r.respWrite.WriteHeaderAndJson(successCode.StatusCode(), successCode, restful.MIME_JSON)
	}
}

const (
	defaultPage     = "0"
	defaultPageSize = "10"
)

// GetPaginationParam get `page` and `pageSize` from query
func GetPaginationParam(pageStr, pageSizeStr string, minPageSize, maxPageSize int) (int, int, error) {
	if pageStr == "" {
		pageStr = defaultPage
	}

	if pageSizeStr == "" {
		pageSizeStr = defaultPageSize
	}

	page, err := util.StrTo(pageStr).Int()
	if err != nil {
		return 0, 0, fmt.Errorf("invalid page %s: %v", pageStr, err)
	}

	pageSize, err := util.StrTo(pageSizeStr).Int()
	if err != nil {
		return 0, 0, fmt.Errorf("invalid pageSize %s: %v", pageStr, err)
	}

	if page < 0 {
		page = 0
	}

	if pageSize < minPageSize {
		pageSize = maxPageSize
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize, nil
}
