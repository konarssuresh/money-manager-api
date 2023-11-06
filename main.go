package main

import (
	"fmt"
	"log"
	"money-manager/db"
	"money-manager/router"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading environment variables")
	}

	db := db.Init()

	handler := router.NewDB(db)

	router.InitializeRouter(handler)

	port := os.Getenv("HOST_PORT")
	if port == "" {
		port = ":8080" // Default to 8080 if the environment variable is not set
	}

	er := http.ListenAndServe(port, router.Router)

	if er != nil {
		fmt.Println("error starting server")
		fmt.Println("port")
		fmt.Println(port)
		fmt.Println("server start error")
		fmt.Println(er)
		log.Fatal(er)
	}

}
