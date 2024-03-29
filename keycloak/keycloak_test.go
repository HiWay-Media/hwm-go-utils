package keycloak_test

import (
	"context"
	"os"
	"testing"

	"github.com/HiWay-Media/hwm-go-utils/keycloak"
)

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "" {
		err := os.Setenv("APP_ENV", "test")
		if err != nil {
			panic("could not set test env")
		}
	}
	
	m.Run()
}

func TestIKeycloak(t *testing.T) {
	realm := os.Getenv("KEYCLOAK_REALM")
	server := os.Getenv("KEYCLOAK_SERVER")
	clientId := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	k, err := keycloak.NewKeycloak(context.Background(), realm, server, clientId, clientSecret, false)
	if err != nil {
		t.Fatalf(err.Error())
	}
	k.IsDebug()
	//log.Println(k)
}
