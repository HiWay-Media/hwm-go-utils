package nomad_test

import (
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
