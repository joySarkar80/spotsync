package mango

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/mango/dto"
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

func mangoErrorResponse(c *echo.Context, err error) error {
	if errors.Is(err, ErrMangoNotFound) {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Mango not found",
		})
	}
	return c.JSON(http.StatusInternalServerError, httpresponse.Error{
		Code:    http.StatusInternalServerError,
		Message: "Something went wrong",
		Details: err.Error(),
	})
}

func (h *handler) CreateMango(c *echo.Context) error {
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

	response, err := h.service.CreateMango(req)
	if err != nil {
		return mangoErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) GetMangoes(c *echo.Context) error {
	mangoes, err := h.service.GetMangoes()
	if err != nil {
		return mangoErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, mangoes)
}

func (h *handler) GetMangoByID(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid mango id",
			Details: err.Error(),
		})
	}

	response, err := h.service.GetMangoByID(uint(id))
	if err != nil {
		return mangoErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateMango(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid mango id",
			Details: err.Error(),
		})
	}

	var req dto.UpdateRequest

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

	response, err := h.service.UpdateMango(uint(id), req)
	if err != nil {
		return mangoErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, response)
}
