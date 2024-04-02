package health

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type HealthCheck interface {
	Name() string
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}

var Liveness HealthCheck = &liveness{}

type liveness struct {
	startOnce    sync.Once
	lastVerified atomic.Value
}

func (l *liveness) Name() string {
	return "/healthz"
}

func (l *liveness) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	l.startOnce.Do(func() {
		l.lastVerified.Store(time.Now())
		go wait.Forever(func() {
			l.lastVerified.Store(time.Now())
		}, time.Minute)
	})

	lastVerifiedTime := l.lastVerified.Load().(time.Time)
	if time.Since(lastVerifiedTime) < (2 * time.Minute) {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, "health")
	} else {
		writer.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(writer, "unhealth")
	}
}
