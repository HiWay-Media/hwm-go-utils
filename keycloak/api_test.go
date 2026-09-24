package keycloak_test

import (
	"context"
	"os"
	"testing"

	"github.com/HiWay-Media/hwm-go-utils/keycloak"
)

func TestAPI(t *testing.T) {
	realm := os.Getenv("KEYCLOAK_REALM")
	server := os.Getenv("KEYCLOAK_SERVER")
	clientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	if server == "" {
		t.Skip("KEYCLOAK_SERVER not set, skipping integration test")
	}
	k, err := keycloak.NewKeycloak(context.Background(), realm, server, clientId, clientSecret, false)
	if err != nil {
		t.Fatal(err)
	}
	k.IsDebug()
	//log.Println(k)
}
