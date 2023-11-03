package main

import (
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
		log.Fatal("Error loading environment variables")
	}

	db := db.Init()

	handler := router.NewDB(db)

	router.InitializeRouter(handler)

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), router.Router))

}
