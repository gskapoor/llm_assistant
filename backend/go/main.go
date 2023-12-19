package main

import (
  "fmt"
  "net/http"
)

func textHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprint(w, "BOTTOM TEXT\n")
}

func main () {
  http.HandleFunc("/audio_query", textHandler)

  // Form w/ Audio -> Text -> an AI response 

  fmt.Println("Server is running on :8080")
  http.ListenAndServe(":8080", nil)
}
