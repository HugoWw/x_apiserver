package impl

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/resources"
	"github.com/HugoWw/x_apiserver/pkg/client/httpclient"
	"github.com/HugoWw/x_apiserver/pkg/dao/ent"
)

var (
	HttpC *httpclient.RTClient
	DBC   *ent.Client
)

func InjectFunc(opt *resources.RestOption) {
	driver := entsql.OpenDB("mysql", opt.Db)
	DBC = ent.NewClient(ent.Driver(driver))
	HttpC = opt.HttpClient
}
