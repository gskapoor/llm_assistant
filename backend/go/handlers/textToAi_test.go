package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
		assert.Equal(t, AssistantSession{}, session, "Expected an empty session on error")
	})

}

func TestAssistantChat(t *testing.T) {
	t.Run("Test Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `{"response":"Test Response"}`)
		}))
		defer server.Close()

		session := AssistantSession{
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

		session := AssistantSession{
			AssistantID: "test_id",
			ThreadID:    "test_thread",
		}

		res, err := assistantChat(session, "", server.URL)

		assert.Error(t, err, "Expected an error")
		assert.Equal(t, "", res, "Expected an empty response on error")
	})
}
