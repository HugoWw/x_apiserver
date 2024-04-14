package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net/url"
)

type CtrlOptions struct {
	EndPoint string
}

func NewCtlOptions() *CtrlOptions {
	return &CtrlOptions{
		EndPoint: "http://auth-example-server.default:10433",
	}
}

func (o *CtrlOptions) Valid() []error {

	errors := []error{}

	hostURL, err := url.Parse(o.EndPoint)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		hostURL, err = url.Parse("http://" + o.EndPoint)
		if err != nil {
			errors = append(errors, fmt.Errorf("--ctrl-conf %v The host endpoint url must be in normal format http://host", o.EndPoint))
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (o *CtrlOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.EndPoint, "ctrl-conf", o.EndPoint, "The controller endpoint address.")
}
