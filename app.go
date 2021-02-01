package main

import (
	"time"
	"log"
	"fmt"
	"strconv"
	"strings"

	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/gorilla/mux"
  _ "github.com/lib/pq"
	"github.com/go-redis/redis"
)

type App struct{
	Router    *mux.Router
	DB        *sql.DB
	Redis     *redis.Client
	cache_ttl time.Duration 
}

func (a *App) initializeRoutes(){
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET");
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
}

func (a *App) Initialize(dbUser, dbPassword, dbName, redisAddr, redisPassword, redisDB, cache_ttl string){
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName);
	var err error;
	a.DB, err = sql.Open("postgres", connectionString);
	if err != nil {
		log.Fatal(err);
	}

	a.Router = mux.NewRouter();

	redisDB2, _ := strconv.Atoi(redisDB);
	a.Redis = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: redisPassword,
		DB: redisDB2,
	})

	ttl, _ := strconv.Atoi(cache_ttl);
	a.cache_ttl = time.Duration(ttl) * time.Millisecond;

	pong, err := a.Redis.Ping().Result();
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

	id := fmt.Sprintf("%05s", vars["id"]);
	fmt.Println(id);
	
	// id, err := strconv.Atoi(vars["id"]);
	// if err != nil {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid User ID");
	// 	return;
	// }
	defer r.Body.Close()

	u := user{ID: id};

	fmt.Println("Get Method");
	if val, err := u.getUserRedis(a.Redis); err != redis.Nil {
		fmt.Println("Get user from redis");
		data := user{};
		json.Unmarshal([]byte(val), &data);
		respondWithJSON(w, http.StatusOK, data);
	} else{
		if err := u.getUserDB(a.DB); err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Not Found");
				respondWithError(w, http.StatusNotFound, "Not Found");
			} else{
				respondWithError(w, http.StatusInternalServerError, err.Error());
			}
		} else{
			fmt.Println("Get user from DB");
			respondWithJSON(w, http.StatusOK, u);

			fmt.Println("Set user to redis");
			err := u.setUserRedis(a.Redis, a.cache_ttl);
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
		fmt.Println(err);
		respondWithError(w, http.StatusBadRequest, "Invalid request payload.");
		return;
	}
	defer r.Body.Close();

	if(u.Firstname == "" || u.Lastname == ""){
		fmt.Println("Invalid request payload: Name must be not empty.");
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: Name must be not empty.");
		return;
	}

	u.Firstname = strings.TrimSpace(u.Firstname);
	u.Lastname = strings.TrimSpace(u.Lastname);
	u.ID = fmt.Sprintf("%05s", u.ID);
	tmp := u;

	fmt.Println("Post Method");
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