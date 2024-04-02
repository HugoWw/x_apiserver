package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net/url"
)

type LuXuCtrlOptions struct {
	EndPoint string
}

func NewLuXuCtlOptions() *LuXuCtrlOptions {
	return &LuXuCtrlOptions{
		EndPoint: "http://lx-svc-ctrl.default:10433",
	}
}

func (o *LuXuCtrlOptions) Valid() []error {

	errors := []error{}

	hostURL, err := url.Parse(o.EndPoint)
	if err != nil || hostURL.Scheme == "" || hostURL.Host == "" {
		hostURL, err = url.Parse("http://" + o.EndPoint)
		if err != nil {
			errors = append(errors, fmt.Errorf("--lxctrl-endpoint %v The host endpoint url must be in normal format http://host", o.EndPoint))
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (o *LuXuCtrlOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.EndPoint, "lxctrl-endpoint", o.EndPoint, "The luxu controller endpoint address.")
}
