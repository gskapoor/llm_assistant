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

type AssistantSession struct {
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
}
type AssistantSessionInfo struct {
	AssistantSession AssistantSession `json:"assistant_session"`
}

type AssistantMessage struct {
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
	Message     string `json:"message"`
}

type MessageResponse struct {
	Response string `json:"response"`
}

func getLLMUrl() (string, error) {
	const envLocation = ".env"
	const envVarName = "LLM_URL"

	err := godotenv.Load(envLocation)
	if err != nil {
		log.Println("Error reading environment: ", err)
		return "", err
	}

	key := os.Getenv(envVarName)

	return key, nil
}

func assistantInit(url string) (AssistantSession, error) {

	var session AssistantSession

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error accessing API: %v", err)
		return session, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		return session, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return session, err
	}

	var sessionHolder AssistantSessionInfo
	err = json.Unmarshal(body, &sessionHolder)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return session, err
	}

	session = sessionHolder.AssistantSession

	return session, nil

}

func assistantChat(session AssistantSession, message, url string) (string, error) {
	assistantID := session.AssistantID
	threadID := session.ThreadID

	requestForm := AssistantMessage{
		AssistantID: assistantID,
		ThreadID:    threadID,
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
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	var messageRes MessageResponse
	err = json.Unmarshal(body, &messageRes)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return "", err
	}

	return messageRes.Response, nil

}

func assistantKill(session AssistantSession, url string) error {

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

	_, err = client.Do(req)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return err
	}

	return nil
}

func textToAi(message string) (string, error) {

	base_url, err := getLLMUrl()
	if err != nil {
		log.Printf("Error getting url, make sure to set .env variable: %v", err)
		return "", err
	}

	url := base_url + "/assistant"

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
