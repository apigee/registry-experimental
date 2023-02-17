package edge

import (
	"path"
)

const environmentsPath = "environments"
const virtualHostsPath = "virtualhosts"

// EnvironmentsService is an interface for interfacing with the Apigee Edge Admin API
// querying Edge environments.
type EnvironmentsService interface {
	ListNames() ([]string, *Response, error)
	Get(env string) (*Environment, *Response, error)
	ListVirtualHosts(env string) ([]string, *Response, error)
	GetVirtualHost(env, name string) (*VirtualHost, *Response, error)
}

type EnvironmentsServiceOp struct {
	client *EdgeClient
}

var _ EnvironmentsService = &EnvironmentsServiceOp{}

// Environment contains information about an environment within an Edge organization.
type Environment struct {
	Name           string    `json:"name,omitempty"`
	CreatedBy      string    `json:"createdBy,omitempty"`
	CreatedAt      Timestamp `json:"createdAt,omitempty"`
	LastModifiedBy string    `json:"lastModifiedBy,omitempty"`
	LastModifiedAt Timestamp `json:"lastModifiedAt,omitempty"`
	APIProxies     []Proxy   `json:"aPIProxy,omitempty"`
	// Properties      PropertyWrapper `json:"properties,omitempty"`
}

type VirtualHost struct {
	Name        string   `json:"name,omitempty"`
	Port        string   `json:"port,omitempty"`
	BaseUrl     string   `json:"baseUrl,omitempty"`
	HostAliases []string `json:"hostAliases,omitempty"`
}

// List retrieves the list of environment names for the organization referred by the ApigeeClient.
func (s *EnvironmentsServiceOp) ListNames() ([]string, *Response, error) {
	req, e := s.client.NewRequestNoEnv("GET", environmentsPath, nil)
	if e != nil {
		return nil, nil, e
	}
	var names []string
	resp, e := s.client.Do(req, &names)
	if e != nil {
		return nil, resp, e
	}
	return names, resp, e
}

// Get retrieves the information about an Environment in an organization, information including
// the properties, and the created and last modified details.
func (s *EnvironmentsServiceOp) Get(env string) (*Environment, *Response, error) {
	path := path.Join(environmentsPath, env)
	req, e := s.client.NewRequestNoEnv("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	returnedEnv := Environment{}
	resp, e := s.client.Do(req, &returnedEnv)
	if e != nil {
		return nil, resp, e
	}
	return &returnedEnv, resp, e
}

func (s *EnvironmentsServiceOp) ListVirtualHosts(env string) ([]string, *Response, error) {
	path := path.Join(environmentsPath, env, virtualHostsPath)
	req, e := s.client.NewRequestNoEnv("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	var names []string
	resp, e := s.client.Do(req, &names)
	if e != nil {
		return nil, resp, e
	}
	return names, resp, e
}

func (s *EnvironmentsServiceOp) GetVirtualHost(env, name string) (*VirtualHost, *Response, error) {
	path := path.Join(environmentsPath, env, virtualHostsPath, name)
	req, e := s.client.NewRequestNoEnv("GET", path, nil)
	if e != nil {
		return nil, nil, e
	}
	vh := &VirtualHost{}
	resp, e := s.client.Do(req, vh)
	if e != nil {
		return nil, resp, e
	}

	return vh, resp, e
}
