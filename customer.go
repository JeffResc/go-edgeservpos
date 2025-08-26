package edgeservpos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Customer represents each customer in the response array
type Customer struct {
	ServerID      int       `json:"serverId"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	EmailAddress  string    `json:"emailAddress"`
	Point         int       `json:"point"`
	PhoneNumbers  []string  `json:"phoneNumbers"`
	LastVisitDate int64     `json:"lastVisitDate"`
	Addresses     []Address `json:"addresses"`
}

// Address represents a customer's address
type Address struct {
	Address  string `json:"address"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zipCode"`
}

// ListCustomers retrieves all customers for the restaurant with automatic token management
func (c *Client) ListCustomers() ([]Customer, error) {
	token, err := c.GetOAuthToken()
	if err != nil {
		return nil, fmt.Errorf("error getting OAuth token: %w", err)
	}

	return c.listCustomers(token)
}

// listCustomers is the internal function to retrieve customer data
func (c *Client) listCustomers(token string) ([]Customer, error) {
	customerListURL := fmt.Sprintf("%s/%s/backofhouse/customer/list", c.Host, c.RestaurantCode)

	requestBody := map[string]interface{}{
		"serverId":        nil,
		"searchValue":     "",
		"addressRequired": false,
		"zipRequired":     false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequest("POST", customerListURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating customer request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching customer list: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading customer response: %w", err)
	}

	var customers []Customer
	if err := json.Unmarshal(body, &customers); err != nil {
		return nil, fmt.Errorf("error parsing customer JSON: %w", err)
	}

	return customers, nil
}
