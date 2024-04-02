package app

var (
	Success         = newResultCode(0, "请求成功")
	ServerErrors    = newResultCode(100500, "请求处理失败")
	InvalidParams   = newResultCode(100400, "入参错误")
	NotFound        = newResultCode(100404, "资源信息不存在")
	TooManyRequests = newResultCode(100429, "http请求过多")
	LackParams      = newResultCode(100480, "入参不全")
	Unauthorized    = newResultCode(100403, "未经授权")
	RequestTimeOut  = newResultCode(100408, "客户端请求超时")

	// UnKnowError are used for user-defined error messages returned with http status of 500
	UnKnowError          = newResultCode(200000, "未知的错误")
	AuthenticationFailed = newResultCode(100401, "认证失败")
)
