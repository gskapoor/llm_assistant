package handlers

import (
  "fmt"
  "net/http"
)

func HandleVoiceInput(w http.ResponseWriter, r *http.Request){

  fmt.Fprintf(w, "Handling voice input\n") 

}

