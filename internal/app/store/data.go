package store

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/aseptimu/GophKeeper/internal/app/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDataItemNotFound = errors.New("data item not found")

type DataStore interface {
	SaveDataItem(ctx context.Context, userID string, dataType models.DataType, data, metadata string) (*models.DataItem, error)
	GetDataItemsByUserID(ctx context.Context, userID string) ([]*models.DataItem, error)
	GetDataItemByID(ctx context.Context, id, userID string) (*models.DataItem, error)
	UpdateDataItem(ctx context.Context, id, userID, data, metadata string) error
	DeleteDataItem(ctx context.Context, id, userID string) error
}

const saveDataItemQuery = `
	INSERT INTO data_items (user_id, type, data, metadata)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
`

const getDataItemsByUserIDQuery = `
	SELECT id, user_id, type, data, metadata, created_at, updated_at
	FROM data_items
	WHERE user_id = $1
	ORDER BY created_at DESC
`

const getDataItemByIDQuery = `
	SELECT id, user_id, type, data, metadata, created_at, updated_at
	FROM data_items
	WHERE id = $1 AND user_id = $2
`

const updateDataItemQuery = `
	UPDATE data_items
	SET data = $3, metadata = $4, updated_at = now()
	WHERE id = $1 AND user_id = $2
	RETURNING updated_at
`

const deleteDataItemQuery = `
	DELETE FROM data_items
	WHERE id = $1 AND user_id = $2
`

type DataStoreImpl struct {
	pool *pgxpool.Pool
}

func NewDataStore(pool *pgxpool.Pool) *DataStoreImpl {
	return &DataStoreImpl{pool: pool}
}

func (s *DataStoreImpl) SaveDataItem(ctx context.Context, userID string, dataType models.DataType, data, metadata string) (*models.DataItem, error) {
	var item models.DataItem

	err := s.pool.QueryRow(ctx, saveDataItemQuery, userID, string(dataType), data, metadata).Scan(
		&item.ID, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to save data item", "error", err)
		return nil, err
	}

	item.UserID = userID
	item.Type = dataType
	item.Data = data
	item.Metadata = metadata

	slog.Info("Saved data item", "id", item.ID, "type", dataType, "user_id", userID)
	return &item, nil
}

func (s *DataStoreImpl) GetDataItemsByUserID(ctx context.Context, userID string) ([]*models.DataItem, error) {
	rows, err := s.pool.Query(ctx, getDataItemsByUserIDQuery, userID)
	if err != nil {
		slog.Error("Failed to get data items", "error", err)
		return nil, err
	}
	defer rows.Close()

	var items []*models.DataItem
	for rows.Next() {
		var item models.DataItem
		var dataType string

		err := rows.Scan(
			&item.ID, &item.UserID, &dataType, &item.Data, &item.Metadata,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			slog.Error("Failed to scan data item", "error", err)
			return nil, err
		}

		item.Type = models.DataType(dataType)
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating data items", "error", err)
		return nil, err
	}

	return items, nil
}

func (s *DataStoreImpl) GetDataItemByID(ctx context.Context, id, userID string) (*models.DataItem, error) {
	var item models.DataItem
	var dataType string

	err := s.pool.QueryRow(ctx, getDataItemByIDQuery, id, userID).Scan(
		&item.ID, &item.UserID, &dataType, &item.Data, &item.Metadata,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDataItemNotFound
		}
		slog.Error("Failed to get data item", "error", err)
		return nil, err
	}

	item.Type = models.DataType(dataType)
	return &item, nil
}

func (s *DataStoreImpl) UpdateDataItem(ctx context.Context, id, userID, data, metadata string) error {
	var updatedAt time.Time

	err := s.pool.QueryRow(ctx, updateDataItemQuery, id, userID, data, metadata).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDataItemNotFound
		}
		slog.Error("Failed to update data item", "error", err)
		return err
	}

	slog.Info("Updated data item", "id", id, "user_id", userID)
	return nil
}

func (s *DataStoreImpl) DeleteDataItem(ctx context.Context, id, userID string) error {
	result, err := s.pool.Exec(ctx, deleteDataItemQuery, id, userID)
	if err != nil {
		slog.Error("Failed to delete data item", "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrDataItemNotFound
	}

	slog.Info("Deleted data item", "id", id, "user_id", userID)
	return nil
}
