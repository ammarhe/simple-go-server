package usecases

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"sync"
//	"verve-task/repositories"
//)
//
//type RequestHandler struct {
//	memory *repositories.MemoryStore
//}
//
//func (h *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request, mu *sync.Mutex) {
//	idParam := r.URL.Query().Get("id")
//	endpoint := r.URL.Query().Get("endpoint")
//	title := r.URL.Query().Get("title") // change to struct
//
//	w.Header().Set("Content-Type", "application/json")
//	if idParam == "" {
//		http.Error(w, "id is required", http.StatusBadRequest)
//		return
//	}
//	if !h.memory.IsLogged(idParam) {
//		go h.memory.LogRequest(idParam)
//	}
//
//	wg := sync.WaitGroup{}
//	var response *http.Response
//	var err error
//	if endpoint != "" {
//
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			endpointErr, resp := callEndpoint(endpoint, title)
//			if endpointErr != nil {
//				err = endpointErr
//			}
//			response = resp
//		}()
//	}
//	wg.Wait()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	if response != nil {
//		defer response.Body.Close()
//		respBody, readErr := io.ReadAll(response.Body)
//		if readErr != nil {
//			http.Error(w, readErr.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		_, writeErr := w.Write(respBody)
//		if writeErr != nil {
//			http.Error(w, writeErr.Error(), http.StatusInternalServerError)
//			return
//		}
//		return
//	}
//
//	_, err = w.Write([]byte("Ok"))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	return
//}
//
//func callEndpoint(endpoint string, title string) (error, *http.Response) {
//	requestBody := map[string]string{"title": title}
//
//	jsonData, err := json.Marshal(requestBody)
//	if err != nil {
//		return err, nil
//	}
//	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
//	if err != nil {
//		return err, nil
//	}
//	fmt.Println(resp)
//	return nil, resp
//}
