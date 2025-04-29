package usecase

import (
	"context"
	"fmt"

	"launcherbackend_api/internal/domain/entity"
	"launcherbackend_api/internal/domain/repository"
)

type OTAUseCase struct {
	otaRepo repository.OTARepository
}

func NewOTAUseCase(otaRepo repository.OTARepository) *OTAUseCase {
	return &OTAUseCase{
		otaRepo: otaRepo,
	}
}

func (uc *OTAUseCase) CreateOTA(ctx context.Context, ota entity.OTA) (entity.OTA, error) {
	if ota.AppID == "" {
		return entity.OTA{}, fmt.Errorf("app ID is required")
	}
	if ota.VersionName == "" {
		return entity.OTA{}, fmt.Errorf("version name is required")
	}
	if ota.VersionCode <= 0 {
		return entity.OTA{}, fmt.Errorf("valid version code is required")
	}
	if ota.URL == "" {
		return entity.OTA{}, fmt.Errorf("URL is required")
	}

	return uc.otaRepo.Create(ctx, ota)
}

func (uc *OTAUseCase) GetOTA(ctx context.Context, id string, appID string) (entity.OTA, string, error) {
	// Check if both id and appID are provided
	if id != "" && appID != "" {
		return nil, "", fmt.Errorf("cannot provide both id and appID, choose one")
	}
	
	// Check if at least one of id or appID is provided
	if id == "" && appID == "" {
		return nil, "", fmt.Errorf("must provide either id or appID")
	}

	return uc.otaRepo.Get(ctx, id, appID)
}

func (uc *OTAUseCase) GetAllOTAs(ctx context.Context, cursor string, limit int) ([]entity.OTA, string, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	
	return uc.otaRepo.GetAll(ctx, cursor, limit)
}

func (uc *OTAUseCase) UpdateOTA(ctx context.Context, ota entity.OTA) (entity.OTA, error) {
	if ota.ID == "" {
		return entity.OTA{}, fmt.Errorf("ID is required")
	}
	if ota.AppID == "" {
		return entity.OTA{}, fmt.Errorf("app ID is required")
	}
	if ota.VersionName == "" {
		return entity.OTA{}, fmt.Errorf("version name is required")
	}
	if ota.VersionCode <= 0 {
		return entity.OTA{}, fmt.Errorf("valid version code is required")
	}
	if ota.URL == "" {
		return entity.OTA{}, fmt.Errorf("URL is required")
	}
	
	return uc.otaRepo.Update(ctx, ota)
}

func (uc *OTAUseCase) DeleteOTA(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}
	return uc.otaRepo.Delete(ctx, id)
}
