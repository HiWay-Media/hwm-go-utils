package keycloak 

import (
	"github.com/Nerzal/gocloak/v10"
)

// *** Realm ***
func (g *gkeycloak) GetRealm( realm string ) (*gocloak.RealmRepresentation, error) {
	token, err := g.adminToken()
	if err != nil {
		return nil, err
	}
	realmRep, err := g.client.GetRealm(g.ctx, token, realm ) 
	if err != nil {
		return nil, err
	}
	return realmRep, nil
}

//
func (g *gkeycloak) GetRealms( ) ([]*gocloak.RealmRepresentation, error) {
	token, err := g.adminToken()
	if err != nil {
		return nil, err
	}
	realms, err := g.client.GetRealms(g.ctx, token ) 
	if err != nil {
		return nil, err
	}
	return realms, nil
}

/*
*/
func (g *gkeycloak) CreateRealm( realm gocloak.RealmRepresentation ) (string, error) {
	token, err := g.adminToken()
	if err != nil {
		return "", err
	}
	resp, err := g.client.CreateRealm(g.ctx, token, realm ) 
	if err != nil {
		return "", err
	}
	return resp, nil
}

/*
*/
func (g *gkeycloak) UpdateRealm( realm gocloak.RealmRepresentation) error{
	token, err := g.adminToken()
	if err != nil {
		return err
	}
 	err = g.client.UpdateRealm(g.ctx, token, realm ) 
	if err != nil {
		return err
	}
	return nil
}

/*
*/
func (g *gkeycloak) DeleteRealm(  realm string ) error {
	token, err := g.adminToken()
	if err != nil {
		return err
	}
	err = g.client.DeleteRealm(g.ctx, token, realm ) 
	if err != nil {
		return err
	}
	return nil
}