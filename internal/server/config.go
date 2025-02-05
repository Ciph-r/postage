package server

import (
	"github.com/ciph-r/postage/internal/services/health"
	"github.com/ciph-r/postage/internal/services/sockets"
)

// config is the main config of the entire server. All other configs are imported
// here and given a corresponding namespace using the envPrefix tag.
type config struct {
	Health  health.Config  `envPrefix:"HEALTH_"`
	Sockets sockets.Config `envPrefix:"SOCKETS_"`
}
