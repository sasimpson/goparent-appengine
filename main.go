package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/datastore"
	"google.golang.org/appengine"
)

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "ok")
}

//This file is specifically for running in GCP AppEngine.
func main() {
	env := &goparent.Env{
		Service: goparent.Service{},
		DB:      &datastore.DBEnv{},
		Auth: goparent.Authentication{
			SigningKey: []byte("supersecretsquirrel")},
	}

	serviceHandler := api.Handler{
		Env:                   env,
		UserService:           &datastore.UserService{Env: env},
		UserInvitationService: &datastore.UserInviteService{Env: env},
		FamilyService:         &datastore.FamilyService{Env: env},
		ChildService:          &datastore.ChildService{Env: env},
		FeedingService:        &datastore.FeedingService{Env: env},
		SleepService:          &datastore.SleepService{Env: env},
		WasteService:          &datastore.WasteService{Env: env},
	}

	r := api.BuildAPIRouting(&serviceHandler)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	r.Use(handlers.CORS(originsOk, headersOk, methodsOk))

	log.Printf("starting appengine service...")
	http.Handle("/", r)

	appengine.Main()
}
