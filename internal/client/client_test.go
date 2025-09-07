package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("localhost:8080")
	if client == nil {
		t.Fatal("expected client but got nil")
	}
	if client.serverAddr != "localhost:8080" {
		t.Errorf("expected serverAddr %s, got %s", "localhost:8080", client.serverAddr)
	}
	if client.httpClient == nil {
		t.Error("expected httpClient but got nil")
	}
	if client.token != "" {
		t.Errorf("expected empty token, got %s", client.token)
	}
}

func TestClient_IsAuthenticated(t *testing.T) {
	client := NewClient("localhost:8080")

	if client.IsAuthenticated() {
		t.Error("expected not authenticated initially")
	}

	client.SetToken("test-token")
	if !client.IsAuthenticated() {
		t.Error("expected authenticated after setting token")
	}
}

func TestClient_Logout(t *testing.T) {
	client := NewClient("localhost:8080")
	client.SetToken("test-token")

	if !client.IsAuthenticated() {
		t.Error("expected authenticated before logout")
	}

	client.Logout()
	if client.IsAuthenticated() {
		t.Error("expected not authenticated after logout")
	}
}

func TestClient_Register(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/register" {
			t.Errorf("expected path /auth/register, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			t.Errorf("failed to decode user: %v", err)
		}

		if user.Login != "testuser" {
			t.Errorf("expected login testuser, got %s", user.Login)
		}
		if user.Password != "password123" {
			t.Errorf("expected password password123, got %s", user.Password)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "test-jwt-token",
		})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix

	err := client.Register("testuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Error("expected authenticated after registration")
	}
}

func TestClient_Login(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" {
			t.Errorf("expected path /auth/login, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			t.Errorf("failed to decode user: %v", err)
		}

		if user.Login != "testuser" {
			t.Errorf("expected login testuser, got %s", user.Login)
		}
		if user.Password != "password123" {
			t.Errorf("expected password password123, got %s", user.Password)
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "test-jwt-token",
		})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix

	err := client.Login("testuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !client.IsAuthenticated() {
		t.Error("expected authenticated after login")
	}
}

func TestClient_GetDataItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("expected path /data, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("expected Authorization header Bearer test-token, got %s", authHeader)
		}

		response := DataItemsResponse{
			Items: []*DataItem{
				{
					ID:        "item1",
					UserID:    "user1",
					Type:      "login_password",
					Data:      `{"login":"test","password":"pass"}`,
					Metadata:  "test metadata",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	items, err := client.GetDataItems()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	if items[0].ID != "item1" {
		t.Errorf("expected item ID item1, got %s", items[0].ID)
	}
}

func TestClient_GetDataItems_NotAuthenticated(t *testing.T) {
	client := NewClient("localhost:8080")

	_, err := client.GetDataItems()
	if err == nil {
		t.Error("expected error for unauthenticated request")
	}
}

func TestClient_CreateLoginPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("expected path /data, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("expected Authorization header Bearer test-token, got %s", authHeader)
		}

		var req CreateDataItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Type != "login_password" {
			t.Errorf("expected type login_password, got %s", req.Type)
		}

		response := DataItemResponse{
			Item: &DataItem{
				ID:        "new-item",
				UserID:    "user1",
				Type:      "login_password",
				Data:      req.Data,
				Metadata:  req.Metadata,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	item, err := client.CreateLoginPassword("testuser", "password123", "test metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatal("expected item but got nil")
	}

	if item.ID != "new-item" {
		t.Errorf("expected item ID new-item, got %s", item.ID)
	}
}

func TestClient_CreateBankCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("expected path /data, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var req CreateDataItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Type != "bank_card" {
			t.Errorf("expected type bank_card, got %s", req.Type)
		}

		response := DataItemResponse{
			Item: &DataItem{
				ID:        "new-card",
				UserID:    "user1",
				Type:      "bank_card",
				Data:      req.Data,
				Metadata:  req.Metadata,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	item, err := client.CreateBankCard("4111111111111111", "12/25", "123", "John Doe", "card metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatal("expected item but got nil")
	}

	if item.ID != "new-card" {
		t.Errorf("expected item ID new-card, got %s", item.ID)
	}
}

func TestClient_CreateTextData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("expected path /data, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var req CreateDataItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Type != "text" {
			t.Errorf("expected type text, got %s", req.Type)
		}

		response := DataItemResponse{
			Item: &DataItem{
				ID:        "new-text",
				UserID:    "user1",
				Type:      "text",
				Data:      req.Data,
				Metadata:  req.Metadata,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	item, err := client.CreateTextData("test text content", "text metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatal("expected item but got nil")
	}

	if item.ID != "new-text" {
		t.Errorf("expected item ID new-text, got %s", item.ID)
	}
}

func TestClient_CreateBinaryData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data" {
			t.Errorf("expected path /data, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var req CreateDataItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Type != "binary" {
			t.Errorf("expected type binary, got %s", req.Type)
		}

		response := DataItemResponse{
			Item: &DataItem{
				ID:        "new-binary",
				UserID:    "user1",
				Type:      "binary",
				Data:      req.Data,
				Metadata:  req.Metadata,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	testData := []byte("test binary data")
	item, err := client.CreateBinaryData(testData, "text/plain", "test.txt", "binary metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatal("expected item but got nil")
	}

	if item.ID != "new-binary" {
		t.Errorf("expected item ID new-binary, got %s", item.ID)
	}
}

func TestClient_DownloadBinaryFile(t *testing.T) {
	testData := []byte("test binary file content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/data/test-id/download" {
			t.Errorf("expected path /data/test-id/download, got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("expected Authorization header Bearer test-token, got %s", authHeader)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=\"test.txt\"")
		w.Write(testData)
	}))
	defer server.Close()

	client := NewClient(server.URL[7:]) // Remove "http://" prefix
	client.SetToken("test-token")

	fileData, binaryData, err := client.DownloadBinaryFile("test-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(fileData) != string(testData) {
		t.Errorf("expected file data %s, got %s", string(testData), string(fileData))
	}

	if binaryData.FileName != "test.txt" {
		t.Errorf("expected file name test.txt, got %s", binaryData.FileName)
	}

	if binaryData.ContentType != "text/plain" {
		t.Errorf("expected content type text/plain, got %s", binaryData.ContentType)
	}
}
