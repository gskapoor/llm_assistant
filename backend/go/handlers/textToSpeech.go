package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	uberDuckAPIKey = "YOUR_UBERDUCK_API_KEY"
	uberDuckAPIURL = "https://api.uberduck.ai/tts"
)

type TextToSpeechRequest struct {
	Text string `json:"text"`
}

func TextToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	// Read the request body containing the text to be converted to speech
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
	}

	// Create a request object for UberDuck's API
	requestData := TextToSpeechRequest{Text: string(body)}
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		http.Error(w, "Failed to create request JSON", http.StatusInternalServerError)
	}

	// Create an HTTP POST request to UberDuck's API
	req, err := http.NewRequest(http.MethodPost, uberDuckAPIURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		http.Error(w, "Failed to create HTTP request", http.StatusInternalServerError)
	}
	req.Header.Set("Authorization", "Bearer "+uberDuckAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request to UberDuck's API
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to make request to UberDuck's API", http.StatusInternalServerError)
	}
	defer response.Body.Close()

	// Read the response from UberDuck's API
	audioData, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response from UberDuck's API", http.StatusInternalServerError)
	}

	// Set the Content-Type header for the audio response
	w.Header().Set("Content-Type", "audio/mpeg") // Adjust the content type based on the actual response

	// Write the audio data as the response
	_, err = w.Write(audioData)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/text-to-speech", TextToSpeechHandler)
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
