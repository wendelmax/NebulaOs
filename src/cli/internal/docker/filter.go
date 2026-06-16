package docker

import (
	"strings"

	"github.com/docker/docker/api/types/container"
)

func IsNebulaContainer(c container.Summary) bool {
	for _, name := range c.Names {
		n := strings.ToLower(name)
		if strings.Contains(n, "nebula") ||
			strings.Contains(n, "nats") ||
			strings.Contains(n, "traefik") ||
			strings.Contains(n, "postgres") ||
			strings.Contains(n, "keycloak") ||
			strings.Contains(n, "vault") ||
			strings.Contains(n, "moto") {
			return true
		}
	}
	return false
}
