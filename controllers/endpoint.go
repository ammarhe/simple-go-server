package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"verve-task/repositories"
)

type HttpController struct {
	MemoryRepo *repositories.MemoryStore
}

func NewHttpController(memoryRepo *repositories.MemoryStore) *HttpController {
	return &HttpController{
		MemoryRepo: memoryRepo,
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
	if !h.MemoryRepo.IsLogged(id) {
		h.MemoryRepo.LogRequest(id)
	}

	// calling endpoint if provided
	if endpoint != "" {
		response, err := h.callEndpoint(endpoint, title)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		respBody, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, "failed", http.StatusInternalServerError)
			return
		}
		w.Write(respBody)
		return
	}

	w.Write([]byte("Ok"))
}

func (h *HttpController) callEndpoint(endpoint string, title string) (*http.Response, error) {
	requestBody := map[string]string{"title": title, "count": strconv.Itoa(len(h.MemoryRepo.LoggedReqIds))}

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
