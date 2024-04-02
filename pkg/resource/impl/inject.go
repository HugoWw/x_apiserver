package impl

import (
	"github.com/HugoWw/x_apiserver/pkg/apiserver/resources"
	"github.com/HugoWw/x_apiserver/pkg/client/httpclient"
)

var HttpC *httpclient.RTClient

func InjectFunc(opt *resources.RestOption) {
	HttpC = opt.HttpClient
}
