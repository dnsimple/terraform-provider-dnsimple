package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Remove this once the official client supports the new endpoints

type DNSimpleClient struct {
	AccountID string
	Token     string
	baseURL   string
}

type RegistrantChangeCheck struct {
	ContactID           int64                    `json:"contact_id"`
	DomainID            int64                    `json:"domain_id"`
	ExtendedAttributes  []map[string]interface{} `json:"extended_attributes"`
	RegistryOwnerChange bool                     `json:"registry_owner_change"`
}

type RegistrantChangeCheckResponse struct {
	Data RegistrantChangeCheck `json:"data"`
}

type RegistrantChange struct {
	Id                  int64                    `json:"id"`
	ContactID           int64                    `json:"contact_id"`
	DomainID            int64                    `json:"domain_id"`
	ExtendedAttributes  []map[string]interface{} `json:"extended_attributes"`
	RegistryOwnerChange bool                     `json:"registry_owner_change"`
	State               string                   `json:"state"`
	IRTLockLiftedBy     string                   `json:"irt_lock_lifted_by"`
	CreatedAt           string                   `json:"created_at"`
	UpdatedAt           string                   `json:"updated_at"`
}

type RegistrantChangeResponse struct {
	Data RegistrantChange `json:"data"`
}

type RegistrantChangeInput struct {
	ContactID          int64                  `json:"contact_id"`
	DomainID           string                 `json:"domain_id"`
	ExtendedAttributes map[string]interface{} `json:"extended_attributes"`
}

func NewDNSimpleClient(accountID string, token string, sandbox bool) *DNSimpleClient {
	var baseURL string

	if sandbox {
		baseURL = "https://api.sandbox.dnsimple.com"
	} else {
		baseURL = "https://api.dnsimple.com"
	}

	return &DNSimpleClient{
		AccountID: accountID,
		Token:     token,
		baseURL:   baseURL,
	}
}

// Create a simple HTTP client
func CreateClient() *http.Client {
	return &http.Client{}
}

// Create a simple HTTP request
func CreateRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Create a simple HTTP request with a JSON body
func CreateJSONRequest(method string, url string, body interface{}) (*http.Request, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := CreateRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Create a simple HTTP request with a JSON body and an API token
func (c *DNSimpleClient) CreateJSONRequestWithToken(method string, url string, body interface{}, token string) (*http.Request, error) {
	req, err := CreateJSONRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

// Make a POST request to check_registrant_change endpoint at /%s/registrar/registrant_changes/check
func (c *DNSimpleClient) CheckRegistrantChange(domain string, contactID int64) (*RegistrantChangeCheckResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/registrar/registrant_changes/check", c.baseURL, c.AccountID)
	body := map[string]interface{}{
		"domain":     domain,
		"contact_id": contactID,
	}

	req, err := c.CreateJSONRequestWithToken("POST", url, body, c.Token)
	if err != nil {
		return nil, err
	}

	res, err := CreateClient().Do(req)
	if err != nil {
		return nil, err
	}

	// Serialize the response body to a RegistrantChangeCheck struct
	var registrantChangeCheck RegistrantChangeCheckResponse
	err = json.NewDecoder(res.Body).Decode(&registrantChangeCheck)
	if err != nil {
		return nil, err
	}

	return &registrantChangeCheck, nil
}

// Make a POST request to initiate a registrant change at /%s/registrar/registrant_changes
func (c *DNSimpleClient) CreateRegistrantChange(input RegistrantChangeInput) (*RegistrantChangeResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/registrar/registrant_changes", c.baseURL, c.AccountID)

	req, err := c.CreateJSONRequestWithToken("POST", url, input, c.Token)
	if err != nil {
		return nil, err
	}

	res, err := CreateClient().Do(req)
	if err != nil {
		return nil, err
	}

	// Serialize the response body to a RegistrantChange struct
	var registrantChange RegistrantChangeResponse
	err = json.NewDecoder(res.Body).Decode(&registrantChange)
	if err != nil {
		return nil, err
	}

	return &registrantChange, nil
}

// Make a GET request to get a registrant change at /%s/registrar/registrant_changes/%s
func (c *DNSimpleClient) GetRegistrantChange(registrantChangeID string) (*RegistrantChangeResponse, error) {
	url := fmt.Sprintf("%s/v2/%s/registrar/registrant_changes/%s", c.baseURL, c.AccountID, registrantChangeID)

	req, err := c.CreateJSONRequestWithToken("GET", url, nil, c.Token)
	if err != nil {
		return nil, err
	}

	res, err := CreateClient().Do(req)
	if err != nil {
		return nil, err
	}

	// Serialize the response body to a RegistrantChange struct
	var registrantChange RegistrantChangeResponse
	err = json.NewDecoder(res.Body).Decode(&registrantChange)
	if err != nil {
		return nil, err
	}

	return &registrantChange, nil
}

// Make a DELETE request to cancel a registrant change at /%s/registrar/registrant_changes/%s
func (c *DNSimpleClient) CancelRegistrantChange(registrantChangeID string) error {
	url := fmt.Sprintf("%s/v2/%s/registrar/registrant_changes/%s", c.baseURL, c.AccountID, registrantChangeID)

	req, err := c.CreateJSONRequestWithToken("DELETE", url, nil, c.Token)
	if err != nil {
		return err
	}

	_, err = CreateClient().Do(req)
	if err != nil {
		return err
	}

	return nil
}
