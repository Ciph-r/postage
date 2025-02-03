package server

import "github.com/ciph-r/postage/internal/services/health"

type config struct {
	Health health.Config `envPrefix:"HEALTH_"`
}
