package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load();
	if err != nil {
		log.Fatal("Error loading .env file");
	}

	fmt.Println(os.Getenv("APP_DB_USERNAME"), os.Getenv("APP_DB_PASSWORD"));
	
	a := App{};
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_REDIS_ADDR"),
		os.Getenv("APP_REDIS_PASSWORD"),
		os.Getenv("APP_REDIS_DB"),
		os.Getenv("CACHE_TTL"));
	a.Run(":8010");
}