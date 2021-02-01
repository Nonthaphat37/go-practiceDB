package main

import (
	"os"
	"testing"
	"log"
	"fmt"
)

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS users
(
	id TEXT PRIMARY KEY,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL
)`

var a App;

func TestMain(m *testing.M){
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
		fmt.Println("dwadawdwa");
		log.Fatal(err);
	}
}

func clearTable(){
	a.DB.Exec("DELETE FROM users");
}