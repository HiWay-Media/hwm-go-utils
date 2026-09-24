package nomad_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/HiWay-Media/hwm-go-utils/log"
	"github.com/HiWay-Media/hwm-go-utils/nomad"
)

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "" {
		err := os.Setenv("APP_ENV", "test")
		if err != nil {
			panic("could not set test env")
		}
	}
	//env.Load()
	m.Run()
}

func getNomad() nomad.IService {
	options := nomad.Options{
		BaseUrl:  "url",
		LogLevel: "debug",
		Logger:   log.GetLogger("debug"),
	}
	return nomad.NewService(options)
}

func TestRequestPaths(t *testing.T) {
	var got []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.EscapedPath())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	want := []string{
		"GET /v1/job/job%201",
		"GET /v1/client/allocation/alloc/stats",
		"GET /v1/node/node/allocations",
		"DELETE /v1/job/job%201",
		"POST /v1/job/job%201/scale",
		"POST /v1/jobs",
	}

	for _, base := range []string{srv.URL, srv.URL + "/", srv.URL + "/v1", srv.URL + "/v1/"} {
		got = nil
		s := nomad.NewService(nomad.Options{BaseUrl: base})
		_, _ = s.GetDefinition("job 1", "eu")
		_, _ = s.AllocationStats("alloc", "eu")
		_, _ = s.GetAllocations("node", "eu")
		_ = s.DeleteJob("job 1", "eu", false)
		_ = s.ScaleJob("job 1", 1, "eu")
		_ = s.RunJob(nomad.JobDefinition{}, "eu")

		if !reflect.DeepEqual(got, want) {
			t.Errorf("base %q:\n got  %v\n want %v", base, got, want)
		}
	}
}

func TestTokenAndScaleGroup(t *testing.T) {
	var token string
	var body nomad.ScaleRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token = r.Header.Get("X-Nomad-Token")
		_ = json.NewDecoder(r.Body).Decode(&body)
		_, _ = w.Write([]byte("{}"))
	}))
	defer srv.Close()

	_ = nomad.NewService(nomad.Options{BaseUrl: srv.URL}).ScaleJob("job", 2, "eu")
	if token != "" || body.Target.Group != nomad.DefaultScaleGroup || body.Count != 2 {
		t.Errorf("defaults: token %q, body %+v", token, body)
	}

	_ = nomad.NewService(nomad.Options{BaseUrl: srv.URL, Token: "acl-secret", ScaleGroup: "encoder"}).ScaleJob("job", 1, "eu")
	if token != "acl-secret" || body.Target.Group != "encoder" {
		t.Errorf("options: token %q, body %+v", token, body)
	}
}

func TestGetAllocationsDecodesNomadArray(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		// shape of GET /v1/node/:id/allocations: a JSON array of allocations
		_, _ = w.Write([]byte(`[
			{"ID": "a1", "NodeID": "node-1", "JobID": "restreamer-1", "TaskGroup": "restreamer", "ClientStatus": "running"},
			{"ID": "a2", "NodeID": "node-1", "JobID": "encoder-7", "TaskGroup": "encoder", "ClientStatus": "complete"}
		]`))
	}))
	defer srv.Close()

	allocs, err := nomad.NewService(nomad.Options{BaseUrl: srv.URL}).GetAllocations("node-1", "eu")
	if err != nil {
		t.Fatalf("GetAllocations: %v", err)
	}
	if path != "/v1/node/node-1/allocations" {
		t.Errorf("path %q", path)
	}
	if len(allocs.NomadAllocations) != 2 || allocs.NomadAllocations[0].ID != "a1" || allocs.NomadAllocations[1].JobID != "encoder-7" {
		t.Errorf("unexpected allocations: %+v", allocs.NomadAllocations)
	}
}

func TestNomadAllocationsJSONRoundTrip(t *testing.T) {
	in := nomad.NomadAllocations{NomadAllocations: []nomad.NomadAlloc{{ID: "a1"}}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out nomad.NomadAllocations
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("object form: %v", err)
	}
	if len(out.NomadAllocations) != 1 || out.NomadAllocations[0].ID != "a1" {
		t.Errorf("round trip lost data: %s -> %+v", b, out)
	}

	var empty nomad.NomadAllocations
	if err := json.Unmarshal([]byte(" [] "), &empty); err != nil || len(empty.NomadAllocations) != 0 {
		t.Errorf("empty array: err=%v, %+v", err, empty)
	}
}
