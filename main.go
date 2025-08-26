package edgeservpos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents an EdgeServPOS API client
type Client struct {
	Host           string
	RestaurantCode string
	ClientID       string
	ClientSecret   string
	Username       string
	Password       string
	HTTPClient     *http.Client
}

// OAuthResponse represents the expected JSON response from the OAuth token endpoint
type OAuthResponse struct {
	Value string `json:"value"`
}

// NewClient creates a new EdgeServPOS API client
func NewClient(host, restaurantCode, clientID, clientSecret, username, password string) *Client {
	return &Client{
		Host:           host,
		RestaurantCode: restaurantCode,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Username:       username,
		Password:       password,
		HTTPClient:     &http.Client{},
	}
}

// GetOAuthToken retrieves an OAuth token for the client
func (c *Client) GetOAuthToken() (string, error) {
	tokenURL := fmt.Sprintf("%s/%s/oauth/token?grant_type=password&client_id=%s&client_secret=%s&username=%s&password=%s",
		c.Host, c.RestaurantCode, c.ClientID, c.ClientSecret, c.Username, c.Password)

	resp, err := c.HTTPClient.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("error fetching token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading token response: %w", err)
	}

	var tokenResp OAuthResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("error parsing token JSON: %w", err)
	}

	return tokenResp.Value, nil
}
