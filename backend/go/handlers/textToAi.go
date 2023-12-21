package handlers

import (
	"fmt"
	"net/http"
)

func HandleTextInput(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Sample AI Text\n")
}
