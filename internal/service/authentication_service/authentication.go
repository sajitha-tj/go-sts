package authentication_service

import (
	"html/template"
	"log"
	"net/http"

	"github.com/sajitha-tj/go-sts/internal/repository/user_repository"
)

const TEMPLATES_DIR = "/home/sajithaj/my-sts-project/go-sts/internal/service/authentication_service/templates"

type AuthenticationData struct {
	ResponseType string
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
}

type AuthenticationService struct {
	userStore *user_repository.UserStore
}

func NewAuthenticationService(userStore *user_repository.UserStore) *AuthenticationService {
	return &AuthenticationService{
		userStore: userStore,
	}
}

func (s *AuthenticationService) HandleAuthentication(w http.ResponseWriter, r *http.Request, data AuthenticationData) bool {
	if r.Method == http.MethodPost {
		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if s.authenticateUser(username, password) {
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
func (s *AuthenticationService) authenticateUser(username string, password string) bool {
	user, err := s.userStore.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return false
	}
	if user == nil || user.Password != password {
		return false
	}
	return true
}

// renderLoginPage renders the login page with the provided authentication data.
func renderLoginPage(w http.ResponseWriter, data AuthenticationData) {
	tmpl, err := template.ParseFiles(TEMPLATES_DIR + "/login.html")
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
