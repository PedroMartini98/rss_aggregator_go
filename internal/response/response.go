package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func WithJson(w http.ResponseWriter, code int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error: %v ,Marshaling data: %s", err, payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func WithError(w http.ResponseWriter, code int, err string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", err)
	}

	// Me disseram que o struct é mais best pratices so...

	type errResponse struct {
		Error string `json:"error"`
	}

	WithJson(w, code, errResponse{Error: err})
}
