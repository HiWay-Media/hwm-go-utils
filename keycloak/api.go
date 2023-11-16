package keycloak

import (
	"gocloak.JWT"
)

func (g *gkeycloak) GetToken( tokenOptions gocloak.TokenOptions ) (*JWT, error) {
	token, err :=  g.GCloakClient.GetToken(g.ctx, g.realm, tokenOptions)
	if err != nil {
		return, nil, err
	}
	return token, nil
}