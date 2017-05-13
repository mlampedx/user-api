package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type API struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *API) Initialize(user, password, dbname string) {
	connectionStr := fmt.Sprintf("Attempting to connect to dbname=%s with user=%s and password=%s", dbname, user, password)
	var err error
	a.DB, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
}

func (a *API) Run(port string) {

}
