package handler

import (
	"github.com/labstack/echo/v4"
	"mine-stats/jobs"
	"mine-stats/models"
	"mine-stats/protocol/minecraft"
	"net/http"
	"strconv"
	"time"
)

type ServerPost struct {
	ID      int           `json:"id" form:"id"`
	Name    string        `json:"name" form:"name"`
	URL     string        `json:"url" form:"url"`
	Port    uint16        `json:"port" form:"port"`
	Timeout time.Duration `json:"timeout" form:"timeout"`
	Every   time.Duration `json:"every" form:"every"`
}

func (h *Handler) ListOwnServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

	srvList, err := h.Store.GetMinecraftServerListByUser(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error":   "Could not get server list for your user try again later",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, srvList)
}

func (h *Handler) OneOwnServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

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

	server, err := h.Store.GetMinecraftServerForUserByID(user.ID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error":   err.Error(),
			"message": "Server Not Found",
		})
	}

	return c.JSON(http.StatusOK, server)
}

func (h *Handler) AddServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

	serverToAdd := new(ServerPost)
	if err := c.Bind(serverToAdd); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"Error":   err.Error(),
			"message": "Error getting correct struct",
		})
	}

	mineProto := minecraftProtocol.MinecraftServer{
		Name:    serverToAdd.Name,
		Every:   serverToAdd.Every,
		Timeout: serverToAdd.Timeout,
		Address: serverToAdd.URL,
		Port:    serverToAdd.Port,
	}

	srv, err := h.Store.AddServer(&mineProto, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error":   err.Error(),
			"message": "Error saving server",
		})
	}

	j := jobs.NewJob(srv)
	jobs.AddJob(j)

	return c.JSON(http.StatusOK, srv)
}

func (h *Handler) UpdateServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

	serverToUpdate := new(ServerPost)
	if err := c.Bind(serverToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"Error":   err.Error(),
			"message": "Error getting correct struct",
		})
	}

	updateData := models.Server{
		Url:     serverToUpdate.URL,
		Every:   serverToUpdate.Every,
		Port:    serverToUpdate.Port,
		Timeout: serverToUpdate.Timeout,
		Name:    serverToUpdate.Name,
	}

	srv, err := h.Store.UpdateServer(user.ID, serverToUpdate.ID, updateData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, srv)
}

func (h *Handler) DeleteServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

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

	err = h.Store.DeleteServerForUserByID(user.ID, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			"error":   "error deleting server, try again",
			"message": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
