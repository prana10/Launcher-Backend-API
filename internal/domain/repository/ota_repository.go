package repository

import (
	"context"

	"launcherbackend_api/internal/domain/entity"
)

type OTARepository interface {
	Create(ctx context.Context, ota entity.OTA) (entity.OTA, error)
	Get(ctx context.Context, id string, appID string) (entity.OTA, string, error)
	GetAll(ctx context.Context, cursor string, limit int) ([]entity.OTA, string, int64, error)
	Update(ctx context.Context, ota entity.OTA) (entity.OTA, error)
	Delete(ctx context.Context, id string) error
}
