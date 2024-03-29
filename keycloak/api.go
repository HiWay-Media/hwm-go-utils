package keycloak

import (
	"fmt"
	"github.com/Nerzal/gocloak/v10"
)

/*
*/
func (g *gkeycloak) Login(username string, password string) (*gocloak.JWT, error) {
	token, err := g.client.Login(g.ctx, g.clientId, g.clientSecret, g.realm, username, password )
	if err != nil {
		return nil, err
	}
	return token, nil
}
/*
 */
func (g *gkeycloak) GetToken(tokenOptions gocloak.TokenOptions) (*gocloak.JWT, error) {
	token, err := g.client.GetToken(g.ctx, g.realm, tokenOptions)
	if err != nil {
		return nil, err
	}
	return token, nil
}

/*
*/
func (g *gkeycloak) RefreshToken( refreshToken string ) (*gocloak.JWT, error) {
	tokenRefreshed, err := g.client.RefreshToken(g.ctx, refreshToken, g.clientId, g.clientSecret, g.realm) 
	if err != nil {
		return nil, err
	}
	return tokenRefreshed, nil
}

/*
*/
func (g * gkeycloak) Logout(refreshToken string) error {
	err := g.client.Logout(g.ctx, g.realm, refreshToken, g.clientId, g.clientSecret)
	if err != nil {
		return err
	}
	return nil
}

/*
*/
func (g *gkeycloak) GetUserEmail( email string ) (*gocloak.User, error) {
	users, err := g.client.GetUsers(g.ctx, g.adminJWT.AccessToken, g.realm, gocloak.GetUsersParams{Email: &email})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

/*
*/
func (g *gkeycloak) CreateUser( user gocloak.User) (string, error) {
	resp, err := g.client.CreateUser(g.ctx, g.adminJWT.AccessToken, g.realm, user)
	if err != nil {
		return "", err
	}
	return resp, nil
}

/*
*/
func (g *gkeycloak) UpdateUser( firstName string, lastName string, username string, attributes map[string][]string, realmRoles []string) (bool, error) {
	g.debugPrint("into keycloak Updateuser")
	//getting user first
	user, err := g.GetUserEmail( username,)
	if err != nil {
		fmt.Errorf("failed to getting user: %s", err.Error())
		return false, err
	}
	g.debugPrint("attributes ", attributes)
	user.RealmRoles = &realmRoles
	user.FirstName = &firstName
	user.LastName = &lastName
	user.Email = &username
	//
	err = g.client.UpdateUser(g.ctx, g.adminJWT.AccessToken, g.realm, *user)
	if err != nil {
		return false, err
	}
	//
	return true, nil
}

/*
*/
func (g *gkeycloak) SetPassword(userID, realm, password string, temporary bool) error {
	g.debugPrint("into keycloak SetPassword")
	err := g.client.SetPassword(g.ctx, g.adminJWT.AccessToken, userID, g.realm, password, temporary)
	if err != nil {
		return err
	}
	return nil
}

/*
*/
func (g *gkeycloak) LogoutUserSession( session string ) error {
	err := g.client.LogoutUserSession( g.ctx, g.adminJWT.AccessToken, g.realm, session )
	if err != nil {
		return err
	}
	return nil
}

/*
*/
func (g *gkeycloak) CreateGroup( group gocloak.Group ) (string, error) {
	r, err := g.client.CreateGroup( g.ctx, g.adminJWT.AccessToken, g.realm, group )
	if err != nil {
		return "", err
	}
	return r, nil
}