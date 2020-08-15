package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

// app encapsulates Env, router and middleware

type App struct {
	Router *mux.Router
}

type shortenReq struct {
	URL                string `json:"url" validate:"nonzero"`
	ExpirationInMinute int64  `json:"expiration_in_minute" validate:"min=0"`
}

func (a *App) Initialize() {
	// set log formatter
	// LstdFlags: print time and date
	// Lshortfile: line number, file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.InitializeRoutes()

}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	a.Router.HandleFunc("/api/into", a.getShortLinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	defer r.Body.Close()
	if err := validator.Validate(req); err != nil {
		return
	}

	fmt.Printf("%v\n", req)
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	fmt.Printf("%s\n", s)
}

// Run starts listen and server
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
