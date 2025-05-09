package util

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {

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

func RespondWithError(w http.ResponseWriter, code int, err string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", err)
	}

	// Me disseram que o struct é mais best pratices so...

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, errResponse{Error: err})
}
