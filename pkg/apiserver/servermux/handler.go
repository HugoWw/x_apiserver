package servermux

import (
	"github.com/HugoWw/x_apiserver/pkg/clog"
	"github.com/HugoWw/x_apiserver/pkg/constant"
	"github.com/emicklei/go-restful/v3"
	"net/http"
	"strings"
)

type APIServerHandler struct {
	GoRestfulContainer *restful.Container
	NoGoRestfulMux     *PathRecorderMux
}

func NewAPIServerHandler(name string) *APIServerHandler {
	nonGoRestfulMux := NewPathRecorderMux(name)

	gorestfulContainer := restful.NewContainer()
	gorestfulContainer.ServeMux = http.NewServeMux()
	gorestfulContainer.Router(restful.CurlyRouter{})

	// Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders: []string{
			"Content-Type", constant.X_AUTH_TOKEN,
			"Access-Control-Allow-Origin", "Access-Control-Allow-Headers",
		},
		AllowedHeaders: []string{"Content-Type", "Accept", "Content-Length", "Accept-Encoding"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		CookiesAllowed: true,
		AllowedDomains: []string{".*"},
		Container:      gorestfulContainer,
	}

	gorestfulContainer.Filter(cors.Filter)

	return &APIServerHandler{
		GoRestfulContainer: gorestfulContainer,
		NoGoRestfulMux:     nonGoRestfulMux,
	}

}

func (a *APIServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	//clog.Logger.Infof("request url: %s %q", r.Method, r.URL.String())

	for _, ws := range a.GoRestfulContainer.RegisteredWebServices() {
		switch {
		case strings.HasPrefix(path, ws.RootPath()):
			clog.Logger.Infof("%v: %v %q satisfied by gorestful with webservice %v", a.NoGoRestfulMux.name, r.Method, path, ws.RootPath())
			a.GoRestfulContainer.Dispatch(w, r)
			return

		}
	}

	clog.Logger.Infof("%v: %v %q satisfied by nonGoRestful", a.NoGoRestfulMux.name, r.Method, path)
	a.NoGoRestfulMux.ServeHTTP(w, r)
}
