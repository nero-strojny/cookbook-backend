package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"server/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS()(r)))
}
