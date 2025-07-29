package unlink

import (
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/unlink"
	"github.com/go-chi/chi/v5"
)

func InitUnlinkHanldler(router chi.Router, handler Inbound.UnlinkHandler)
