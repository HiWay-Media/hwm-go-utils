package keycloak

import (
	"sync"
	"github.com/Nerzal/gocloak/v10"
)

type gkeycloak struct {
	debug 			bool
	clientId     	string
	clientSecret 	string
	realm        	string
	server       	string
	gocloak.GoCloak
	Mu sync.Mutex
}

type IKeycloak interface {
	IsDebug()
}

func NewKeycloak() (IKeycloak, error)  {
	return nil, nil
}

func (g *gkeycloak) IsDebug() bool {
	return g.debug
}