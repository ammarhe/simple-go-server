package main

import (
	"log"
	"net/http"
	"sync"
	"verve-task/controllers"
	"verve-task/repositories"
)

var mu sync.Mutex

func main() {

	memory := repositories.NewMemoryStore()
	handler := controllers.NewHttpController(memory)
	http.HandleFunc("/api/verve/accept", handler.HandleRequest)
	go memory.LogCounter()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
