package handler

import "mine-stats/store"

type Handler struct {
	store *store.Store
}

func NewHandler(st *store.Store) *Handler {
	return &Handler{
		store: st,
	}
}