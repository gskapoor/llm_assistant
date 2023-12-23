package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// getOpenAIKey: Gets the environment variable OPENAI_API_KEY
func getOpenAIKey() (string, error) {

	const envLocation = ".env"
	const envVarName = "OPENAI_API_KEY"

	err := godotenv.Load(envLocation)
	if err != nil {
		log.Println("Error reading environment: ", err)
		return "", err
	}

	key := os.Getenv(envVarName)

	return key, nil
}

// transcribe: Turns an audio file, at a given path, into text
// audioPath (string): The path to the audio file
func transcribe(audioPath string) (string, error) {

	token, err := getOpenAIKey()
	if err != nil {
		return "", err
	}

	// TODO: Understand this code and write something explaining it
	// This is basically just boilerplate given by the "openai" library
	trClient := openai.NewClient(token)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioPath,
	}

	resp, err := trClient.CreateTranscription(ctx, req)
	if err != nil {
		log.Printf("Error getting transcription: %v", err)
		return "", err
	}

	return resp.Text, nil
}

// deleteFile: Deletes file at a given path
// filePath (string): The path to the file
func deleteFile(filePath string) {

	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to delete the file: %v", err)
	}
}

// createTempDirectory: Creates a temporary directory at a given path
// filePath (string): The path to the file
func createTempDirectory(filePath string) error {

	// These are UNIX, representing an octal number
	// Read here for reference: https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation
	const permissions = 0700

	err := os.MkdirAll(filePath, permissions)
	return err
}

// textToAI: Takes a transcribed message and outputs an AI response
func textToAI(message string) (string, error) {

	// TODO: Implement this with text endpoint
	_, err := http.NewRequest("POST", "localhost:8000/text", bytes.NewBufferString(message))
	if err != nil {
		return "", err
	}

	return "", nil
}

// HandleVoiceInput: a HTTP handler for Speech to Text
func HandleVoiceInput(w http.ResponseWriter, r *http.Request) {

	const storageDir = "uploads"

	// File size is 20 megabytes, so 20 << 20 bytes
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("audio")
	defer file.Close()

	err = createTempDirectory(storageDir)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().UnixNano()

	audioFilePath := fmt.Sprintf(storageDir+"/audiofile_%d.wav", timestamp)

	out, err := os.Create(audioFilePath)
	if err != nil {
		log.Println("Failed to create temporary file:", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer deleteFile(audioFilePath)
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("Failed to save audio file", err)
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	transcribedText, err := transcribe(audioFilePath)

	if err != nil {
		log.Println("ERROR")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(transcribedText))

	// TODO: Send a response from the text endpoint

}
