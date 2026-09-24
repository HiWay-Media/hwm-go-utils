package generic

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type owner struct {
	ID   uint
	Name string
}

type pet struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	OwnerID uint   `json:"owner_id"`
	Owner   *owner `json:"owner"`
}

type fakeService struct {
	created *pet
	err     error
	listArg [2]int
}

func (f *fakeService) Get(id any) (*pet, error) { return &pet{}, f.err }
func (f *fakeService) Delete(id any) error      { return f.err }
func (f *fakeService) Create(obj *pet) error    { f.created = obj; return f.err }
func (f *fakeService) List(start, limit int) ([]pet, error) {
	f.listArg = [2]int{start, limit}
	return nil, f.err
}

func newHandlerApp(svc IService[pet]) *fiber.App {
	app := fiber.New()
	setRoutes[pet]("pets", app, NewHandler[pet](svc, nil))
	return app
}

func call(t *testing.T, app *fiber.App, method, target, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b)
}

func TestCreateClearsServerOwnedFields(t *testing.T) {
	svc := &fakeService{}
	status, _ := call(t, newHandlerApp(svc), http.MethodPost, "/pets/",
		`{"id": 99, "name": "rex", "owner_id": 7, "owner": {"ID": 7, "Name": "mallory"}}`)

	if status != http.StatusOK {
		t.Fatalf("status %d", status)
	}
	if svc.created == nil {
		t.Fatal("service not called")
	}
	if svc.created.ID != 0 || svc.created.Owner != nil {
		t.Errorf("primary key / association not cleared: %+v", svc.created)
	}
	if svc.created.Name != "rex" || svc.created.OwnerID != 7 {
		t.Errorf("regular fields must be kept: %+v", svc.created)
	}
}

func TestErrorsDoNotLeakDetails(t *testing.T) {
	secret := errors.New("Error 1045: Access denied for user 'app'@'10.0.0.5'")

	cases := []struct {
		name       string
		err        error
		method     string
		target     string
		wantStatus int
		wantBody   string
	}{
		{"get internal", secret, http.MethodGet, "/pets/1", 500, "internal error"},
		{"get not found", gorm.ErrRecordNotFound, http.MethodGet, "/pets/1", 404, "not found"},
		{"delete not found", gorm.ErrRecordNotFound, http.MethodDelete, "/pets/abc", 404, "not found"},
		{"list internal", secret, http.MethodGet, "/pets/", 500, "internal error"},
		{"create internal", secret, http.MethodPost, "/pets/", 500, "internal error"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			status, body := call(t, newHandlerApp(&fakeService{err: c.err}), c.method, c.target, `{"name":"rex"}`)
			if status != c.wantStatus || !strings.Contains(body, c.wantBody) || strings.Contains(body, "Access denied") {
				t.Errorf("got %d %s, want %d containing %q", status, body, c.wantStatus, c.wantBody)
			}
		})
	}
}

func TestListRejectsNegativePaging(t *testing.T) {
	for _, q := range []string{"start=-1", "limit=-5", "start=x"} {
		svc := &fakeService{}
		if status, _ := call(t, newHandlerApp(svc), http.MethodGet, "/pets/?"+q, ""); status != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", q, status)
		}
	}

	svc := &fakeService{}
	if status, _ := call(t, newHandlerApp(svc), http.MethodGet, "/pets/?start=10&limit=20", ""); status != http.StatusOK || svc.listArg != [2]int{10, 20} {
		t.Errorf("status %d, args %v", status, svc.listArg)
	}
}
