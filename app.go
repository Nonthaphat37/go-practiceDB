package main

import (
	"log"
	"fmt"
	"strconv"

	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/gorilla/mux"
  _ "github.com/lib/pq"
)

type App struct{
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) initializeRoutes(){
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET");
}

func (a *App) Initialize(user, password, dbName string){
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName);
	var err error;
	a.DB, err = sql.Open("postgres", connectionString);
	if err != nil {
		log.Fatal(err);
	}
	
	a.Router = mux.NewRouter();

	a.initializeRoutes();
	fmt.Println(user, password, dbName);
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r);
    fmt.Println("asbsaf")
	id, err := strconv.Atoi(vars["id"]);
	if err != nil {
		fmt.Println(err);
		respondWithError(w, http.StatusBadRequest, "Invalid User ID");
        return
	}
	
	u := user{id: id};

	if err := u.getUser(a.DB); err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Not Found");
		} else{
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	} else{
		respondWithJSON(w, http.StatusOK, u)
	}
}


func (a *App) Run(addr string){
	log.Fatal(http.ListenAndServe(addr, a.Router));
}