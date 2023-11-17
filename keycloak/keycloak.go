package keycloak

import (
	"context"
	"log"
	"sync"

	"github.com/Nerzal/gocloak/v10"
)
/**/
type gkeycloak struct {
	ctx          context.Context
	debug        bool
	clientId     string
	clientSecret string
	realm        string
	server       string
	adminJWT     *gocloak.JWT
	client       gocloak.GoCloak
	Mu           sync.Mutex
}
/**/
type IKeycloak interface {
	//
	IsDebug() bool
	Login(username string, password string) (*gocloak.JWT, error)
	GetToken( tokenOptions gocloak.TokenOptions ) (*gocloak.JWT, error)
	RefreshToken( refreshToken string ) (*gocloak.JWT, error)
	Logout(refreshToken string) error
	GetUserEmail( email string) (*gocloak.User, error)
	UpdateUser( firstName string, lastName string, username string, attributes map[string][]string, realmRoles []string) (bool, error)
	SetPassword(userID, realm, password string, temporary bool) error
	LogoutUserSession( session string ) error
	CreateGroup( group gocloak.Group ) (string, error)
	// *** Client Roles ***
	AddClientRoleToUser( idOfClient, userID string, roles []gocloak.Role ) error
	//
}

func NewKeycloak(ctx context.Context, realm string, server string, clientId string, clientSecret string, isDebug bool) (IKeycloak, error) {
	k := &gkeycloak{
		ctx:          ctx,
		debug:        isDebug,
		clientId:     clientId,
		clientSecret: clientSecret,
		realm:        realm,
		server:       server,
		client:       gocloak.NewClient(server),
		Mu:           sync.Mutex{},
	}
	k.client.RestyClient().SetDebug(isDebug)
	//
	grantType := GRANT_CLIENT_CREDENTIALS
	t, err := k.GetToken(gocloak.TokenOptions{
		GrantType:    &grantType,
		ClientID:     &k.clientId,
		ClientSecret: &k.clientSecret,
	})
	if err != nil {
		return nil, err
	}
	k.debugPrint("token: ", t)
	//
	return k, nil
}

//

func (g *gkeycloak) IsDebug() bool {
	return g.debug
}

func (g *gkeycloak) debugPrint(d ...interface{}) {
	if g.IsDebug() {
		log.Println(d...)
	}
}