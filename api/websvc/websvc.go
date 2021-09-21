package websvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Orchestrator interface {
	StartProcess(requestID string) error
	StopProcess(requestID string) error
}

// NewCfg - creates new config for web service.
func NewCfg(addr string, o Orchestrator) Cfg {
	return Cfg{addr, o}
}

// Cfg - config for web service.
type Cfg struct {
	address      string
	orchestrator Orchestrator
}

// Init - creates API and start new web server.
func Init(cfg Cfg) {
	api := &api{cfg.orchestrator}
	http.HandleFunc("/", api.handle)

	fmt.Println("web server listening at ", cfg.address)
	log.Fatal(http.ListenAndServe(cfg.address, nil))
}

type api struct {
	orchestrator Orchestrator
}

func (api *api) handle(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		respondWithJSON(w, "Error", "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var payload request
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&payload)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			respondWithJSON(w, "Error", "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			respondWithJSON(w, "Error", "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	switch payload.Action {
	case "start":
		fmt.Println("start")
		if err := api.orchestrator.StartProcess(payload.RequestID); err != nil {
			respondWithJSON(w, "Error", err.Error(), http.StatusOK)
			return
		}
		respondWithJSON(w, "Status", "Success", http.StatusOK)
	case "stop":
		fmt.Println("stop")
		if err := api.orchestrator.StopProcess(payload.RequestID); err != nil {
			respondWithJSON(w, "Error", err.Error(), http.StatusOK)
			return
		}
		respondWithJSON(w, "Status", "Success", http.StatusOK)
	}

	return
}

func respondWithJSON(w http.ResponseWriter, content string, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	payload := map[string]string{content: message}
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
