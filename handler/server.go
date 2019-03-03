package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) ListServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) AddServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) UpdateServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) DeleteServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}
