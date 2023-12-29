package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssistantInit(t *testing.T) {
	t.Run("Test Success", func(t *testing.T) {
		// Create a mock HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"assistant_session":{"assistant_id":"test_id","thread_id":"test_thread"}}`)
		}))
		defer server.Close()

		// Call the function being tested with the mock server URL
		session, err := assistantInit(server.URL)
		assert.NoError(t, err, "Expected no error")

		// Verify the result
		assert.Equal(t, "test_id", session.AssistantID, "Unexpected assistant ID")
		assert.Equal(t, "test_thread", session.ThreadID, "Unexpected thread ID")
	})

	t.Run("HTTP Error", func(t *testing.T) {
		// Create a mock HTTP server that returns an error status code
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer server.Close()

		session, err := assistantInit(server.URL)

		assert.Error(t, err, "Expected an error")
		assert.Equal(t, assistantSession{}, session, "Expected an empty session on error")
	})

}

func TestAssistantChat(t *testing.T) {
	t.Run("Test Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"response":"Test Response"}`)
		}))
		defer server.Close()

		session := assistantSession{
			AssistantID: "test_id",
			ThreadID:    "test_thread",
		}

		res, err := assistantChat(session, "", server.URL)
		assert.NoError(t, err, "Expected no error")

		assert.Equal(t, res, "Test Response", "Unexpected Response")

	})

	t.Run("Test HTTP Failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer server.Close()

		session := assistantSession{
			AssistantID: "test_id",
			ThreadID:    "test_thread",
		}

		res, err := assistantChat(session, "", server.URL)

		assert.Error(t, err, "Expected an error")
		assert.Equal(t, "", res, "Expected an empty response on error")
	})
}

func TestAssistantKill(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"response":"Test Response"}`)
		}))
		defer server.Close()

		session := assistantSession{
			AssistantID: "test_id",
			ThreadID:    "test_thread",
		}

		err := assistantKill(session, server.URL)
		assert.Nil(t, err, "Expected no error")
	})

	t.Run("HTTP Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer server.Close()

		session := assistantSession{
			AssistantID: "test_id",
			ThreadID:    "test_thread",
		}

		err := assistantKill(session, server.URL)
		assert.Error(t, err, "Expected an error")
	})

}

func TestTextToAi(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		mockResponse := `{"response": "MockedResponse"}`

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(mockResponse))
		}))
		defer mockServer.Close()

		err := os.Setenv("LLM_URL", mockServer.URL)
		if err != nil {
			t.Fatalf("Error setting environment variable: %v", err)
		}
		defer os.Unsetenv("LLM_URL")

		testMessage := "Test message"
		result, err := textToAi(testMessage)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result != "MockedResponse" {
			t.Errorf("Expected 'MockedResponse', got '%s'", result)
		}
	})

	t.Run("Failure to Initialize", func(t *testing.T) {

		mockResponse := `{"response": "MockedResponse"}`

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				w.Write([]byte(mockResponse))
			}
		}))
		defer mockServer.Close()

		err := os.Setenv("LLM_URL", mockServer.URL)
		if err != nil {
			t.Fatalf("Error setting environment variable: %v", err)
		}
		defer os.Unsetenv("LLM_URL")

		testMessage := "Test message"
		_, err = textToAi(testMessage)

		if err == nil {
			t.Errorf("Expected an error, but got nil")
		}
		// TODO: Write test for correct HTTP Error, not doable until Error library is made
	})
	t.Run("Failed to Chat", func(t *testing.T) {

		mockResponse := `{"response": "MockedResponse"}`

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				w.Write([]byte(mockResponse))
			}
		}))
		defer mockServer.Close()

		err := os.Setenv("LLM_URL", mockServer.URL)
		if err != nil {
			t.Fatalf("Error setting environment variable: %v", err)
		}
		defer os.Unsetenv("LLM_URL")

		testMessage := "Test message"
		_, err = textToAi(testMessage)

		if err == nil {
			t.Errorf("Expected an error, but got nil")
		}
	})

	t.Run("Failed to Kill", func(t *testing.T) {

		mockResponse := `{"response": "MockedResponse"}`

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "DELETE" {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				w.Write([]byte(mockResponse))
			}
		}))
		defer mockServer.Close()

		err := os.Setenv("LLM_URL", mockServer.URL)
		if err != nil {
			t.Fatalf("Error setting environment variable: %v", err)
		}
		defer os.Unsetenv("LLM_URL")

		testMessage := "Test message"
		result, err := textToAi(testMessage)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result != "MockedResponse" {
			t.Errorf("Expected 'MockedResponse', got '%s'", result)
		}
	})
}
