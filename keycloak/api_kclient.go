package keycloak

import (
	"github.com/Nerzal/gocloak/v10"
)

/*
*/
func (g *gkeycloak) AddClientRoleToUser( idOfClient, userID string, roles []gocloak.Role ) error {
	err := g.client.AddClientRoleToUser(g.ctx, g.adminJWT.AccessToken, g.realm, idOfClient, userID, roles )
	if err != nil {
		return nil, err
	}
	return nil
}
/*
*/