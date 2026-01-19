package api

import (
	"errors"
	"goapp/internal/pkg/database"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UsersForm struct {
	Name  string
	Email string
	Age   string
}

type UsersPageData struct {
	Users []database.User
	Form  UsersForm
	Error string

	Page     int
	Limit    int
	PrevPage int
	NextPage int
}

type EditPageData struct {
	User  *database.User
	Error string
}

func (api *Api) GetUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		parsed, err := strconv.Atoi(p)
		if err != nil || parsed < 1 {
			w.WriteHeader(http.StatusBadRequest)
			api.renderTemplate(w, "users.html", UsersPageData{
				Error: "invalid page",
				Page:  1,
				Limit: limit,
			})
			return
		}
		page = parsed
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 {
			w.WriteHeader(http.StatusBadRequest)
			api.renderTemplate(w, "users.html", UsersPageData{
				Error: "invalid limit",
				Page:  page,
				Limit: 10,
			})
			return
		}
		if parsed > 100 {
			parsed = 100
		}
		limit = parsed
	}

	offset := (page - 1) * limit

	users, err := api.db.GetUsers(r.Context(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.renderTemplate(w, "users.html", UsersPageData{
			Error: "failed to fetch users",
			Page:  page,
			Limit: limit,
		})
		return
	}

	prevPage := 0
	if page > 1 {
		prevPage = page - 1
	}
	nextPage := page + 1

	api.renderTemplate(w, "users.html", UsersPageData{
		Users:    users,
		Page:     page,
		Limit:    limit,
		PrevPage: prevPage,
		NextPage: nextPage,
	})
}

func (api *Api) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.renderTemplate(w, "edit.html", EditPageData{
			Error: "invalid id",
		})
		return
	}

	user, err := api.db.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			api.renderTemplate(w, "edit.html", EditPageData{
				Error: "user not found",
			})
			return
		}
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		api.renderTemplate(w, "edit.html", EditPageData{
			Error: "failed to fetch user",
		})
		return
	}

	api.renderTemplate(w, "edit.html", EditPageData{
		User: user,
	})
}

func (api *Api) CreateUser(w http.ResponseWriter, r *http.Request) {
	const page = 1
	const limit = 10
	const offset = 0

	render := func(status int, msg string, form UsersForm) {
		users, err := api.db.GetUsers(r.Context(), limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			api.renderTemplate(w, "users.html", UsersPageData{
				Error:    "failed to fetch users",
				Page:     page,
				Limit:    limit,
				PrevPage: 0,
				NextPage: 2,
			})
			return
		}

		w.WriteHeader(status)
		api.renderTemplate(w, "users.html", UsersPageData{
			Users:    users,
			Form:     form,
			Error:    msg,
			Page:     page,
			Limit:    limit,
			PrevPage: 0,
			NextPage: 2,
		})
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")

	if name == "" || email == "" {
		render(http.StatusBadRequest, "name and email are required", UsersForm{
			Name:  name,
			Email: email,
			Age:   ageStr,
		})
		return
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		render(http.StatusBadRequest, "age must be a number", UsersForm{
			Name:  name,
			Email: email,
			Age:   ageStr,
		})
		return
	}
	if age <= 0 {
		render(http.StatusBadRequest, "age must be greater than 0", UsersForm{
			Name:  name,
			Email: email,
			Age:   ageStr,
		})
		return
	}

	user := &database.User{
		Name:  name,
		Email: email,
		Age:   age,
	}

	if err := api.db.CreateUser(r.Context(), user); err != nil {
		log.Print(err)
		render(http.StatusInternalServerError, "failed to create user", UsersForm{
			Name:  name,
			Email: email,
			Age:   ageStr,
		})
		return
	}

	http.Redirect(w, r, "/users?page=1&limit=10", http.StatusSeeOther)
}

func (api *Api) EditUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.renderTemplate(w, "edit.html", EditPageData{
			Error: "invalid id",
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	ageStr := r.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.renderTemplate(w, "edit.html", EditPageData{
			User:  &database.User{ID: id, Name: name, Email: email},
			Error: "age must be a number",
		})
		return
	}

	if name == "" || email == "" {
		w.WriteHeader(http.StatusBadRequest)
		api.renderTemplate(w, "edit.html", EditPageData{
			User:  &database.User{ID: id, Name: name, Email: email, Age: age},
			Error: "name and email are required",
		})
		return
	}
	if age <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		api.renderTemplate(w, "edit.html", EditPageData{
			User:  &database.User{ID: id, Name: name, Email: email, Age: age},
			Error: "age must be greater than 0",
		})
		return
	}

	user := &database.User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}

	if err := api.db.UpdateUser(r.Context(), user); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			api.renderTemplate(w, "edit.html", EditPageData{
				User:  user,
				Error: "user not found",
			})
			return
		}
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		api.renderTemplate(w, "edit.html", EditPageData{
			User:  user,
			Error: "failed to update user",
		})
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (api *Api) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	err = api.db.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		log.Printf("failed to delete user %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (api *Api) renderTemplate(w http.ResponseWriter, name string, data any) {
	if err := api.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template execution failed (%s): %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
