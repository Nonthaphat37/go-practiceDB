package main

import (
	"os"
)

func main(){
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