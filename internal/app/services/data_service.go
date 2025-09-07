package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/aseptimu/GophKeeper/internal/app/store"
	"github.com/aseptimu/GophKeeper/internal/app/utils"
)

type DataService struct {
	store       store.DataStore
	fileStorage *store.FileStorage
}

func NewDataService(store store.DataStore, fileStorage *store.FileStorage) *DataService {
	return &DataService{
		store:       store,
		fileStorage: fileStorage,
	}
}

func (s *DataService) SaveLoginPassword(ctx context.Context, userID string, login, password, metadata string) (*models.DataItem, error) {
	data := models.LoginPasswordData{
		Login:    login,
		Password: password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return s.store.SaveDataItem(ctx, userID, models.DataTypeLoginPassword, string(jsonData), metadata)
}

func (s *DataService) SaveTextData(ctx context.Context, userID, text, metadata string) (*models.DataItem, error) {
	return s.store.SaveDataItem(ctx, userID, models.DataTypeText, text, metadata)
}

func (s *DataService) SaveBinaryData(ctx context.Context, userID string, data []byte, contentType, fileName, metadata string) (*models.DataItem, error) {
	if len(data) == 0 {
		return nil, errors.New("binary data cannot be empty")
	}
	if contentType == "" {
		return nil, errors.New("content type cannot be empty")
	}
	if fileName == "" {
		return nil, errors.New("file name cannot be empty")
	}

	filePath, err := s.fileStorage.SaveFile(data, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileSize, err := s.fileStorage.GetFileSize(filePath)
	if err != nil {
		s.fileStorage.DeleteFile(filePath)
		return nil, fmt.Errorf("failed to get file size: %w", err)
	}

	binaryData := models.BinaryData{
		FilePath:    filePath,
		ContentType: contentType,
		FileName:    fileName,
		FileSize:    fileSize,
	}

	jsonData, err := json.Marshal(binaryData)
	if err != nil {
		s.fileStorage.DeleteFile(filePath)
		return nil, fmt.Errorf("failed to marshal binary data: %w", err)
	}

	return s.store.SaveDataItem(ctx, userID, models.DataTypeBinary, string(jsonData), metadata)
}

func (s *DataService) SaveBankCard(ctx context.Context, userID, cardNumber, expiryDate, cvv, cardholder, metadata string) (*models.DataItem, error) {
	if err := utils.ValidateBankCard(cardNumber, expiryDate, cvv); err != nil {
		slog.Error("Bank card validation failed", "error", err)
		return nil, err
	}

	data := models.BankCardData{
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
		Cardholder: cardholder,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return s.store.SaveDataItem(ctx, userID, models.DataTypeBankCard, string(jsonData), metadata)
}

func (s *DataService) GetUserDataItems(ctx context.Context, userID string) ([]*models.DataItem, error) {
	return s.store.GetDataItemsByUserID(ctx, userID)
}

func (s *DataService) GetDataItem(ctx context.Context, id, userID string) (*models.DataItem, error) {
	return s.store.GetDataItemByID(ctx, id, userID)
}

func (s *DataService) UpdateDataItem(ctx context.Context, id, userID string, dataType models.DataType, data, metadata string) error {
	switch dataType {
	case models.DataTypeLoginPassword:
		var loginData models.LoginPasswordData
		if err := json.Unmarshal([]byte(data), &loginData); err != nil {
			return ErrInvalidDataFormat
		}
	case models.DataTypeBankCard:
		var cardData models.BankCardData
		if err := json.Unmarshal([]byte(data), &cardData); err != nil {
			return ErrInvalidDataFormat
		}
		if err := utils.ValidateBankCard(cardData.CardNumber, cardData.ExpiryDate, cardData.CVV); err != nil {
			return err
		}
	case models.DataTypeText, models.DataTypeBinary:
	default:
		return ErrInvalidDataType
	}

	return s.store.UpdateDataItem(ctx, id, userID, data, metadata)
}

func (s *DataService) DeleteDataItem(ctx context.Context, id, userID string) error {
	item, err := s.store.GetDataItemByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if item.Type == models.DataTypeBinary {
		var binaryData models.BinaryData
		if err := json.Unmarshal([]byte(item.Data), &binaryData); err == nil {
			s.fileStorage.DeleteFile(binaryData.FilePath)
		}
	}

	return s.store.DeleteDataItem(ctx, id, userID)
}

func (s *DataService) GetBinaryFile(ctx context.Context, id, userID string) ([]byte, *models.BinaryData, error) {
	item, err := s.store.GetDataItemByID(ctx, id, userID)
	if err != nil {
		return nil, nil, err
	}

	if item.Type != models.DataTypeBinary {
		return nil, nil, errors.New("item is not binary data")
	}

	var binaryData models.BinaryData
	if err := json.Unmarshal([]byte(item.Data), &binaryData); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal binary data: %w", err)
	}

	fileData, err := s.fileStorage.GetFile(binaryData.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %w", err)
	}

	return fileData, &binaryData, nil
}
