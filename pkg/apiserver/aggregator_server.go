package apiserver

import (
	"context"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/options"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/health"
	"github.com/HugoWw/x_apiserver/pkg/apiserver/servermux"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	_ "github.com/HugoWw/x_apiserver/pkg/resource"
	"net/http"
	"time"
)

type AggregatorServer struct {
	baseResource *BaseResourceServer
	handler      *servermux.APIServerHandler
	cfgOpt       *options.ServerRunOptions
}

func Create(opt *options.ServerRunOptions) (*AggregatorServer, error) {
	c := servermux.NewAPIServerHandler("APIContainer")

	baseResourceC, err := CreateBaseResourceServerConfig(c.GoRestfulContainer, opt)
	if err != nil {
		return nil, err
	}
	baseResourceS := baseResourceC.New()

	return &AggregatorServer{
		baseResource: baseResourceS,
		cfgOpt:       opt,
		handler:      c,
	}, nil
}

func (a *AggregatorServer) PrepareRun() error {

	a.addHealthCheck()

	// todo other prepare

	return nil
}

func (a *AggregatorServer) Run(stop <-chan struct{}) error {
	clog.Logger.Infof("Start X-ApiServer UUID is: %s", a.cfgOpt.LeaderElection.ID)

	server := http.Server{
		Addr:           a.cfgOpt.Server.BindAddr,
		Handler:        a.handler,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    time.Duration(a.cfgOpt.Server.ServerIdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		<-stop
		clog.Logger.Error("Shutting down server get signal info..........")
		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *AggregatorServer) addHealthCheck() {
	a.handler.NoGoRestfulMux.Handle(health.Liveness.Name(), health.Liveness)
}
