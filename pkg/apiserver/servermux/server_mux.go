package servermux

import (
	"fmt"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	"github.com/HugoWw/x_apiserver/pkg/util/sets"
	"net/http"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type PathRecorderMux struct {
	name            string
	lock            sync.Mutex
	notFoundHandler http.Handler
	pathToHandler   map[string]http.Handler
	prefixToHandler map[string]http.Handler
	mux             atomic.Value
	exposedPaths    []string
	pathStacks      map[string]string
}

type pathHandler struct {
	muxName         string
	pathToHandler   map[string]http.Handler
	prefixHandlers  []prefixHandler
	notFoundHandler http.Handler
}

type prefixHandler struct {
	prefix  string
	handler http.Handler
}

func NewPathRecorderMux(name string) *PathRecorderMux {
	ret := &PathRecorderMux{
		name:            name,
		pathToHandler:   map[string]http.Handler{},
		prefixToHandler: map[string]http.Handler{},
		mux:             atomic.Value{},
		exposedPaths:    []string{},
		pathStacks:      map[string]string{},
	}

	ret.mux.Store(&pathHandler{notFoundHandler: http.NotFoundHandler()})
	return ret
}

func (m *PathRecorderMux) ListedPaths() []string {
	handledPaths := append([]string{}, m.exposedPaths...)
	sort.Strings(handledPaths)

	return handledPaths
}

func (m *PathRecorderMux) trackCallers(path string) {
	stack := string(debug.Stack())
	if existingStack, ok := m.pathStacks[path]; ok {
		clog.Logger.Errorf("duplicate path registration of %q: original registration from %v\n\nnew registration from %v", path, existingStack, stack)
	}
	m.pathStacks[path] = stack
}

func (m *PathRecorderMux) refreshMuxLocked() {
	newMux := &pathHandler{
		muxName:         m.name,
		pathToHandler:   map[string]http.Handler{},
		prefixHandlers:  []prefixHandler{},
		notFoundHandler: http.NotFoundHandler(),
	}
	if m.notFoundHandler != nil {
		newMux.notFoundHandler = m.notFoundHandler
	}
	for path, handler := range m.pathToHandler {
		newMux.pathToHandler[path] = handler
	}

	keySet := sets.KeySet(m.prefixToHandler)
	keys := sets.List(keySet)

	sort.Sort(sort.Reverse(byPrefixPriority(keys)))
	for _, prefix := range keys {
		newMux.prefixHandlers = append(newMux.prefixHandlers, prefixHandler{
			prefix:  prefix,
			handler: m.prefixToHandler[prefix],
		})
	}

	m.mux.Store(newMux)
}

func (m *PathRecorderMux) NotFoundHandler(notFoundHandler http.Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.notFoundHandler = notFoundHandler

	m.refreshMuxLocked()
}

func (m *PathRecorderMux) Unregister(path string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.pathToHandler, path)
	delete(m.prefixToHandler, path)
	delete(m.pathStacks, path)
	for i := range m.exposedPaths {
		if m.exposedPaths[i] == path {
			m.exposedPaths = append(m.exposedPaths[:i], m.exposedPaths[i+1:]...)
			break
		}
	}

	m.refreshMuxLocked()
}

func (m *PathRecorderMux) Handle(path string, handler http.Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.trackCallers(path)

	m.exposedPaths = append(m.exposedPaths, path)
	m.pathToHandler[path] = handler
	m.refreshMuxLocked()
}

func (m *PathRecorderMux) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(path, http.HandlerFunc(handler))
}

func (m *PathRecorderMux) UnlistedHandle(path string, handler http.Handler) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.trackCallers(path)

	m.pathToHandler[path] = handler
	m.refreshMuxLocked()
}

func (m *PathRecorderMux) UnlistedHandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	m.UnlistedHandle(path, http.HandlerFunc(handler))
}

func (m *PathRecorderMux) HandlePrefix(path string, handler http.Handler) {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Sprintf("%q must end in a trailing slash", path))
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.trackCallers(path)

	m.exposedPaths = append(m.exposedPaths, path)
	m.prefixToHandler[path] = handler
	m.refreshMuxLocked()
}

func (m *PathRecorderMux) UnlistedHandlePrefix(path string, handler http.Handler) {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Sprintf("%q must end in a trailing slash", path))
	}

	m.lock.Lock()
	defer m.lock.Unlock()
	m.trackCallers(path)

	m.prefixToHandler[path] = handler
	m.refreshMuxLocked()
}

func (m *PathRecorderMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.Load().(*pathHandler).ServeHTTP(w, r)
}

func (h *pathHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if exactHandler, ok := h.pathToHandler[r.URL.Path]; ok {

		clog.Logger.Infof("%v: %q satisfied by exact match", h.muxName, r.URL.Path)

		exactHandler.ServeHTTP(w, r)
		return
	}

	for _, prefixHandler := range h.prefixHandlers {
		if strings.HasPrefix(r.URL.Path, prefixHandler.prefix) {
			clog.Logger.Infof("%v: %q satisfied by prefix %v", h.muxName, r.URL.Path, prefixHandler.prefix)
			prefixHandler.handler.ServeHTTP(w, r)
			return
		}
	}

	clog.Logger.Infof("%v: %q satisfied by NotFoundHandler", h.muxName, r.URL.Path)

	h.notFoundHandler.ServeHTTP(w, r)
}

type byPrefixPriority []string

func (s byPrefixPriority) Len() int      { return len(s) }
func (s byPrefixPriority) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byPrefixPriority) Less(i, j int) bool {
	lhsNumParts := strings.Count(s[i], "/")
	rhsNumParts := strings.Count(s[j], "/")
	if lhsNumParts != rhsNumParts {
		return lhsNumParts < rhsNumParts
	}

	lhsLen := len(s[i])
	rhsLen := len(s[j])
	if lhsLen != rhsLen {
		return lhsLen < rhsLen
	}

	return strings.Compare(s[i], s[j]) < 0
}
