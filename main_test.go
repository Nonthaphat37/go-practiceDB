// this program is not cover all the cases;

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
	"bytes"
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

func TestCreateUser(t *testing.T){
	clearTable();
	var jsonStr = []byte(`{"id" : "00001", "firstname" : "nonthaphat", "lastname" : "wongwattanakij"}`);
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr));
	req.Header.Set("Content-Type", "application/json");

	response := executeRequest(req);
	checkResponseCode(t, http.StatusCreated, response.Code);

	var m map[string]interface{};
	json.Unmarshal(response.Body.Bytes(), &m);
	
	if m["id"] != "00001" {
		t.Errorf("Expected id name to be '00001'. Got ‘%v’", m["id"]);
	}

	if m["firstname"] != "nonthaphat" {
		t.Errorf("Expected firstname to be 'nonthaphat'. Got ‘%v’", m["firstname"]);
	}

	if m["lastname"] != "wongwattanakij" {
		t.Errorf("Expected lastname to be 'wongwattanakij'. Got ‘%v’", m["lastname"]);
	}
}

func addUser(id, firstname, lastname string){
	a.DB.Exec("INSERT INTO users(id, firstname, lastname) VALUES($1, $2, $3) RETURNING id", id, firstname, lastname);
}

func TestGetUser(t *testing.T){
	clearTable();
	addUser("00001", "nonthaphat", "wongwattanakij");

	req, _ := http.NewRequest("GET", "/user/1", nil);
	response := executeRequest(req);

	checkResponseCode(t, http.StatusOK, response.Code);
}

func TestUpdateUser(t *testing.T){
	clearTable();
	addUser("00001", "nonthaphat", "wongwattanakij");

	req, _ := http.NewRequest("GET", "/user/1", nil);
	response := executeRequest(req);

	var originalUser map[string]interface{};
	json.Unmarshal(response.Body.Bytes(), &originalUser);

	var jsonStr = []byte(`{"id" : "00001", "firstname" : "nonthaphat2", "lastname" : "wongwattanakij2"}`);
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr));
	req.Header.Set("Content-Type", "application/json");

	response = executeRequest(req);

	checkResponseCode(t, http.StatusOK, response.Code);

	var m map[string]interface{};
	json.Unmarshal(response.Body.Bytes(), &m);

	if m["id"] != originalUser["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"]);
	}

	if m["firstname"] == originalUser["firstname"] {
        t.Errorf("Expected the firstname to change from '%v' to '%v'. Got '%v'", originalUser["firstname"], m["firstname"], m["firstname"]);
    }

    if m["lastname"] == originalUser["lastname"] {
        t.Errorf("Expected the lastname to change from '%v' to '%v'. Got '%v'", originalUser["lastname"], m["lastname"], m["lastname"]);
    }
}