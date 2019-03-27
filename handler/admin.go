package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AdminListServer(c echo.Context) error {
	srvs, err := h.Store.GetMinecraftServerList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, srvs)
}

func (h *Handler) AdminOneServer(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}

	srv, err := h.Store.GetMinecraftServerByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, srv)
}

func (h *Handler) AdminDeleteServer(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}

	err = h.Store.DeleteServerByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) AdminListUser(c echo.Context) error {
	users, err := h.Store.GetUserList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) AdminOneUser(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}

	user, err := h.Store.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) AdminDeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	if idParam == "" {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error": "Could not parse server ID",
		})
	}

	err = h.Store.DeleteUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
