package main

// This is a sample client server that helps to handle the OAuth2 authorization code flow.
// It listens for a callback from the authorization server and handles the redirect.
// This server is not a complete implementation of an OAuth2 client, but rather used for testing purposes.
// It is meant to be run alongside with the STS.

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const CLIENT_ID = "my-client"
const CLIENT_SECRET = "foobar"

var authCode string

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/callback", handleCallback).Methods("GET")
	r.HandleFunc("/getToken", handleGetToken).Methods("GET")

	log.Println("Starting client on port: 3846")
	if err := http.ListenAndServe(":3846", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Render the index page
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<title>OAuth2 Client</title>
		</head>
		<body>
			<h1>Welcome to the OAuth2 Client</h1>
			<a href="http://123e4567-e89b-12d3-a456-426614174000.localhost:8080/authorize?response_type=code&client_id=my-client&redirect_uri=http://localhost:3846/callback&scope=fosite+openid+photos+offline&state=random-state-value&nonce=random-nonce-value&code_challenge=example-code-challenge&code_challenge_method=S256">Authorize</a>
		</body>
		</html>
	`))
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	responseBodyData := make(map[string]string)

	if r.URL.Query().Get("error") != "" {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusOK)
		authCode = r.URL.Query().Get("code")
	}

	for key, value := range r.URL.Query() {
		responseBodyData[key] = value[0]
	}
	jsonData, err := json.Marshal(responseBodyData)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<title>OAuth2 Client Callback</title>
		</head>
		<body>
			<h1>Callback Response</h1>
			<pre>` + string(jsonData) + `</pre>
			<br/>
			<a href="/getToken">Get Token</a>
		</body>
		</html>
	`))
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	if authCode == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	tokenURL := "http://123e4567-e89b-12d3-a456-426614174000.localhost:8080/token"

	data := url.Values{}
	data.Set("code", authCode)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", CLIENT_ID)
	data.Set("client_secret", CLIENT_SECRET)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: received status code %d", resp.StatusCode)
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
