package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type Dialogue struct {
	TranscribedText string `json:"transcribed_text"`
	Response        string `json:"response"`
}

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

// HandleVoiceInput: a HTTP handler for Speech to Text
func HandleVoiceInput(w http.ResponseWriter, r *http.Request) {

	const storageDir = "uploads"

	// File size is 20 megabytes, so 20 << 20 bytes
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Println("Failed to parse form: ", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("audio")
	if err != nil {
		log.Println("Failed to find audio: ", err)
		http.Error(w, "Failed to find audio", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = createTempDirectory(storageDir)
	if err != nil {
		log.Println("Failed to create temporary directory: ", err)
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	timestamp := time.Now().UnixNano()

	audioFilePath := fmt.Sprintf(storageDir+"/audiofile_%d.wav", timestamp)

	out, err := os.Create(audioFilePath)
	if err != nil {
		log.Println("Failed to create temporary file: ", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer deleteFile(audioFilePath)
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("Failed to save audio file: ", err)
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	transcribedText, err := transcribe(audioFilePath)

	if err != nil {
		log.Println("Failed to transcribe audio file: ", err)
		http.Error(w, "Failed to transcribe audio file", http.StatusInternalServerError)
		return
	}

	// TODO: Make this call the endpoint instead, this is easier for now
	response, err := textToAi(transcribedText)
	if err != nil {
		log.Println("Error reaching AI: ", err)
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	dialogueStruct := Dialogue{
		TranscribedText: transcribedText,
		Response:        response,
	}

	jsonDialogue, err := json.Marshal(dialogueStruct)
	if err != nil {
		log.Println("Error Marshaling JSON: ", err)
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonDialogue))

}
