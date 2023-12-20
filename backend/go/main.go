package main

import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gskapoor/llm_assistant/backend/go/handlers"
  "github.com/gskapoor/llm_assistant/backend/go/shared"
)

func textHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprint(w, "BOTTOM TEXT\n")
}

func main () {

  r := mux.NewRouter()

  r.HandleFunc("/voice", middleware.Logging()(handlers.HandleVoiceInput) )
  r.HandleFunc("/text", middleware.Logging()(handlers.HandleTextInput) )

  // Form w/ Audio -> Text -> an AI response 

  fmt.Println("Server is running on :8080")
  err := http.ListenAndServe(":8080", r)
  if err != nil {
    fmt.Println(err)
  }
}
