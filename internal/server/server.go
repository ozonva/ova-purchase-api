package server

import "github.com/ozonva/ova-purchase-api/internal/disposal"

type Server interface {
	disposal.Disposal
	Run() error
}
