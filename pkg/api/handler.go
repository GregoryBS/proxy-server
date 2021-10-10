package api

import (
	"proxy/pkg/database"

	"github.com/aerogo/aero"
)

type Handler struct {
	db *database.DB
}

func (h *Handler) RequestAll(ctx aero.Context) error {
	return nil
}

func (h *Handler) RequestOne(ctx aero.Context) error {
	return nil
}

func (h *Handler) Repeat(ctx aero.Context) error {
	return nil
}

func (h *Handler) Scan(ctx aero.Context) error {
	return nil
}
