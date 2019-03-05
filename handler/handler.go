package handler

import (
	"log"
	"mine-stats/public"
	"mine-stats/store"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Store *store.Store
	Prod  bool
}

func NewHandler(st *store.Store, prod bool) *Handler {
	return &Handler{
		Store: st,
		Prod:  prod,
	}
}

func (h *Handler) ServeIndex(c echo.Context) error {
	htmlb, err := public.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	// convert to string
	html := string(htmlb)
	return c.HTML(http.StatusOK, html)
}
