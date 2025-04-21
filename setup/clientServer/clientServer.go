package main

// This is a sample client server that handles the OAuth2 authorization code flow.
// It listens for a callback from the authorization server and handles the redirect.
// This server is not a complete implementation of an OAuth2 client, but rather used for testing purposes.
// It is meant to be run alongside with the STS.

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth callback received")
		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "application/json")

		responseBodyData := make(map[string]string)
		// error handling
		if r.URL.Query().Get("error") != "" {
			// send json error response
			responseBodyData = map[string]string{
				"error":             r.URL.Query().Get("error"),
				"error_description": r.URL.Query().Get("error_description"),
			}
		} else {
			// success response: send all URL query params as JSON
			for key, values := range r.URL.Query() {
				if len(values) > 0 {
					responseBodyData[key] = values[0]
				}
			}
		}

		jsonData, err := json.Marshal(responseBodyData)
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})

	fmt.Println("Starting server on port 3846...")
	err := http.ListenAndServe(":3846", nil)
	if err != nil {
		panic(err)
	}
}
