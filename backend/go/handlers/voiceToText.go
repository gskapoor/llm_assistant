package handlers

import (
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

func getOpenAIKey() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error reading environment: ", err)
		return "", err
	}

	key := os.Getenv("OPENAI_API_KEY")

	return key, nil
}

func makeSession(audioPath string) (string, error) {

	token, err := getOpenAIKey()
	if err != nil {
		return "", err
	}

	trClient := openai.NewClient(token)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioPath,
	}

	resp, err := trClient.CreateTranscription(ctx, req)
	if err != nil {
		log.Printf("Error getting transcription", err)
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

func createTempDirectory(filePath string) error {
	const permissions = 0700
	var err error = os.MkdirAll(filePath, permissions)
	return err
}

// HandleVoiceInput: a HTTP handler for Speech to Text
func HandleVoiceInput(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(20 << 20)

	log.Println("Recieved form")

	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	log.Println("Parsed form")

	file, _, err := r.FormFile("audio")
	defer file.Close()

	log.Println("Opened Audio File")

	createTempDirectory("uploads")

	timestamp := time.Now().UnixNano()

	audioFilePath := fmt.Sprintf("uploads/audiofile_%d.wav", timestamp)

	out, err := os.Create(audioFilePath)
	if err != nil {
		log.Println("Failed to create temporary file:", err)
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer deleteFile(audioFilePath)
	defer out.Close()

	log.Println("Created File")

	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("Failed to save audio file", err)
		http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
		return
	}

	log.Println("Copied File")

	transcribedText, err := makeSession(audioFilePath)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(transcribedText))

	log.Println("Sent to API")

}
