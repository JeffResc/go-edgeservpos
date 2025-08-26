package edgeservpos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://api.example.com", "restaurant123", "clientid", "clientsecret", "user", "pass")

	if client.Host != "https://api.example.com" {
		t.Errorf("Expected Host to be 'https://api.example.com', got %s", client.Host)
	}
	if client.RestaurantCode != "restaurant123" {
		t.Errorf("Expected RestaurantCode to be 'restaurant123', got %s", client.RestaurantCode)
	}
	if client.ClientID != "clientid" {
		t.Errorf("Expected ClientID to be 'clientid', got %s", client.ClientID)
	}
	if client.ClientSecret != "clientsecret" {
		t.Errorf("Expected ClientSecret to be 'clientsecret', got %s", client.ClientSecret)
	}
	if client.Username != "user" {
		t.Errorf("Expected Username to be 'user', got %s", client.Username)
	}
	if client.Password != "pass" {
		t.Errorf("Expected Password to be 'pass', got %s", client.Password)
	}
	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized, got nil")
	}
}

func TestGetOAuthToken_Success(t *testing.T) {
	mockResponse := OAuthResponse{Value: "test-token-123"}
	responseBody, _ := json.Marshal(mockResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	token, err := client.GetOAuthToken()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if token != "test-token-123" {
		t.Errorf("Expected token to be 'test-token-123', got %s", token)
	}
}

func TestGetOAuthToken_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	_, err := client.GetOAuthToken()

	if err == nil {
		t.Error("Expected an error for HTTP 500 response, got nil")
	}
}

func TestGetOAuthToken_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	_, err := client.GetOAuthToken()

	if err == nil {
		t.Error("Expected an error for invalid JSON, got nil")
	}
}