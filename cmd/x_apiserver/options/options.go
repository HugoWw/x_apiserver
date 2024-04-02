package options

import (
	"github.com/HugoWw/x_apiserver/pkg/apiserver/cli/flag"
)

// ServerRunOptions runs a api server options.
type ServerRunOptions struct {
	Server         *ServerOptions
	HttpClient     *HttpClientOption
	LuXuCtrl       *LuXuCtrlOptions
	LeaderElection *LeaderElectionOptions
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	s := ServerRunOptions{
		Server:         NewServerOptions(),
		HttpClient:     NewHttpClientOptions(),
		LuXuCtrl:       NewLuXuCtlOptions(),
		LeaderElection: NewLeaderElectionOptions(),
	}

	return &s
}

// Validate checks ServerRunOptions and return a slice of found errs.
func (s *ServerRunOptions) Validate() []error {
	errs := []error{}
	errs = append(errs, s.Server.Valid()...)
	errs = append(errs, s.HttpClient.Valid()...)
	errs = append(errs, s.LuXuCtrl.Valid()...)
	errs = append(errs, s.LeaderElection.Valid()...)

	return errs
}

// Flags returns flags for a specific APIServer by section name
func (s *ServerRunOptions) Flags() (fss flag.NamedFlagSets) {

	s.Server.AddFlags(fss.FlagSet("server"))
	s.HttpClient.AddFlags(fss.FlagSet("http-client"))
	s.LuXuCtrl.AddFlags(fss.FlagSet("luxu-controller"))
	s.LeaderElection.AddFlags(fss.FlagSet("leader-election"))

	return fss
}
