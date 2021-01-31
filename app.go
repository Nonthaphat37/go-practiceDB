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
	"github.com/go-redis/redis"
)

type App struct{
	Router *mux.Router
	DB     *sql.DB
	Redis  *redis.Client
}

func (a *App) initializeRoutes(){
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET");
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
}

func (a *App) Initialize(user, password, dbName string){
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbName);
	var err error;
	a.DB, err = sql.Open("postgres", connectionString);
	if err != nil {
		log.Fatal(err);
	}

	a.Router = mux.NewRouter();

	a.Redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	pong, err := a.Redis.Ping().Result()
	fmt.Println("Test ping redis", pong, err);

	a.initializeRoutes();
}

func respondWithError(w http.ResponseWriter, code int, message string){
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r);
	id, err := strconv.Atoi(vars["id"]);
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID");
		return;
	}
	defer r.Body.Close()

	u := user{ID: id};

	fmt.Println("Get Mathod");
	if val, err := u.getUserRedis(a.Redis); err != redis.Nil {
		data := user{};
		json.Unmarshal([]byte(val), &data);
		fmt.Println("Get user from redis");
		respondWithJSON(w, http.StatusOK, data);
	} else{
		if err := u.getUserDB(a.DB); err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Not Found");
			} else{
				respondWithError(w, http.StatusInternalServerError, err.Error());
			}
		} else{
			fmt.Println("Get user from DB");
			respondWithJSON(w, http.StatusOK, u);

			fmt.Println("Set user to redis");
			err := u.setUserRedis(a.Redis);
			if err != nil {
				fmt.Println("Error set from redis", err);
			}
		}
	}
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request){
	var u user;
	decoder := json.NewDecoder(r.Body);
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload");
		return;
	}
	defer r.Body.Close();

	tmp := u;
	fmt.Println("Post Mathod");
	if err := u.getUserDB(a.DB); err != nil {
		if err == sql.ErrNoRows {
			if err := u.createUser(a.DB); err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error());
			} else{
				respondWithJSON(w, http.StatusCreated, u);
			}
		} else{
			respondWithError(w, http.StatusInternalServerError, err.Error());
		}
	} else{
		u = tmp;
		if err := u.updateUser(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error());
		} else{
			respondWithJSON(w, http.StatusOK, u);
			err := u.delUserRedis(a.Redis);
			if err != nil {
				fmt.Println("Error delete from redis", err);
			}
		}
	}
}

func (a *App) Run(addr string){
	log.Fatal(http.ListenAndServe(addr, a.Router));
}