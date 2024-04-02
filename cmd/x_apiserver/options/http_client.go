package options

import "github.com/spf13/pflag"

type HttpClientOption struct {
	TimeOut         int
	MaxIdleConns    int
	IdleConnTimeout int
	MaxRetries      int
}

func NewHttpClientOptions() *HttpClientOption {
	return &HttpClientOption{
		TimeOut:         30,
		MaxIdleConns:    100,
		IdleConnTimeout: 90,
		MaxRetries:      5,
	}
}

func (h *HttpClientOption) Valid() []error {
	return nil
}

func (h *HttpClientOption) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&h.TimeOut, "http-client-timeout", h.TimeOut, "The Http client timeout.")

	fs.IntVar(&h.MaxIdleConns, "http-client-max-idle-conn", h.MaxIdleConns, "The Http client MaxIdleConns "+
		"controls the maximum number of idle keep-alive connections across all hosts. Zero means no limit.")

	fs.IntVar(&h.IdleConnTimeout, "http-client-idle-conn-timeout", h.IdleConnTimeout, "The Http client IdleConnTimeout "+
		"is the maximum amount of time an idle keep-alive connection will remain idle before closing itself. "+
		"Zero means no limit.")

	fs.IntVar(&h.MaxRetries, "http-client-maxretries", h.MaxRetries, "The Http client max retries")
}
