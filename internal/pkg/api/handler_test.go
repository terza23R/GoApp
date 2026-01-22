package api

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"goapp/internal/pkg/database"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type fakeUserRepo struct {
	getUsersFn    func(ctx context.Context, limit, offset int) ([]database.User, error)
	getUserByIDFn func(ctx context.Context, id int64) (*database.User, error)
	createUserFn  func(ctx context.Context, u *database.User) error
	updateUserFn  func(ctx context.Context, u *database.User) error
	deleteUserFn  func(ctx context.Context, id int64) error
}

func (f *fakeUserRepo) GetUsers(ctx context.Context, limit, offset int) ([]database.User, error) {
	if f.getUsersFn != nil {
		return f.getUsersFn(ctx, limit, offset)
	}
	return []database.User{}, nil
}

func (f *fakeUserRepo) GetUserByID(ctx context.Context, id int64) (*database.User, error) {
	if f.getUserByIDFn != nil {
		return f.getUserByIDFn(ctx, id)
	}
	return nil, database.ErrUserNotFound
}

func (f *fakeUserRepo) CreateUser(ctx context.Context, u *database.User) error {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, u)
	}
	return nil
}

func (f *fakeUserRepo) UpdateUser(ctx context.Context, u *database.User) error {
	if f.updateUserFn != nil {
		return f.updateUserFn(ctx, u)
	}
	return nil
}

func (f *fakeUserRepo) DeleteUser(ctx context.Context, id int64) error {
	if f.deleteUserFn != nil {
		return f.deleteUserFn(ctx, id)
	}
	return nil
}

func newTestAPI(repo UserRepository) *Api {
	tpl := template.Must(template.New("root").Parse(`
		{{define "users.html"}}ERROR={{.Error}}{{end}}
		{{define "edit.html"}}ERROR={{.Error}}{{end}}
	`))

	return &Api{
		router:    mux.NewRouter(),
		db:        repo,
		templates: tpl,
	}
}

func TestHealth(t *testing.T) {
	api := newTestAPI(&fakeUserRepo{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	api.Health(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "ok" {
		t.Fatalf("expected body 'ok', got %q", w.Body.String())
	}
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	api := newTestAPI(&fakeUserRepo{
		getUsersFn: func(ctx context.Context, limit, offset int) ([]database.User, error) {
			return []database.User{}, nil
		},
	})

	form := url.Values{}
	form.Set("name", "Mahir")
	form.Set("email", "not-an-email")
	form.Set("age", "24")

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	api.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid email format") {
		t.Fatalf("expected invalid email format, got %q", w.Body.String())
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	api := newTestAPI(&fakeUserRepo{
		getUsersFn: func(ctx context.Context, limit, offset int) ([]database.User, error) {
			return []database.User{}, nil
		},
		createUserFn: func(ctx context.Context, u *database.User) error {
			return &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
		},
	})

	form := url.Values{}
	form.Set("name", "Mahir")
	form.Set("email", "mahir@test.com")
	form.Set("age", "24")

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	api.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "email already exists") {
		t.Fatalf("expected email already exists, got %q", w.Body.String())
	}
}

func TestGetUsers_DBError(t *testing.T) {
	api := newTestAPI(&fakeUserRepo{
		getUsersFn: func(ctx context.Context, limit, offset int) ([]database.User, error) {
			return nil, errors.New("db down")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	api.GetUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "failed to fetch users") {
		t.Fatalf("expected failed to fetch users, got %q", w.Body.String())
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	api := newTestAPI(&fakeUserRepo{
		deleteUserFn: func(ctx context.Context, id int64) error {
			return database.ErrUserNotFound
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/users/123", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "123"})
	w := httptest.NewRecorder()

	api.DeleteUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "user not found") {
		t.Fatalf("expected user not found, got %q", w.Body.String())
	}
}
