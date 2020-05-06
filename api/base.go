package api

import (
	"encoding/json"
	"net/http"

	"github.com/Z-M-Huang/Tools/data"
)

func writeResponse(w http.ResponseWriter, response *data.Response) {
	jsonBody, _ := json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}

func writeUnexpectedError(w http.ResponseWriter, response *data.Response) {
	response.Alert.IsDanger = true
	response.Alert.Message = "Um... Your data got eaten by the cyber space... Would you like to try again?"
	writeResponse(w, response)
}
