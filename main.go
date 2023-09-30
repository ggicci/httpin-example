package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ggicci/httpin"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	router = mux.NewRouter() // use gorilla/mux for routing

	usersDB []*User
	reposDB []*Repository
)

type User struct {
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
	IsMember  bool      `json:"is_member"`
	Age       int       `json:"age"`
}

type Repository struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

// API (Create a user): POST /users

type CreateUserInput struct {
	*User `in:"body"`
}

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(httpin.Input).(*CreateUserInput)
	input.CreatedAt = time.Now()
	usersDB = append(usersDB, input.User)

	sendJSON(rw, map[string]interface{}{
		"input": input,
		"users": usersDB,
	})
}

// API (List users): GET /users?is_member=true&sort_by[]=age&sort_desc[]=false

type ListUsersInput struct {
	IsMember bool     `in:"query=is_member,vip"`
	SortBy   []string `in:"query=sort_by[]"`
	SortDesc []bool   `in:"query=sort_desc[]"`
}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	var (
		input = r.Context().Value(httpin.Input).(*ListUsersInput)
		res   []*User
	)

	for _, user := range usersDB {
		if user.IsMember == input.IsMember {
			res = append(res, user)
		}
	}

	sendJSON(rw, map[string]interface{}{
		"users": res,
		"input": input,
	})
}

// API (Create a repo): POST /users/{login}/repos (require access token)

type TokenInput struct {
	Token string `in:"header=X-Api-Token,Authorization;query=token,access_token;required"`
}

type CreateRepositoryInput struct {
	TokenInput                // NOTE: httpin does support embedded structs
	NewRepository *Repository `in:"body"`
}

func CreateRepository(rw http.ResponseWriter, r *http.Request) {
	var (
		input = r.Context().Value(httpin.Input).(*CreateRepositoryInput)
	)

	if input.Token != "secret" {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	reposDB = append(reposDB, input.NewRepository)
	sendJSON(rw, map[string]interface{}{
		"input": input,
		"repos": reposDB,
	})
}

// API (List repos): GET /users/{login}/repos (require access token)

type ListRepositoriesOfUserInput struct {
	TokenInput        // NOTE: httpin does support embedded structs
	Login      string `in:"path=login"`
	Language   string `in:"query=lang"`
}

func ListRepositoriesOfUser(rw http.ResponseWriter, r *http.Request) {
	var (
		input = r.Context().Value(httpin.Input).(*ListRepositoriesOfUserInput)
	)

	if input.Login != "ggicci" {
		http.NotFound(rw, r)
		return
	}

	if input.Token != "secret" {
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	sendJSON(rw, map[string]interface{}{
		"repos": reposDB,
		"input": input,
	})
}

func init() {
	httpin.UseGorillaMux("path", mux.Vars)

	router.Handle("/users", alice.New(
		// httpin.NewInput creates a middleware of a specified type of struct
		// that parses the request and stores the result in the context
		httpin.NewInput(CreateUserInput{}),
	).ThenFunc(CreateUser)).Methods("POST")

	router.Handle("/users", alice.New(
		httpin.NewInput(ListUsersInput{}),
	).ThenFunc(ListUsers)).Methods("GET")

	router.Handle("/users/{login}/repos", alice.New(
		httpin.NewInput(CreateRepositoryInput{}),
	).ThenFunc(CreateRepository)).Methods("POST")

	router.Handle("/users/{login}/repos", alice.New(
		httpin.NewInput(ListRepositoriesOfUserInput{}),
	).ThenFunc(ListRepositoriesOfUser)).Methods("GET")
}

func main() {
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}

func sendJSON(rw http.ResponseWriter, v interface{}, statusCode ...int) (err error) {
	rw.Header().Add("Content-Type", "application/json")

	if len(statusCode) == 1 {
		rw.WriteHeader(statusCode[0])
	}
	return json.NewEncoder(rw).Encode(v)
}
