package handler

import (
	"github.com/labstack/echo/v4"
	"mine-stats/models"
	"mine-stats/protocol/minecraft"
	"net/http"
	"strconv"
	"time"
)

type AddServerPost struct {
	//AddedBy int `json:"added_by" form:"added_by"`
	Name string `json:"name" form:"name"`
	URL string `json:"url" form:"url"`
	Port uint16 `json:"port" form:"port"`
	Timeout time.Duration `json:"timeout" form:"timeout"`
	Every time.Duration `json:"every" form:"every"`
}

func (h *Handler) ListServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) OneServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) AddServer(c echo.Context) error {
	user := c.Get("user").(*models.User)

	serverToAdd := new(AddServerPost)
	if err := c.Bind(serverToAdd); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
			"message": "Error getting correct struct",
		})
	}

	mineProto := minecraftProtocol.MinecraftServer{
		Name: serverToAdd.Name,
		Every: serverToAdd.Every,
		Timeout: serverToAdd.Timeout,
		Address: serverToAdd.URL,
		Port: serverToAdd.Port,
	}

	srv, err := h.Store.AddServer(&mineProto, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"Error": err.Error(),
			"message": "Error saving server",
		})
	}

	// TODO: Add server to job list

	return c.JSON(http.StatusOK, srv)
}

func (h *Handler) UpdateServer(c echo.Context) error {
	return c.HTML(http.StatusOK, "OK!")
}

func (h *Handler) DeleteServer(c echo.Context) error {
	idParam := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Could not parse server ID",
		})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Could not parse server ID",
		})
	}

	h.Store.DeleteServerByID(id)

	return c.HTML(http.StatusOK, c.Param("id"))
}
