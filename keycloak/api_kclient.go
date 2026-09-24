package keycloak

import (
	"github.com/Nerzal/gocloak/v10"
)

/*
*/
func (g *gkeycloak) AddClientRoleToUser( idOfClient, userID string, roles []gocloak.Role ) error {
	token, err := g.adminToken()
	if err != nil {
		return err
	}
	err = g.client.AddClientRoleToUser(g.ctx, token, g.realm, idOfClient, userID, roles )
	if err != nil {
		return  err
	}
	return nil
}
/*
*/