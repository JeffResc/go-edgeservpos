package edgeservpos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListCustomers_Success(t *testing.T) {
	mockCustomers := []Customer{
		{
			ServerID:      1,
			FirstName:     "John",
			LastName:      "Doe",
			EmailAddress:  "john.doe@example.com",
			Point:         100,
			PhoneNumbers:  []string{"555-1234"},
			LastVisitDate: 1640995200,
			Addresses: []Address{
				{
					Address:  "123 Main St",
					Address2: "Apt 4B",
					City:     "Anytown",
					State:    "CA",
					ZipCode:  "12345",
				},
			},
		},
	}

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if strings.Contains(r.URL.Path, "oauth/token") {
			tokenResponse := OAuthResponse{Value: "test-token"}
			json.NewEncoder(w).Encode(tokenResponse)
		} else if strings.Contains(r.URL.Path, "customer/list") {
			json.NewEncoder(w).Encode(mockCustomers)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	customers, err := client.ListCustomers()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(customers) != 1 {
		t.Errorf("Expected 1 customer, got %d", len(customers))
	}

	customer := customers[0]
	if customer.FirstName != "John" {
		t.Errorf("Expected FirstName to be 'John', got %s", customer.FirstName)
	}
	if customer.LastName != "Doe" {
		t.Errorf("Expected LastName to be 'Doe', got %s", customer.LastName)
	}
	if customer.EmailAddress != "john.doe@example.com" {
		t.Errorf("Expected EmailAddress to be 'john.doe@example.com', got %s", customer.EmailAddress)
	}
	if customer.Point != 100 {
		t.Errorf("Expected Point to be 100, got %d", customer.Point)
	}
	if len(customer.PhoneNumbers) != 1 || customer.PhoneNumbers[0] != "555-1234" {
		t.Errorf("Expected PhoneNumbers to be ['555-1234'], got %v", customer.PhoneNumbers)
	}
	if customer.LastVisitDate != 1640995200 {
		t.Errorf("Expected LastVisitDate to be 1640995200, got %d", customer.LastVisitDate)
	}
	if len(customer.Addresses) != 1 {
		t.Errorf("Expected 1 address, got %d", len(customer.Addresses))
	}

	address := customer.Addresses[0]
	if address.Address != "123 Main St" {
		t.Errorf("Expected Address to be '123 Main St', got %s", address.Address)
	}
	if address.City != "Anytown" {
		t.Errorf("Expected City to be 'Anytown', got %s", address.City)
	}
	if address.State != "CA" {
		t.Errorf("Expected State to be 'CA', got %s", address.State)
	}
	if address.ZipCode != "12345" {
		t.Errorf("Expected ZipCode to be '12345', got %s", address.ZipCode)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 HTTP calls (token + customers), got %d", callCount)
	}
}

func TestListCustomers_TokenError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "oauth/token") {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	_, err := client.ListCustomers()

	if err == nil {
		t.Error("Expected an error when token request fails, got nil")
	}
}

func TestListCustomers_CustomerListError(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		if strings.Contains(r.URL.Path, "oauth/token") {
			w.WriteHeader(http.StatusOK)
			tokenResponse := OAuthResponse{Value: "test-token"}
			json.NewEncoder(w).Encode(tokenResponse)
		} else if strings.Contains(r.URL.Path, "customer/list") {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	_, err := client.ListCustomers()

	if err == nil {
		t.Error("Expected an error when customer list request fails, got nil")
	}
}

func TestListCustomers_InvalidCustomerJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if strings.Contains(r.URL.Path, "oauth/token") {
			tokenResponse := OAuthResponse{Value: "test-token"}
			json.NewEncoder(w).Encode(tokenResponse)
		} else if strings.Contains(r.URL.Path, "customer/list") {
			w.Write([]byte("invalid json"))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "restaurant123", "clientid", "clientsecret", "user", "pass")
	_, err := client.ListCustomers()

	if err == nil {
		t.Error("Expected an error for invalid customer JSON, got nil")
	}
}