package zone

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/zone/dto"
	"spotsync/internal/httpresponse"
	"strconv"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{service: s}
}

func (h *handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	resp, err := h.service.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Ok("Parking zone created successfully", resp))
}

func (h *handler) GetAllZones(c *echo.Context) error {
	resp, err := h.service.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve parking zones",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Ok("Parking zones retrieved successfully", resp))
}

func (h *handler) GetZoneByID(c *echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid zone id",
		})
	}

	resp, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Parking zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve parking zone",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Ok("Parking zone retrieved successfully", resp))
}
