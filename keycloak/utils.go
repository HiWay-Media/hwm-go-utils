package keycloak

import (
	"github.com/Nerzal/gocloak/v10"
))

func (g *gkeycloak) LoginTokenOptions() gocloak.TokenOptions {
	return gocloak.TokenOptions{
		GrantType:    &GRANT_PASSWORD,
		ClientID:     &g.clientId,
		ClientSecret: &g.clientSecret,
	}
}