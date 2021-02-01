package main

import (
	"os"
	"testing"
	"log"

	"fmt"
	"encoding/json"

	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
	id TEXT PRIMARY KEY,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL
)`

var a App;

func TestMain(m *testing.M){
	err := godotenv.Load();
	if err != nil {
		log.Fatal("Error loading .env file");
	}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_REDIS_ADDR"),
		os.Getenv("APP_REDIS_PASSWORD"),
		os.Getenv("APP_REDIS_DB"),
		os.Getenv("CACHE_TTL"));

	ensureTableExist();
	code := m.Run();
	clearTable();
	os.Exit(code);
}

func ensureTableExist(){
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		fmt.Println(err);
		log.Fatal(err);
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder{
	rr := httptest.NewRecorder();
	a.Router.ServeHTTP(rr, req);
	return rr;
}

func checkResponseCode(t *testing.T, expected, actual int){
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTable(){
	a.DB.Exec("DELETE FROM users");
}

func TestGetNonExistUser(t *testing.T){
	clearTable();

	req, _ := http.NewRequest("GET", "/user/1", nil);
	response := executeRequest(req);
	
	checkResponseCode(t, http.StatusNotFound, response.Code);

	var m map[string]string;
	json.Unmarshal(response.Body.Bytes(), &m);
	if m["error"] != "Not Found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Not Found'. Got '%s'", m["error"]);
	}
}