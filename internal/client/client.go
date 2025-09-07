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

type Client struct {
	serverAddr string
	httpClient *http.Client
	token      string
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DataItem struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DataItemsResponse struct {
	Items []*DataItem `json:"items"`
}

type DataItemResponse struct {
	Item *DataItem `json:"item"`
}

type CreateDataItemRequest struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

type LoginPasswordData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BankCardData struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
	Cardholder string `json:"cardholder"`
}

type TextData struct {
	Text string `json:"text"`
}

type BinaryData struct {
	FilePath    string `json:"file_path"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
}

type BinaryDataRequest struct {
	Data        []byte `json:"data"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
}

func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr: serverAddr,
		httpClient: &http.Client{},
		token:      "",
	}
}

func (c *Client) Register(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	return c.makeAuthRequest("/auth/register", user)
}

func (c *Client) Login(login, password string) error {
	user := User{
		Login:    login,
		Password: password,
	}

	return c.makeAuthRequest("/auth/login", user)
}

func (c *Client) Logout() {
	c.token = ""
}

func (c *Client) IsAuthenticated() bool {
	return c.token != ""
}

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
