package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/datastore"
	"google.golang.org/appengine"
)

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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Authorization", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	http.Handle("/", handlers.CORS(originsOk, headersOk, methodsOk)(r)) //wire mux to DefaultServeMux for appengine
	appengine.Main()
}
