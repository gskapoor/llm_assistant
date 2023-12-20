package handlers

import (
  "fmt"
  "log"
  "net/http"
)

func HandleVoiceInput(w http.ResponseWriter, r *http.Request){

  fmt.Fprintf(w, "Handling voice input\n") 
  r.ParseMultipartForm(20 << 20)
  // if err != nil {
  //   http.Error(w, "Unable to parse form", http.StatusBadRequest)
  // }
  // file, _, err := r.FormFile("audio")
  // defer file.Close()

  log.Println("Sent to API")

}

