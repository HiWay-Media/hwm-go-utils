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
func (g *gkeycloak) GetRealms( ) ([]*gocloak.RealmRepresentation, error) {
	realms, err := g.client.GetRealms(g.ctx, g.adminJWT.AccessToken ) 
	if err != nil {
		return nil, err
	}
	return realms, nil
}

/*
*/
func (g *gkeycloak) CreateRealm( realm gocloak.RealmRepresentation ) (string, error) {
	resp, err := g.client.CreateRealm(g.ctx, g.adminJWT.AccessToken, realm ) 
	if err != nil {
		return "", err
	}
	return resp, nil
}

/*
*/
func (g *gkeycloak) UpdateRealm( realm gocloak.RealmRepresentation) error{
 	err := g.client.UpdateRealm(g.ctx, g.adminJWT.AccessToken, realm ) 
	if err != nil {
		return err
	}
	return nil
}

/*
*/
func (g *gkeycloak) DeleteRealm(  realm string ) error {
	err := g.client.DeleteRealm(g.ctx, g.adminJWT.AccessToken, realm ) 
	if err != nil {
		return err
	}
	return nil
}