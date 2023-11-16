package keycloak

import "github.com/Nerzal/gocloak/v10"

/*
 */
func (g *gkeycloak) GetToken(tokenOptions gocloak.TokenOptions) (*gocloak.JWT, error) {
	token, err := g.client.GetToken(g.ctx, g.realm, tokenOptions)
	if err != nil {
		return nil, err
	}
	return token, nil
}
