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
	IsDebug() bool
}

func NewKeycloak(realm, server, clientId, realm, clientSecret string, isDebug bool) (IKeycloak, error)  {
	k := gkeycloak{
		debug: 			isDebug,
		clientId:     	clientId,
		clientSecret: 	clientSecret,
		realm:        	realm,
		server:       	server,
		GoCloak:		gocloak.NewClient(server),
		Mu:           	sync.Mutex{},
	}
	k.GoCloak.RestyClient().SetDebug(isDebug)
	return k, nil
}

func (g *gkeycloak) IsDebug() bool {
	return g.debug
}