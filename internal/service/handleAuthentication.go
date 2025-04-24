package service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sajitha-tj/go-sts/config"
)

type AuthenticationData struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
	Nonce        string
}

func handleAuthentication(w http.ResponseWriter, r *http.Request, data AuthenticationData) bool {
	if r.Method == http.MethodPost {
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if authenticateUser(username, password) {
			return true
		} else {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return false
		}
	}
	// If the method is not POST, render the login page
	renderLoginPage(w, data)
	return false
}

// authenticateUser checks the provided username and password against the stored credentials.
func authenticateUser(username string, password string) bool {
	if username == "peter" && password == "secret" {
		return true
	}
	return false
}

// renderLoginPage renders the login page with the provided authentication data.
func renderLoginPage(w http.ResponseWriter, data AuthenticationData) {
	tmpl, err := template.ParseFiles(config.TEMPLATES_DIR + "/login.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
