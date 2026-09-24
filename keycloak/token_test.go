package keycloak

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func fakeTokenServer(t *testing.T, calls *int32) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/protocol/openid-connect/token") {
			http.NotFound(w, r)
			return
		}
		n := atomic.AddInt32(calls, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":"token-%d","expires_in":3600}`, n)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestAdminTokenIsStoredCachedAndRefreshed(t *testing.T) {
	var calls int32
	srv := fakeTokenServer(t, &calls)

	k, err := NewKeycloak(context.Background(), "realm", srv.URL, "client", "secret", false)
	if err != nil {
		t.Fatalf("NewKeycloak: %v", err)
	}
	g := k.(*gkeycloak)

	if g.adminJWT == nil || g.adminJWT.AccessToken != "token-1" {
		t.Fatalf("admin token not stored after NewKeycloak: %+v", g.adminJWT)
	}

	tok, err := g.adminToken()
	if err != nil || tok != "token-1" || atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("want cached token-1 with 1 call, got %q err=%v calls=%d", tok, err, calls)
	}

	// token about to expire: must be refreshed
	g.adminExpiry = time.Now().Add(adminTokenLeeway / 2)
	tok, err = g.adminToken()
	if err != nil || tok != "token-2" || atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("want refreshed token-2 with 2 calls, got %q err=%v calls=%d", tok, err, calls)
	}
}
