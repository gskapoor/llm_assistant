package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/gskapoor/llm_assistant/backend/go/handlers"
	midw "github.com/gskapoor/llm_assistant/backend/go/shared"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/text", midw.Chain(handlers.HandleTextInput, midw.LoggingStart(), midw.Logging()))
	r.HandleFunc("/voice", midw.Chain(handlers.HandleVoiceInput, midw.LoggingStart(), midw.Logging()))
	r.HandleFunc("/tts", midw.Chain(handlers.TextToSpeechHandler, midw.LoggingStart(), midw.Logging()))

	fmt.Println("Server is running on :8080")

	handler := cors.Default().Handler(r)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
	}
}
