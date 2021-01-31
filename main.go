package main

import (
	// "os"
)

func main(){
	a := App{};
	a.Initialize("postgres", "12345", "postgres")
	a.Run(":8010");
}