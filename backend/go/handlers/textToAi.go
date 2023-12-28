package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type assistantSession struct {
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
}
type assistantSessionInfo struct {
	AssistantSession assistantSession `json:"assistant_session"`
}

type assistantMessage struct {
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
	Message     string `json:"message"`
}

type messageResponse struct {
	Response string `json:"response"`
}

func assistantInit(url string) (assistantSession, error) {

	var session assistantSession

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error accessing API: %v", err)
		return session, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		// TODO: Make HTTP Errors a Library, this is too common a pattern
		return session, fmt.Errorf("request failed with status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return session, err
	}

	var sessionHolder assistantSessionInfo
	err = json.Unmarshal(body, &sessionHolder)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return session, err
	}

	session = sessionHolder.AssistantSession

	return session, nil

}

func assistantChat(session assistantSession, message, url string) (string, error) {

	requestForm := assistantMessage{
		AssistantID: session.AssistantID,
		ThreadID:    session.ThreadID,
		Message:     message,
	}

	// Converts the struct into a json object
	jsonRequest, err := json.Marshal(requestForm)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Printf("Error making post request to LLM url: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		return "", fmt.Errorf("request failed with status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	var messageRes messageResponse
	err = json.Unmarshal(body, &messageRes)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return "", err
	}

	return messageRes.Response, nil

}

func assistantKill(session assistantSession, url string) error {

	jsonSession, err := json.Marshal(session)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonSession))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %v", resp.StatusCode)
	}

	return nil
}

func textToAi(message string) (string, error) {

	baseUrl := os.Getenv("LLM_URL")

	url := baseUrl + "/assistant"

	session, err := assistantInit(url)
	defer assistantKill(session, url)
	if err != nil {
		log.Printf("Error initializing assistant: %v", err)
		return "", err
	}

	response, err := assistantChat(session, message, url)
	if err != nil {
		log.Printf("Error chatting: %v", err)
		return "", err
	}

	return response, nil
}

func HandleTextInput(w http.ResponseWriter, r *http.Request) {

	envLocation := ".env"
	err := godotenv.Load(envLocation)
	if err != nil {
		http.Error(w, "Error reading environment:", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body:", http.StatusInternalServerError)
		return
	}

	str, err := textToAi(string(body))
	if err != nil {
		http.Error(w, "Error submitting request", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, str)

}
