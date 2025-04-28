package response

import (
	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	Total      int64  `json:"total,omitempty"`
	Count      int    `json:"count"`
}

type Meta struct {
	Status     string     `json:"status"`
	Code       int        `json:"code"`
	Message    string     `json:"message"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Meta: Meta{
			Status:  "success",
			Code:    fiber.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Meta: Meta{
			Status:  "success",
			Code:    fiber.StatusCreated,
			Message: message,
		},
		Data: data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Meta: Meta{
			Status:  "error",
			Code:    statusCode,
			Message: message,
		},
	})
}

func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}

func ValidationErrorResponse(c *fiber.Ctx, message string) error {
	return BadRequestResponse(c, message)
}

func PaginatedResponse(c *fiber.Ctx, message string, data interface{}, hasNext bool, hasPrev bool, nextCursor string, prevCursor string, total int64, count int) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Meta: Meta{
			Status:  "success",
			Code:    fiber.StatusOK,
			Message: message,
			Pagination: &Pagination{
				HasNext:    hasNext,
				HasPrev:    hasPrev,
				NextCursor: nextCursor,
				PrevCursor: prevCursor,
				Total:      total,
				Count:      count,
			},
		},
		Data: data,
	})
} 