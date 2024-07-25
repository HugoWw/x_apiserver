package apiserver

import (
	"database/sql"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/options"
	apiserverHttp "github.com/HugoWw/x_apiserver/pkg/apiserver/http"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/resources"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/storage"
	rst "github.com/HugoWw/x_apiserver/pkg/client/httpclient"
	"github.com/HugoWw/x_apiserver/pkg/constant"
	restful "github.com/emicklei/go-restful/v3"
)

func CreateBaseResourceServerConfig(c *restful.Container, option *options.ServerRunOptions) (*BaseResourceServerConfig, error) {
	genericServerCfg := &GenericServerConfig{
		Containers:  c,
		FilterChain: nil,
	}

	client := apiserverHttp.InitHttpClient(option.HttpClient)
	resetClient, err := rst.NewHttpClient(client, option.HttpClient.TimeOut, option.CtrlConf.EndPoint, option.HttpClient.MaxRetries)
	if err != nil {
		return nil, err
	}

	sqlDb, err := storage.InitMysqlDBObj(option.MysqlClient)
	if err != nil {
		return nil, err
	}

	gws, err := resources.Default.GetAPIRegisterResource(constant.BaseResource)
	if err != nil {
		return nil, err
	}

	return &BaseResourceServerConfig{
		GenericServerConfig: genericServerCfg,
		webservice:          gws.GenericWebService(),
		httpC:               resetClient,
		Db:                  sqlDb,
	}, nil

}

type BaseResourceServerConfig struct {
	*GenericServerConfig
	webservice *restful.WebService
	httpC      *rst.RTClient
	Db         *sql.DB
	//todo other obj config init
}

func (c *BaseResourceServerConfig) New() *BaseResourceServer {
	if c == nil {
		return nil
	}

	genericServer := c.GenericServerConfig.New()
	baseRsServer := &BaseResourceServer{
		GenericServer: genericServer,
		RestOption: &resources.RestOption{
			HttpClient: c.httpC,
			Db:         c.Db,
		},
	}

	inject, _ := resources.Default.GetInjectObj(constant.BaseResource)
	inject.Inject(baseRsServer.RestOption)

	baseRsServer.installRESTAPI(c.webservice)
	return baseRsServer
}

type BaseResourceServer struct {
	*GenericServer
	RestOption *resources.RestOption
}

func (s *BaseResourceServer) installRESTAPI(ws *restful.WebService) {
	s.GenericServer.addRestAPI(ws)
}
