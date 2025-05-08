package idp_service

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sajitha-tj/go-sts/internal/repository/user_repository"
)

const TEMPLATES_DIR = "/home/sajithaj/my-sts-project/go-sts/internal/service/idp_service/templates"

type IdPService struct {
	userStore *user_repository.UserStore
}

func NewIdPService(userStore *user_repository.UserStore) *IdPService {
	return &IdPService{
		userStore: userStore,
	}
}

func (s *IdPService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	flowId, err := getFlowId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := LoginFormData{
		FlowId: flowId,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles(TEMPLATES_DIR + "/login.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *IdPService) HandleLoginCallback(w http.ResponseWriter, r *http.Request) {
	flowId, err := getFlowId(r)
	if err != nil {
		log.Printf("Error getting flowId: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userClaims, err := s.authenticateUser(username, password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	payload := AcceptLoginRequestPayload{
		FlowId:     flowId,
		Success:    true,
		UserClaims: userClaims,
	}

	// Send the payload to the login accepted endpoint
	resp, err := SendLoginAcceptedRequest("http://localhost:8080/auth/login/accept", payload)
	if err != nil {
		log.Printf("Error sending login accepted request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// err = json.NewEncoder(w).Encode(resp)
	// if err != nil {
	// 	log.Printf("Error encoding JSON response: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// Redirect to the /authorize endpoint with flowId and other parameters such as client_id, redirect_uri, etc (from resp).
	authorizeURL := fmt.Sprintf("http://53d45148-f68f-4c1e-8aa8-0a2108a06daa.localhost:8080/authorize?response_type=code&flowId=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s", flowId, resp.ClientID, resp.RedirectURI.String(), strings.Join(resp.Scope, "+"),resp.State)
	http.Redirect(w, r, authorizeURL, http.StatusFound)
}

func (s *IdPService) authenticateUser(username, password string) (*UserClaims, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("Invalid credentials")
	}
	user, err := s.userStore.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("Error fetching user: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("User not found")
	}
	if user.Password != password {
		return nil, fmt.Errorf("Invalid password")
	}
	// User authenticated successfully
	return convertToUserClaims(user), nil
}

func convertToUserClaims(user *user_repository.User) *UserClaims {
	return &UserClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func getFlowId(r *http.Request) (string, error) {
	flowId := r.URL.Query().Get("flowId")
	if flowId == "" {
		flowId = r.FormValue("flowId")
	}
	if flowId == "" {
		return "", fmt.Errorf("flowId is required")
	}
	return flowId, nil
}
