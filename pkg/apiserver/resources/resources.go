package resources

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
)

var Default = schema{
	apiResource:    map[string]APIWebService{},
	injectResource: map[string]Injecter{},
}

type Injecter interface {
	Inject(opt *RestOption)
}

type APIWebService interface {
	GenericWebService() *restful.WebService
}

type InjectFunc func(opt *RestOption)

func (f InjectFunc) Inject(opt *RestOption) {
	f(opt)
	return
}

type GenericWebSvcFunc func() *restful.WebService

func (r GenericWebSvcFunc) GenericWebService() *restful.WebService {
	return r()
}

type schema struct {
	apiResource    map[string]APIWebService
	injectResource map[string]Injecter
}

func (s *schema) AddAPIRegisterResource(name string, register APIWebService) {
	if s.apiResource == nil {
		s.apiResource = make(map[string]APIWebService)
	}
	s.apiResource[name] = register
}

func (s *schema) AddInjectResourceObj(name string, inject Injecter) {
	if s.injectResource == nil {
		s.injectResource = make(map[string]Injecter)
	}
	s.injectResource[name] = inject
}

// GetAPIRegisterResource return APIWebService by name
func (s *schema) GetAPIRegisterResource(name string) (APIWebService, error) {
	r, ok := s.apiResource[name]
	if !ok {
		return nil, fmt.Errorf("%s API Register Reousrce not exist", name)
	}

	return r, nil
}

// GetInjectObj return Injecter by name
func (s *schema) GetInjectObj(name string) (Injecter, error) {
	r, ok := s.injectResource[name]
	if !ok {
		return nil, fmt.Errorf("%s Inject Object not exist", name)
	}

	return r, nil
}

// RegisteredWebServices returns the collections of added WebServices
func (s *schema) RegisteredWebServices() []*restful.WebService {

	result := []*restful.WebService{}

	for _, ws := range s.apiResource {
		result = append(result, ws.GenericWebService())
	}

	return result
}
