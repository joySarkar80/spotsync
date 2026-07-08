package order

import (
	"errors"
	"haddibanga/internal/domain/mango"
	"haddibanga/internal/domain/order/dto"
	"haddibanga/internal/httpresponse"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(s *service) *handler {
	return &handler{service: s}
}

func getCurrentUserID(c *echo.Context) (uint, bool) {
	userId, ok := c.Get("user_id").(uint)
	return userId, ok
}

func orderErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrOrderNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Order not found",
		})
	}
	if errors.Is(err, mango.ErrMangoNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Mango not found",
		})
	}
	if errors.Is(err, ErrNotEnoughStock) {
		return c.JSON(http.StatusConflict, httpresponse.Error{
			Code:    http.StatusConflict,
			Message: "Not enough stock available",
		})
	}
	if errors.Is(err, ErrOrderAlreadyCancelled) {
		return c.JSON(http.StatusConflict, httpresponse.Error{
			Code:    http.StatusConflict,
			Message: "Order is already cancelled",
		})
	}
	if errors.Is(err, ErrForbiddenOrderAccess) {
		return c.JSON(http.StatusForbidden, httpresponse.Error{
			Code:    http.StatusForbidden,
			Message: "You do not own this order",
		})
	}
	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateOrder(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var req dto.CreateRequest
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

	response, err := h.service.CreateOrder(userId, req)
	if err != nil {
		return orderErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetMyOrders(c *echo.Context) error {
	userId, ok := getCurrentUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	orders, err := h.service.GetMyOrders(userId)
	if err != nil {
		return orderErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, orders)
}
