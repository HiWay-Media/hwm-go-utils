package keycloak

import (
	"context"
	"sync"
	"github.com/Nerzal/gocloak/v10"
	"gopkg.in/square/go-jose.v2"
)

type gkeycloak struct {
	debug 			bool
	clientId     	string
	clientSecret 	string
	realm        	string
	server       	string
	adminJWT 		*gocloak.JWT
	gocloak.GoCloak
	Mu sync.Mutex
}

type IKeycloak interface {
	IsDebug() bool
	GetToken( tokenOptions gocloak.TokenOptions ) (*JWT, error)
}

func NewKeycloak(ctx context.Context, realm string, server string, clientId string, clientSecret string, isDebug bool) (IKeycloak, error)  {
	k := &gkeycloak{
		ctx: 			ctx,
		debug: 			isDebug,
		clientId:     	clientId,
		clientSecret: 	clientSecret,
		realm:        	realm,
		server:       	server,
		GoCloak:		gocloak.NewClient(server),
		Mu:           	sync.Mutex{},
	}
	k.GoCloak.RestyClient().SetDebug(isDebug)
	//
	grantType := "client_credentials"
	t, err := k.GetToken(gocloak.TokenOptions{
		GrantType:    &grantType,
		ClientID:     &GCloakClient.clientId,
		ClientSecret: &GCloakClient.clientSecret,
	})
	if err != nil {
		return, nil, err
	}
	//
	return k, nil
}

// 

func (g *gkeycloak) IsDebug() bool {
	return g.debug
}