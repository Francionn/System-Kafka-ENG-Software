package consumer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendToAPI(data SensorData) {
	apiURL := "http://localhost:3000/sensor" // endpoint API Fiber

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[API] Error converting JSON: %v", err)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[API] Erro send to API: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[API] Unexpected response: %s", resp.Status)
	}
}
