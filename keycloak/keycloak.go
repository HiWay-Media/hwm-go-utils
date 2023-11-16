package keycloak

import (
	"context"
	"log"
	"sync"

	"github.com/Nerzal/gocloak/v10"
)

type gkeycloak struct {
	ctx          context.Context
	debug        bool
	clientId     string
	clientSecret string
	realm        string
	server       string
	adminJWT     *gocloak.JWT
	client       gocloak.GoCloak
	Mu           sync.Mutex
}

type IKeycloak interface {
	IsDebug() bool
	GetToken(tokenOptions gocloak.TokenOptions) (*gocloak.JWT, error)
}

func NewKeycloak(ctx context.Context, realm string, server string, clientId string, clientSecret string, isDebug bool) (IKeycloak, error) {
	k := &gkeycloak{
		ctx:          ctx,
		debug:        isDebug,
		clientId:     clientId,
		clientSecret: clientSecret,
		realm:        realm,
		server:       server,
		client:       gocloak.NewClient(server),
		Mu:           sync.Mutex{},
	}
	k.client.RestyClient().SetDebug(isDebug)
	//
	grantType := "client_credentials"
	t, err := k.GetToken(gocloak.TokenOptions{
		GrantType:    &grantType,
		ClientID:     &k.clientId,
		ClientSecret: &k.clientSecret,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("token: %+v", t)
	//
	return k, nil
}

//

func (g *gkeycloak) IsDebug() bool {
	return g.debug
}
