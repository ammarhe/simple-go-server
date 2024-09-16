package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"verve-task/services"
)

type HttpController struct {
	MemoryService *services.MemoryStore
	KafkaService  *services.KafkaProducer
}

func NewHttpController(memoryRepo *services.MemoryStore, kafkaService *services.KafkaProducer) *HttpController {
	return &HttpController{
		MemoryService: memoryRepo,
		KafkaService:  kafkaService,
	}
}
func (h *HttpController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	endpoint := r.URL.Query().Get("endpoint")
	title := r.URL.Query().Get("title")

	if idParam == "" {
		http.Error(w, "failed", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
		return
	}

	// Check if the request is logged
	if !h.MemoryService.IsLogged(id) {
		err := h.MemoryService.LogRequest(id)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		// write to streaming service
		err = h.KafkaService.WriteMsg(idParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// calling endpoint if provided
	if endpoint != "" {
		response, err := h.callEndpoint(endpoint, title)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()
		err = h.MemoryService.OpenLogFile()
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		log.Printf("Endpoint response status code %v", response.StatusCode)
		err = h.MemoryService.CloseLogFile()
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Ok"))
}

func (h *HttpController) callEndpoint(endpoint string, title string) (*http.Response, error) {
	requestBody := map[string]string{"title": title, "count": strconv.Itoa(len(h.MemoryService.LoggedReqIds))}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	return resp, nil
}
