package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

type ServerOptions struct {
	BindAddr          string
	ServerIdleTimeout int
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		BindAddr:          "0.0.0.0:8866",
		ServerIdleTimeout: 90,
	}
}

func (o *ServerOptions) Valid() []error {
	errors := []error{}

	host, _, err := net.SplitHostPort(o.BindAddr)
	if err != nil {
		errors = append(errors, fmt.Errorf("--bind-addr %v bind address format should be host:port", o.BindAddr))
	}

	if validIP := net.ParseIP(host); validIP == nil {
		errors = append(errors, fmt.Errorf("--bind-addr %v the bound ip address does not comply with ipv4 specifications", o.BindAddr))
	}

	return errors
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.BindAddr, "bind-addr", o.BindAddr, "The bind address used to serve the http APIs.")
	fs.IntVar(&o.ServerIdleTimeout, "server-idleTimeout", o.ServerIdleTimeout, "Http Server IdleTimeout is the maximum amount of time to wait "+
		"for the next request when keep-alives are enabled. If IdleTimeout is zero, the value of ReadTimeout is used. "+
		"If both are zero, there is no timeout.")
}
