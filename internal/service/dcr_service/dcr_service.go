package dcr_service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"net/http"

	"github.com/sajitha-tj/go-sts/internal/repository/client_repository"
	"golang.org/x/crypto/bcrypt"
)

type DcrService struct {
	clientStore *client_repository.ClientStore
}

func NewDcrService(clientStore *client_repository.ClientStore) *DcrService {
	return &DcrService{
		clientStore: clientStore,
	}
}

func (d *DcrService) RegisterClient(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var request ClientRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, e := d.clientStore.GetClient(ctx, request.ClientID)
	if e == nil {
		http.Error(w, "Client ID already exists", http.StatusConflict)
		return
	}

	secret := rand.Text()
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &client_repository.Client{
		ClientID:      request.ClientID,
		ClientSecret:  string(hashedSecret),
		RedirectURIs:  request.RedirectURIs,
		GrantTypes:    request.GrantTypes,
		ResponseTypes: request.ResponseTypes,
		Scopes:        request.Scopes,
		Public:        request.Public,
		Audience:      request.Audience,
	}

	if err := d.clientStore.CreateClient(ctx, client); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}
