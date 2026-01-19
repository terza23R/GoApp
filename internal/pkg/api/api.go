package api

import (
	"context"
	"goapp/internal/pkg/database"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Api struct {
	address   string
	router    *mux.Router
	server    *http.Server
	db        *database.DB
	templates *template.Template
}

func NewApi(hostPort string, db *database.DB, templatesPath string) *Api {
	r := mux.NewRouter()
	tpl := template.Must(template.ParseGlob(templatesPath))

	api := &Api{
		address:   hostPort,
		router:    r,
		db:        db,
		templates: tpl,
	}

	api.registerHandlers()

	api.server = &http.Server{
		Addr:    api.address,
		Handler: api.router,
	}

	return api

}

func (api *Api) registerHandlers() {
	api.router.HandleFunc("/users", api.GetUsers).Methods(http.MethodGet)
	api.router.HandleFunc("/users/{id}", api.GetUser).Methods(http.MethodGet)
	api.router.HandleFunc("/users", api.CreateUser).Methods(http.MethodPost)
	api.router.HandleFunc("/users/{id}", api.EditUser).Methods(http.MethodPost)
	api.router.HandleFunc("/users/{id}/delete", api.DeleteUser).Methods(http.MethodPost)
	api.router.HandleFunc("/health", api.Health).Methods(http.MethodGet)
}

func (api *Api) Start() {
	go func() {
		if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()
}

func (api *Api) Stop() {
	if api.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = api.server.Shutdown(ctx)
}
