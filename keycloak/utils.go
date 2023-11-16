package keycloak

import (
	"github.com/Nerzal/gocloak/v10"
)

func (g *gkeycloak) LoginTokenOptions() gocloak.TokenOptions {
	grant := GRANT_PASSWORD
	return gocloak.TokenOptions{
		GrantType:    grant,
		ClientID:     &g.clientId,
		ClientSecret: &g.clientSecret,
	}
}