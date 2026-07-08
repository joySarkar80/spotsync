package reservation

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/reservation/dto"
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

func (h *handler) Reserve(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var req dto.CreateReservationRequest
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

	resp, err := h.service.Reserve(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrZoneNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Parking zone not found",
			})
		case errors.Is(err, ErrZoneFull):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Parking zone is full",
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create reservation",
				Details: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusCreated, httpresponse.Ok("Reservation confirmed successfully", resp))
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	resp, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Ok("My reservations retrieved successfully", resp))
}

func (h *handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	role, _ := c.Get("user_role").(string)

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid reservation id",
		})
	}

	if err := h.service.CancelReservation(userID, role, uint(id)); err != nil {
		switch {
		case errors.Is(err, ErrReservationNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code:    http.StatusNotFound,
				Message: "Reservation not found",
			})
		case errors.Is(err, ErrForbidden):
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Code:    http.StatusForbidden,
				Message: "You are not allowed to cancel this reservation",
			})
		default:
			return c.JSON(http.StatusInternalServerError, httpresponse.Error{
				Code:    http.StatusInternalServerError,
				Message: "Failed to cancel reservation",
				Details: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, httpresponse.Ok("Reservation cancelled successfully", nil))
}

func (h *handler) GetAllReservations(c *echo.Context) error {
	resp, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Ok("All reservations retrieved successfully", resp))
}
