package keycloak 

import (
	"github.com/Nerzal/gocloak/v10"
)

// *** Realm ***
func (g *gkeycloak) GetRealm( realm string ) (*gocloak.RealmRepresentation, error) {
	realmRep, err := g.client.GetRealm(g.ctx, g.adminJWT.AccessToken, realm ) 
	if err != nil {
		return nil, err
	}
	return realmRep, nil
}

//