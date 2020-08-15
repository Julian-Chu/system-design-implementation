package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
)

// app encapsulates Env, router and middleware

type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	config      *Env
}

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

func (a *App) InitializeRoutes() {
	//a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	//a.Router.HandleFunc("/api/info", a.getShortLinkInfo).Methods("GET")
	//a.Router.HandleFunc("/{shortlink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")

	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLinkInfo)).Methods("GET")
	a.Router.Handle("/api/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) Initialize(e *Env) {
	// set log formatter
	// LstdFlags: print time and date
	// Lshortfile: line number, file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	a.config = e
	a.Router = mux.NewRouter()
	a.Middlewares = &Middleware{}
	a.InitializeRoutes()
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("parse parameter failed %v", r.Body)})
		return
	}
	defer r.Body.Close()
	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("validate parameter failed %v", r.Body)})
		return
	}

	fmt.Printf("%v\n", req)
}

func respondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s", e.Status(), e)
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	s := vals.Get("shortlink")
	fmt.Printf("%s\n", s)
	// simulate panic
	panic(s)
}

// Run starts listen and server
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortlink"])
}

func (a App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
