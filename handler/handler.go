package handler

import (
	"mine-stats/store"
)

type Response map[string]interface{}

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
