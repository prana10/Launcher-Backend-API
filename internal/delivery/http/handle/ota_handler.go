package handle

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"launcherbackend_api/internal/common/response"
	"launcherbackend_api/internal/domain/entity"
	"launcherbackend_api/internal/usecase"
)

type OTACreateRequest struct {
	AppID        string `json:"app_id" validate:"required" example:"com.yapindo.launcher"`
	VersionName  string `json:"version_name" validate:"required" example:"1.0.0"`
	VersionCode  int    `json:"version_code" validate:"required,gt=0" example:"100"`
	ReleaseNotes string `json:"release_notes" example:"Initial release with basic features"`
	URL          string `json:"url" validate:"required,url" example:"https://storage.example.com/apps/launcher-1.0.0.apk"`
}

type OTAUpdateRequest struct {
	AppID        string `json:"app_id" validate:"required" example:"com.yapindo.launcher"`
	VersionName  string `json:"version_name" validate:"required" example:"1.0.1"`
	VersionCode  int    `json:"version_code" validate:"required,gt=0" example:"101"`
	ReleaseNotes string `json:"release_notes" example:"Bug fixes and performance improvements"`
	URL          string `json:"url" validate:"required,url" example:"https://storage.example.com/apps/launcher-1.0.1.apk"`
}

type OTAHandler struct {
	otaUseCase *usecase.OTAUseCase
}

func NewOTAHandler(otaUseCase *usecase.OTAUseCase) *OTAHandler {
	return &OTAHandler{
		otaUseCase: otaUseCase,
	}
}

func (h *OTAHandler) RegisterRoutes(router fiber.Router) {
	otaRouter := router.Group("/otas")

	otaRouter.Post("/", h.CreateOTA)
	otaRouter.Get("/", h.GetAllOTAs)
	otaRouter.Get("/get", h.GetOTA)
	otaRouter.Put("/:id", h.UpdateOTA)
	otaRouter.Delete("/:id", h.DeleteOTA)
}

func (h *OTAHandler) CreateOTA(c *fiber.Ctx) error {
	var req OTACreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequestResponse(c, "Invalid request body")
	}

	ota := entity.OTA{
		AppID:        req.AppID,
		VersionName:  req.VersionName,
		VersionCode:  req.VersionCode,
		ReleaseNotes: req.ReleaseNotes,
		URL:          req.URL,
	}

	createdOTA, err := h.otaUseCase.CreateOTA(c.Context(), ota)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create OTA: "+err.Error())
	}

	return response.CreatedResponse(c, "OTA created successfully", createdOTA)
}

func (h *OTAHandler) GetOTA(c *fiber.Ctx) error {
	id := c.Query("id", "")
	appID := c.Query("app_id", "")

	ota, err := h.otaUseCase.GetOTA(c.Context(), id, appID)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, "OTA retrieved successfully", ota)
}

func (h *OTAHandler) GetAllOTAs(c *fiber.Ctx) error {
	cursor := c.Query("cursor", "")
	limitStr := c.Query("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	otas, nextCursor, total, err := h.otaUseCase.GetAllOTAs(c.Context(), cursor, limit)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get OTAs: "+err.Error())
	}

	hasNext := nextCursor != ""
	hasPrev := cursor != ""

	return response.PaginatedResponse(c, "OTAs retrieved successfully", otas, hasNext, hasPrev, nextCursor, cursor, total, len(otas))
}

func (h *OTAHandler) UpdateOTA(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequestResponse(c, "ID is required")
	}

	var req OTAUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequestResponse(c, "Invalid request body")
	}

	otas, _, err := h.otaUseCase.GetOTA(c.Context(), id, "", "", 0)
	if err != nil || len(otas) == 0 {
		return response.NotFoundResponse(c, "OTA not found")
	}

	ota := entity.OTA{
		ID:           id,
		AppID:        req.AppID,
		VersionName:  req.VersionName,
		VersionCode:  req.VersionCode,
		ReleaseNotes: req.ReleaseNotes,
		URL:          req.URL,
	}

	updatedOTA, err := h.otaUseCase.UpdateOTA(c.Context(), ota)
	if err != nil {
		return response.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update OTA: "+err.Error())
	}

	return response.SuccessResponse(c, "OTA updated successfully", updatedOTA)
}

func (h *OTAHandler) DeleteOTA(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequestResponse(c, "ID is required")
	}

	err := h.otaUseCase.DeleteOTA(c.Context(), id)
	if err != nil {
		return response.NotFoundResponse(c, "Failed to delete OTA: "+err.Error())
	}

	return response.SuccessResponse(c, "OTA deleted successfully", nil)
}
