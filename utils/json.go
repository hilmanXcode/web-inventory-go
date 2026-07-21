package utils

import (
	"encoding/json"
	"net/http"
)

func MapStringToJson(data map[string]string, w http.ResponseWriter) []byte {
	jsonMessage, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return jsonMessage
}
