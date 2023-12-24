package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type AssistantSession struct {
	AssistantSession struct {
		AssistantID string `json:"assistant_id"`
		ThreadID    string `json:"thread_id"`
	} `json:"assistant_session"`
}

type AssistantMessage struct {
	AssistantID string `json:"assistant_id"`
	ThreadID    string `json:"thread_id"`
	Message     string `json:"message"`
}

type MessageResponse struct {
	Response string `json:"response"`
}

func assistantInit() (AssistantSession, error) {

	var session AssistantSession

	resp, err := http.Get("http://localhost:8000/assistant")
	if err != nil {
		log.Printf("Error accessing API: %v", err)
		return session, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		return session, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return session, err
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return session, err
	}

	return session, nil

}

func assistantChat(session AssistantSession, message string) (string, error) {
	assistantID := session.AssistantSession.AssistantID
	threadID := session.AssistantSession.ThreadID

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

	// TODO: keep the URL somewhere else/pass in a parameter
	resp, err := http.Post("http://localhost:8000/assistant", "application/json", bytes.NewBuffer(jsonRequest))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status code: %v", resp.StatusCode)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
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

func assistantKill(session AssistantSession) error {

	jsonSession, err := json.Marshal(session)

	_, err = http.NewRequest("DELETE", "http://localhost:8000/assistant", bytes.NewBuffer(jsonSession))
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return err
	}

	return nil
}

func textToAi(message string) (string, error) {
	session, err := assistantInit()
	defer assistantKill(session)

	if err != nil {
		log.Printf("Error initializing assistant: %v", err)
	}
	response, err := assistantChat(session, message)

	return response, nil
}

func HandleTextInput(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body:", http.StatusInternalServerError)
		return
	}


	str, err := textToAi(string(body))
	if err != nil {
		http.Error(w, "Error submitting request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, str)

}
