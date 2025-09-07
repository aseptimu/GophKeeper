// Package client provides a client for interacting with the GophKeeper server.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents a client connection to the GophKeeper server.
type Client struct {
	serverAddr string       // Address of the GophKeeper server
	httpClient *http.Client // HTTP client for making requests
	token      string       // JWT token for authentication
}

// User represents user credentials for authentication.
type User struct {
	Login    string `json:"login"`    // User's login/username
	Password string `json:"password"` // User's password
}

// DataItem represents a piece of user data stored on the server.
type DataItem struct {
	ID        string    `json:"id"`         // Unique identifier for the data item
	UserID    string    `json:"user_id"`    // ID of the user who owns this data
	Type      string    `json:"type"`       // Type of data (login_password, text, binary, bank_card)
	Data      string    `json:"data"`       // The actual data content
	Metadata  string    `json:"metadata"`   // Optional metadata describing the data
	CreatedAt time.Time `json:"created_at"` // When the item was created
	UpdatedAt time.Time `json:"updated_at"` // When the item was last updated
}

// DataItemsResponse represents the response containing multiple data items.
type DataItemsResponse struct {
	Items []*DataItem `json:"items"` // List of data items
}

// DataItemResponse represents the response containing a single data item.
type DataItemResponse struct {
	Item *DataItem `json:"item"` // The data item
}

// CreateDataItemRequest represents a request to create a new data item.
type CreateDataItemRequest struct {
	Type     string `json:"type"`     // Type of data to create
	Data     string `json:"data"`     // The data content
	Metadata string `json:"metadata"` // Optional metadata
}

// LoginPasswordData represents login and password information.
type LoginPasswordData struct {
	Login    string `json:"login"`    // Username or login identifier
	Password string `json:"password"` // Password for the login
}

// BankCardData represents bank card information.
type BankCardData struct {
	CardNumber string `json:"card_number"` // Bank card number
	ExpiryDate string `json:"expiry_date"` // Card expiry date in MM/YY format
	CVV        string `json:"cvv"`         // Card verification value
	Cardholder string `json:"cardholder"`  // Name of the cardholder
}

// TextData represents arbitrary text data.
type TextData struct {
	Text string `json:"text"` // The text content
}

// BinaryData represents information about a binary file.
type BinaryData struct {
	FilePath    string `json:"file_path"`    // Path to the stored file
	ContentType string `json:"content_type"` // MIME type of the file
	FileName    string `json:"file_name"`    // Original filename
	FileSize    int64  `json:"file_size"`    // Size of the file in bytes
}

// BinaryDataRequest represents a request to create binary data.
type BinaryDataRequest struct {
	Data        []byte `json:"data"`         // The binary data
	ContentType string `json:"content_type"` // MIME type of the data
	FileName    string `json:"file_name"`    // Original filename
}

// NewClient creates a new client instance for connecting to the GophKeeper server.
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		httpClient: &http.Client{},
		token:      "",
	}
}

// Register registers a new user with the server.
// Returns an error if registration fails.
func (c *Client) Register(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	return c.makeAuthRequest("/auth/register", user)
}

// Login authenticates a user with the server.
// Returns an error if authentication fails.
func (c *Client) Login(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	return c.makeAuthRequest("/auth/login", user)
}

// Logout clears the authentication token.
func (c *Client) Logout() {
	c.token = ""
}

// IsAuthenticated returns true if the client has a valid authentication token.
func (c *Client) IsAuthenticated() bool {
	return c.token != ""
}

// SetToken sets the authentication token for the client.
func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) makeAuthRequest(endpoint string, user User) error {
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %w", err)
	}

	url := fmt.Sprintf("http://%s%s", c.serverAddr, endpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var token string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" || cookie.Value != "" {
			token = cookie.Value
			break
		}
	}

	if token == "" {
		return fmt.Errorf("no token received from server")
	}

	c.token = token
	return nil
}

func (c *Client) makeAuthenticatedRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("http://%s%s", c.serverAddr, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	return c.httpClient.Do(req)
}

func (c *Client) GetDataItems() ([]*DataItem, error) {
	resp, err := c.makeAuthenticatedRequest("GET", "/data", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return []*DataItem{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var dataResp DataItemsResponse
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return dataResp.Items, nil
}

func (c *Client) CreateLoginPassword(login, password, metadata string) (*DataItem, error) {
	data := LoginPasswordData{
		Login:    login,
		Password: password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	req := CreateDataItemRequest{
		Type:     "login_password",
		Data:     string(jsonData),
		Metadata: metadata,
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/data", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var dataResp DataItemResponse
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return dataResp.Item, nil
}

func (c *Client) CreateBankCard(cardNumber, expiryDate, cvv, cardholder, metadata string) (*DataItem, error) {
	data := BankCardData{
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
		Cardholder: cardholder,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	req := CreateDataItemRequest{
		Type:     "bank_card",
		Data:     string(jsonData),
		Metadata: metadata,
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/data", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var dataResp DataItemResponse
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return dataResp.Item, nil
}

func (c *Client) CreateTextData(text, metadata string) (*DataItem, error) {
	data := TextData{
		Text: text,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	req := CreateDataItemRequest{
		Type:     "text",
		Data:     string(jsonData),
		Metadata: metadata,
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/data", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var dataResp DataItemResponse
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return dataResp.Item, nil
}

func (c *Client) CreateBinaryData(data []byte, contentType, fileName, metadata string) (*DataItem, error) {
	binaryData := BinaryDataRequest{
		Data:        data,
		ContentType: contentType,
		FileName:    fileName,
	}

	jsonData, err := json.Marshal(binaryData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	req := CreateDataItemRequest{
		Type:     "binary",
		Data:     string(jsonData),
		Metadata: metadata,
	}

	resp, err := c.makeAuthenticatedRequest("POST", "/data", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var dataResp DataItemResponse
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return dataResp.Item, nil
}

func (c *Client) DownloadBinaryFile(id string) ([]byte, *BinaryData, error) {
	resp, err := c.makeAuthenticatedRequest("GET", fmt.Sprintf("/data/%s/download", id), nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("server error: %s (status: %d)", string(body), resp.StatusCode)
	}

	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file data: %w", err)
	}

	contentType := resp.Header.Get("Content-Type")
	fileName := resp.Header.Get("Content-Disposition")
	if fileName != "" {
		fileName = strings.TrimPrefix(fileName, "attachment; filename=\"")
		fileName = strings.TrimSuffix(fileName, "\"")
	}

	binaryData := &BinaryData{
		ContentType: contentType,
		FileName:    fileName,
		FileSize:    int64(len(fileData)),
	}

	return fileData, binaryData, nil
}
