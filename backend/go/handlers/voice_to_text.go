package handlers

import (
	"fmt"
	"log"
	"net/http"
  "time"
  "io"
	"os"
)

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
	// defer file.Close()


	log.Println("Opened Audio File")

  createTempDirectory("uploads")

  timestamp := time.Now().UnixNano()

  audioFilePath := fmt.Sprintf("uploads/audiofile_%d.wav", timestamp)

  out, err := os.Create(audioFilePath)
  if err != nil {
    http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
    return
  }
  // defer deleteFile(audioFilePath)
  defer out.Close()

	log.Println("Created File")

  if out == nil {
    log.Println("Out is nil")
  }

  if file == nil {
    log.Println("File is nil")
  }

  _, err = io.Copy(out, file)
  if err != nil {
    log.Println(err)
    log.Println("SUS")
    http.Error(w, "Failed to save audio file", http.StatusInternalServerError)
    return
  }

	log.Println("Copied File")

	log.Println("Sent to API")

}
