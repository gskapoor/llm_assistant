package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func TextToSpeechHandler(w http.ResponseWriter, r *http.Request) {
	openAIURL := os.Getenv("TTS_URL")
	apiKey := os.Getenv("TTS_API_KEY")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define the input text for TTS
	inputText := "Today is a wonderful day to build something people love!"

	// Prepare the request body
	requestBody := fmt.Sprintf(`{"model":"tts-1", "voice":"alloy", "input":"%s"}`, inputText)
	// Create HTTP POST request to OpenAI TTS API
	req, err := http.NewRequest("POST", openAIURL, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		log.Printf("HTTP error: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body (audio content)
	audioContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Write audio content to file (speech.mp3)
	err = os.WriteFile("speech.mp3", audioContent, 0644)
	if err != nil {
		log.Printf("HTTP error: %v", err)
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Speech saved to speech.mp3")
}
