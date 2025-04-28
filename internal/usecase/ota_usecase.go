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

func (uc *OTAUseCase) GetOTAByID(ctx context.Context, id string) (entity.OTA, error) {
	if id == "" {
		return entity.OTA{}, fmt.Errorf("ID is required")
	}
	return uc.otaRepo.GetByID(ctx, id)
}

func (uc *OTAUseCase) GetOTAsByAppID(ctx context.Context, appID string, cursor string, limit int) ([]entity.OTA, string, error) {
	if appID == "" {
		return nil, "", fmt.Errorf("app ID is required")
	}
	
	if limit <= 0 {
		limit = 10
	}
	
	return uc.otaRepo.GetByAppID(ctx, appID, cursor, limit)
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
