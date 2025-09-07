package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aseptimu/GophKeeper/internal/app/middleware"
	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/aseptimu/GophKeeper/internal/app/services"
	"github.com/aseptimu/GophKeeper/internal/app/store"
	"github.com/go-chi/chi/v5"
)

type DataHandler struct {
	dataService *services.DataService
}

func NewDataHandler(dataService *services.DataService) *DataHandler {
	return &DataHandler{
		dataService: dataService,
	}
}

func (h *DataHandler) RegisterRoutes(r chi.Router) {
	r.Route("/data", func(r chi.Router) {
		r.Get("/", h.GetDataItems)
		r.Post("/", h.CreateDataItem)
		r.Get("/{id}", h.GetDataItem)
		r.Get("/{id}/download", h.DownloadBinaryFile)
		r.Put("/{id}", h.UpdateDataItem)
		r.Delete("/{id}", h.DeleteDataItem)
	})
}

type CreateDataItemRequest struct {
	Type     string `json:"type"`
	Data     string `json:"data"`
	Metadata string `json:"metadata"`
}

type CreateLoginPasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

type CreateBankCardRequest struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
	Cardholder string `json:"cardholder"`
	Metadata   string `json:"metadata"`
}

type CreateTextDataRequest struct {
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}

type CreateBinaryDataRequest struct {
	Data        []byte `json:"data"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
}

type DataItemsResponse struct {
	Items []*models.DataItem `json:"items"`
}

type DataItemResponse struct {
	Item *models.DataItem `json:"item"`
}

func (h *DataHandler) GetDataItems(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	items, err := h.dataService.GetUserDataItems(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to get data items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DataItemsResponse{Items: items}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) GetDataItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	item, err := h.dataService.GetDataItem(r.Context(), id, userID)
	if err != nil {
		if err == store.ErrDataItemNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("Failed to get data item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DataItemResponse{Item: item}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) CreateDataItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req CreateDataItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item *models.DataItem

	switch models.DataType(req.Type) {
	case models.DataTypeLoginPassword:
		var loginReq CreateLoginPasswordRequest
		if err := json.Unmarshal([]byte(req.Data), &loginReq); err != nil {
			http.Error(w, "Invalid login/password data format", http.StatusBadRequest)
			return
		}
		item, err = h.dataService.SaveLoginPassword(r.Context(), userID, loginReq.Login, loginReq.Password, req.Metadata)

	case models.DataTypeBankCard:
		var cardReq CreateBankCardRequest
		if err := json.Unmarshal([]byte(req.Data), &cardReq); err != nil {
			http.Error(w, "Invalid bank card data format", http.StatusBadRequest)
			return
		}
		item, err = h.dataService.SaveBankCard(r.Context(), userID, cardReq.CardNumber, cardReq.ExpiryDate, cardReq.CVV, cardReq.Cardholder, req.Metadata)

	case models.DataTypeText:
		var textReq CreateTextDataRequest
		if err := json.Unmarshal([]byte(req.Data), &textReq); err != nil {
			http.Error(w, "Invalid text data format", http.StatusBadRequest)
			return
		}
		item, err = h.dataService.SaveTextData(r.Context(), userID, textReq.Text, req.Metadata)

	case models.DataTypeBinary:
		var binaryReq CreateBinaryDataRequest
		if err := json.Unmarshal([]byte(req.Data), &binaryReq); err != nil {
			http.Error(w, "Invalid binary data format", http.StatusBadRequest)
			return
		}
		item, err = h.dataService.SaveBinaryData(r.Context(), userID, binaryReq.Data, binaryReq.ContentType, binaryReq.FileName, req.Metadata)

	default:
		http.Error(w, "Invalid data type", http.StatusBadRequest)
		return
	}

	if err != nil {
		slog.Error("Failed to create data item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := DataItemResponse{Item: item}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *DataHandler) UpdateDataItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	var req CreateDataItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.dataService.UpdateDataItem(r.Context(), id, userID, models.DataType(req.Type), req.Data, req.Metadata)
	if err != nil {
		if err == store.ErrDataItemNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("Failed to update data item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *DataHandler) DeleteDataItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	err := h.dataService.DeleteDataItem(r.Context(), id, userID)
	if err != nil {
		if err == store.ErrDataItemNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("Failed to delete data item", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DataHandler) DownloadBinaryFile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	id := chi.URLParam(r, "id")

	fileData, binaryData, err := h.dataService.GetBinaryFile(r.Context(), id, userID)
	if err != nil {
		if err == store.ErrDataItemNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		slog.Error("Failed to get binary file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", binaryData.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+binaryData.FileName+"\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", binaryData.FileSize))
	w.Write(fileData)
}
