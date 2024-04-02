package apiserver

import restful "github.com/emicklei/go-restful/v3"

type GenericServerConfig struct {
	Containers  *restful.Container
	FilterChain []restful.FilterFunction
}

func (g *GenericServerConfig) New() *GenericServer {
	return &GenericServer{
		Containers:  g.Containers,
		FilterChain: g.FilterChain,
	}
}

type GenericServer struct {
	Containers  *restful.Container
	FilterChain []restful.FilterFunction
}

func (s *GenericServer) addRestAPI(ws *restful.WebService) {
	if s.FilterChain != nil {
		for _, f := range s.FilterChain {
			ws.Filter(f)
		}
	}

	s.Containers.Add(ws)
}
