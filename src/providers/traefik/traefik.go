package traefik

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
	"gopkg.in/yaml.v3"
)

type TraefikProvider struct {
	ConfigDir string
}

type DynamicConfig struct {
	HTTP struct {
		Routers  map[string]Router  `yaml:"routers"`
		Services map[string]Service `yaml:"services"`
	} `yaml:"http"`
}

type Router struct {
	Rule    string `yaml:"rule"`
	Service string `yaml:"service"`
	TLS     *TLS   `yaml:"tls,omitempty"`
}

type TLS struct {
	CertResolver string `yaml:"certResolver"`
}

type Service struct {
	LoadBalancer LoadBalancer `yaml:"loadBalancer"`
}

type LoadBalancer struct {
	Servers []Server `yaml:"servers"`
}

type Server struct {
	URL string `yaml:"url"`
}

func NewTraefikProvider(configDir string) *TraefikProvider {
	if configDir == "" {
		configDir = "./configs/traefik"
	}
	_ = os.MkdirAll(configDir, 0755)
	return &TraefikProvider{ConfigDir: configDir}
}

func (p *TraefikProvider) ConfigureIngress(ctx context.Context, domainName string, targetService string) error {
	log.Printf("[TraefikProvider] Configuring ingress for domain %s -> %s", domainName, targetService)

	cfg := &DynamicConfig{}
	routerName := fmt.Sprintf("%s-router", domainName)
	serviceName := fmt.Sprintf("%s-service", domainName)

	cfg.HTTP.Routers = make(map[string]Router)
	cfg.HTTP.Routers[routerName] = Router{
		Rule:    fmt.Sprintf("Host(`%s`)", domainName),
		Service: serviceName,
		TLS: &TLS{
			CertResolver: "le", // Default Let's Encrypt resolver
		},
	}

	cfg.HTTP.Services = make(map[string]Service)
	cfg.HTTP.Services[serviceName] = Service{
		LoadBalancer: LoadBalancer{
			Servers: []Server{
				{URL: fmt.Sprintf("http://%s", targetService)},
			},
		},
	}

	return p.saveConfig(domainName, cfg)
}

func (p *TraefikProvider) RequestCertificate(ctx context.Context, domainName string) error {
	log.Printf("[TraefikProvider] SSL automation initiated for %s", domainName)
	return nil // Handled by acme in Traefik router config
}

func (p *TraefikProvider) saveConfig(name string, cfg *DynamicConfig) error {
	filename := filepath.Join(p.ConfigDir, fmt.Sprintf("%s.yaml", name))
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (p *TraefikProvider) GetNetworkStatus(ctx context.Context, domainName string) (string, error) {
	return "active", nil
}

// Ensure TraefikProvider implements domain.NetworkProvider
var _ domain.NetworkProvider = (*TraefikProvider)(nil)
