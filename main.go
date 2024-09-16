package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"verve-task/controllers"
	"verve-task/infrastructure"
	"verve-task/services"
)

var mu sync.Mutex

func main() {
	ctx := context.Background()
	// redis
	redisClient, err := infrastructure.NewRedis(ctx)

	if err != nil {
		log.Fatal("error connecting to Redis:", err)
		return
	}
	// kafka
	kafka, err := infrastructure.NewKafka(ctx, "counter")
	if err != nil {
		log.Fatal("error connecting to Kafka:", err)
		return
	}
	kafkaProducer := services.NewKafkaProducer(kafka.Conn)
	// app
	memory := services.NewMemoryStore(redisClient, ctx)
	handler := controllers.NewHttpController(memory, kafkaProducer)
	http.HandleFunc("/api/verve/accept", handler.HandleRequest)
	go memory.LogCounter()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
