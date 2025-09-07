package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/aseptimu/GophKeeper/internal/app/store"
)

type mockDataStore struct {
	items  map[string]*models.DataItem
	nextID int
}

func newMockDataStore() *mockDataStore {
	return &mockDataStore{
		items:  make(map[string]*models.DataItem),
		nextID: 1,
	}
}

func (m *mockDataStore) SaveDataItem(_ context.Context, userID string, dataType models.DataType, data, metadata string) (*models.DataItem, error) {
	item := &models.DataItem{
		ID:       fmt.Sprintf("test-id-%d", m.nextID),
		UserID:   userID,
		Type:     dataType,
		Data:     data,
		Metadata: metadata,
	}
	m.nextID++
	m.items[item.ID] = item
	return item, nil
}

func (m *mockDataStore) GetDataItemsByUserID(_ context.Context, userID string) ([]*models.DataItem, error) {
	var items []*models.DataItem
	for _, item := range m.items {
		if item.UserID == userID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (m *mockDataStore) GetDataItemByID(_ context.Context, id, userID string) (*models.DataItem, error) {
	item, exists := m.items[id]
	if !exists {
		return nil, store.ErrDataItemNotFound
	}
	if item.UserID != userID {
		return nil, store.ErrDataItemNotFound
	}
	return item, nil
}

func (m *mockDataStore) UpdateDataItem(_ context.Context, id, userID, data, metadata string) error {
	item, exists := m.items[id]
	if !exists {
		return store.ErrDataItemNotFound
	}
	if item.UserID != userID {
		return store.ErrDataItemNotFound
	}
	item.Data = data
	item.Metadata = metadata
	return nil
}

func (m *mockDataStore) DeleteDataItem(_ context.Context, id, userID string) error {
	item, exists := m.items[id]
	if !exists {
		return store.ErrDataItemNotFound
	}
	if item.UserID != userID {
		return store.ErrDataItemNotFound
	}
	delete(m.items, id)
	return nil
}

type mockFileStorage struct {
	files map[string][]byte
}

func newMockFileStorage() *mockFileStorage {
	return &mockFileStorage{
		files: make(map[string][]byte),
	}
}

func (m *mockFileStorage) SaveFile(data []byte, fileName string) (string, error) {
	filePath := "/mock/path/" + fileName
	m.files[filePath] = data
	return filePath, nil
}

func (m *mockFileStorage) GetFile(filePath string) ([]byte, error) {
	data, exists := m.files[filePath]
	if !exists {
		return nil, errors.New("file not found")
	}
	return data, nil
}

func (m *mockFileStorage) DeleteFile(filePath string) error {
	delete(m.files, filePath)
	return nil
}

func (m *mockFileStorage) GetFileSize(filePath string) (int64, error) {
	data, exists := m.files[filePath]
	if !exists {
		return 0, errors.New("file not found")
	}
	return int64(len(data)), nil
}

func TestDataService_SaveLoginPassword(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	item, err := service.SaveLoginPassword(context.Background(), "user1", "testuser", "password123", "test metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatalf("expected item but got nil")
	}

	if item.Type != models.DataTypeLoginPassword {
		t.Errorf("expected type %s, got %s", models.DataTypeLoginPassword, item.Type)
	}

	if item.UserID != "user1" {
		t.Errorf("expected userID %s, got %s", "user1", item.UserID)
	}

	if item.Metadata != "test metadata" {
		t.Errorf("expected metadata %s, got %s", "test metadata", item.Metadata)
	}
}

func TestDataService_SaveTextData(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	item, err := service.SaveTextData(context.Background(), "user1", "test text data", "text metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatalf("expected item but got nil")
	}

	if item.Type != models.DataTypeText {
		t.Errorf("expected type %s, got %s", models.DataTypeText, item.Type)
	}

	if item.Data != "test text data" {
		t.Errorf("expected data %s, got %s", "test text data", item.Data)
	}
}

func TestDataService_SaveBinaryData(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	testData := []byte("test binary data")
	item, err := service.SaveBinaryData(context.Background(), "user1", testData, "text/plain", "test.txt", "binary metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatalf("expected item but got nil")
	}

	if item.Type != models.DataTypeBinary {
		t.Errorf("expected type %s, got %s", models.DataTypeBinary, item.Type)
	}

	if item.Metadata != "binary metadata" {
		t.Errorf("expected metadata %s, got %s", "binary metadata", item.Metadata)
	}
}

func TestDataService_SaveBinaryData_EmptyData(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	_, err := service.SaveBinaryData(context.Background(), "user1", []byte{}, "text/plain", "test.txt", "metadata")
	if err == nil {
		t.Errorf("expected error for empty data but got none")
	}
}

func TestDataService_SaveBinaryData_EmptyContentType(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	_, err := service.SaveBinaryData(context.Background(), "user1", []byte("data"), "", "test.txt", "metadata")
	if err == nil {
		t.Errorf("expected error for empty content type but got none")
	}
}

func TestDataService_SaveBinaryData_EmptyFileName(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	_, err := service.SaveBinaryData(context.Background(), "user1", []byte("data"), "text/plain", "", "metadata")
	if err == nil {
		t.Errorf("expected error for empty file name but got none")
	}
}

func TestDataService_SaveBankCard(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	item, err := service.SaveBankCard(context.Background(), "user1", "4111111111111111", "12/25", "123", "John Doe", "card metadata")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if item == nil {
		t.Fatalf("expected item but got nil")
	}

	if item.Type != models.DataTypeBankCard {
		t.Errorf("expected type %s, got %s", models.DataTypeBankCard, item.Type)
	}

	if item.Metadata != "card metadata" {
		t.Errorf("expected metadata %s, got %s", "card metadata", item.Metadata)
	}
}

func TestDataService_SaveBankCard_InvalidCard(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	_, err := service.SaveBankCard(context.Background(), "user1", "1234567890123456", "12/25", "123", "John Doe", "card metadata")
	if err == nil {
		t.Errorf("expected error for invalid card number but got none")
	}
}

func TestDataService_GetUserDataItems(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	// Save some test data
	_, err := service.SaveLoginPassword(context.Background(), "user1", "testuser", "password", "metadata1")
	if err != nil {
		t.Fatalf("failed to save login password: %v", err)
	}

	_, err = service.SaveTextData(context.Background(), "user1", "test text", "metadata2")
	if err != nil {
		t.Fatalf("failed to save text data: %v", err)
	}

	items, err := service.GetUserDataItems(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestDataService_GetDataItem(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	// Save test data
	item, err := service.SaveLoginPassword(context.Background(), "user1", "testuser", "password", "metadata")
	if err != nil {
		t.Fatalf("failed to save login password: %v", err)
	}

	// Get the item
	retrievedItem, err := service.GetDataItem(context.Background(), item.ID, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if retrievedItem.ID != item.ID {
		t.Errorf("expected ID %s, got %s", item.ID, retrievedItem.ID)
	}
}

func TestDataService_GetDataItem_NotFound(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	_, err := service.GetDataItem(context.Background(), "nonexistent", "user1")
	if err == nil {
		t.Errorf("expected error for nonexistent item but got none")
	}
}

func TestDataService_DeleteDataItem(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	// Save test data
	item, err := service.SaveLoginPassword(context.Background(), "user1", "testuser", "password", "metadata")
	if err != nil {
		t.Fatalf("failed to save login password: %v", err)
	}

	// Delete the item
	err = service.DeleteDataItem(context.Background(), item.ID, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Try to get the deleted item
	_, err = service.GetDataItem(context.Background(), item.ID, "user1")
	if err == nil {
		t.Errorf("expected error for deleted item but got none")
	}
}

func TestDataService_GetBinaryFile(t *testing.T) {
	mockStore := newMockDataStore()
	mockFileStorage := newMockFileStorage()
	service := NewDataService(mockStore, mockFileStorage)

	testData := []byte("test binary data")
	item, err := service.SaveBinaryData(context.Background(), "user1", testData, "text/plain", "test.txt", "metadata")
	if err != nil {
		t.Fatalf("failed to save binary data: %v", err)
	}

	fileData, binaryData, err := service.GetBinaryFile(context.Background(), item.ID, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(fileData) != string(testData) {
		t.Errorf("expected file data %s, got %s", string(testData), string(fileData))
	}

	if binaryData.FileName != "test.txt" {
		t.Errorf("expected file name %s, got %s", "test.txt", binaryData.FileName)
	}
}
