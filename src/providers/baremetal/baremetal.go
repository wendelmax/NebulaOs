package baremetal

import (
	"context"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type BareMetalProvider struct {
	SSHHost     string
	SSHUser     string
	SSHKeyPath  string
	SSHPassword string
	IPMIUser    string
	IPMIPassword string
}

func NewBareMetalProvider(sshHost, sshUser, sshKeyPath, sshPassword, ipmiUser, ipmiPassword string) *BareMetalProvider {
	return &BareMetalProvider{
		SSHHost:      sshHost,
		SSHUser:      sshUser,
		SSHKeyPath:   sshKeyPath,
		SSHPassword:  sshPassword,
		IPMIUser:     ipmiUser,
		IPMIPassword: ipmiPassword,
	}
}

func (p *BareMetalProvider) connectSSH() (*SSHClient, error) {
	host := p.SSHHost
	port := "22"
	if h, p, err := net.SplitHostPort(host); err == nil {
		host = h
		port = p
	}
	return Connect(host, port, p.SSHUser, p.SSHKeyPath, p.SSHPassword)
}

func (p *BareMetalProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	name := getString(resource.Metadata, "name", resource.Name)
	vcpus := getInt(resource.Metadata, "vcpus", 2)
	memory := getInt(resource.Metadata, "memory", 2048)
	diskSize := getInt(resource.Metadata, "disk_size", 20)
	diskPath := getString(resource.Metadata, "disk_path", fmt.Sprintf("/var/lib/libvirt/images/%s.qcow2", name))
	iso := getString(resource.Metadata, "iso", "")
	network := getString(resource.Metadata, "network", "virbr0")

	cmd := fmt.Sprintf(
		"virt-install --name %s --vcpus %d --memory %d --disk path=%s,size=%d --network bridge=%s --graphics none --noautoconsole 2>&1",
		name, vcpus, memory, diskPath, diskSize, network,
	)
	if iso != "" {
		cmd += fmt.Sprintf(" --cdrom %s", iso)
	}

	out, err := client.Run(cmd)
	if err != nil {
		return fmt.Errorf("baremetal: virt-install failed: %w\n%s", err, out)
	}
	log.Printf("[baremetal] VM %s provisioned: %s", name, strings.TrimSpace(out))

	dumpOut, err := client.Run(fmt.Sprintf("virsh dumpxml %s", name))
	if err != nil {
		return fmt.Errorf("baremetal: failed to get domain XML: %w", err)
	}

	resource.ProviderID = extractUUID(dumpOut)
	if resource.ProviderID == "" {
		resource.ProviderID = name
	}
	resource.State = "running"
	return nil
}

func (p *BareMetalProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	if _, err := client.Run(fmt.Sprintf("virsh destroy %s 2>&1", resource.ProviderID)); err != nil {
		log.Printf("[baremetal] virsh destroy warning (may already be off): %s", err)
	}

	if _, err := client.Run(fmt.Sprintf("virsh undefine %s 2>&1", resource.ProviderID)); err != nil {
		return fmt.Errorf("baremetal: virsh undefine failed: %w", err)
	}

	resource.State = "decommissioned"
	return nil
}

func (p *BareMetalProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	client, err := p.connectSSH()
	if err != nil {
		return "", fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	out, err := client.Run(fmt.Sprintf("virsh dominfo %s 2>&1", resourceID))
	if err != nil {
		return "unknown", fmt.Errorf("baremetal: virsh dominfo failed: %w", err)
	}
	return extractState(out), nil
}

func (p *BareMetalProvider) Ping(ctx context.Context) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ping failed to connect: %w", err)
	}
	defer client.Close()

	_, err = client.Run("virsh version")
	if err != nil {
		return fmt.Errorf("baremetal: virsh not available: %w", err)
	}
	return nil
}

func (p *BareMetalProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	client, err := p.connectSSH()
	if err != nil {
		return nil, fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	out, err := client.Run("ls -1 /var/lib/libvirt/images/*.qcow2 /var/lib/libvirt/boot/*.iso 2>/dev/null")
	if err != nil {
		return []domain.Image{}, nil
	}

	var images []domain.Image
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		imgType := "template"
		if strings.HasSuffix(line, ".iso") {
			imgType = "iso"
		}
		images = append(images, domain.Image{
			ID:       line,
			Name:     filepath.Base(line),
			Provider: "baremetal",
			Type:     imgType,
		})
	}
	return images, nil
}

func (p *BareMetalProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	return fmt.Errorf("baremetal: direct container deployment is not supported")
}

func (p *BareMetalProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	cmd := fmt.Sprintf("brctl addbr %s 2>&1 && ip link set %s up 2>&1", n.Name, n.Name)
	if n.CIDR != "" {
		cmd += fmt.Sprintf(" && ip addr add %s dev %s 2>&1", n.CIDR, n.Name)
	}
	_, err = client.Run(cmd)
	if err != nil {
		return fmt.Errorf("baremetal: failed to configure bridge %s: %w", n.Name, err)
	}
	log.Printf("[baremetal] Bridge %s configured (%s)", n.Name, n.CIDR)
	return nil
}

func (p *BareMetalProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	cmd := fmt.Sprintf("ip route add %s via %s 2>&1", r.Destination, r.NextHop)
	_, err = client.Run(cmd)
	if err != nil {
		return fmt.Errorf("baremetal: failed to add route %s via %s: %w", r.Destination, r.NextHop, err)
	}
	return nil
}

func (p *BareMetalProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	client, err := p.connectSSH()
	if err != nil {
		return fmt.Errorf("baremetal: ssh connect failed: %w", err)
	}
	defer client.Close()

	cmd := fmt.Sprintf(
		"iptables -C FORWARD -m comment --comment 'sg-%s' -j ACCEPT 2>/dev/null || iptables -A FORWARD -m comment --comment 'sg-%s' -j ACCEPT 2>&1",
		sgID, sgID,
	)
	_, err = client.Run(cmd)
	if err != nil {
		return fmt.Errorf("baremetal: failed to attach security group %s: %w", sgID, err)
	}
	return nil
}

func getString(m map[string]interface{}, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func getInt(m map[string]interface{}, key string, def int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		case string:
			fmt.Sscanf(val, "%d", &def)
			return def
		}
	}
	return def
}

func extractUUID(xml string) string {
	start := strings.Index(xml, "<uuid>")
	if start == -1 {
		return ""
	}
	start += 6
	end := strings.Index(xml[start:], "</uuid>")
	if end == -1 {
		return ""
	}
	return xml[start : start+end]
}

func extractState(dominfo string) string {
	for _, line := range strings.Split(dominfo, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "State:") {
			state := strings.TrimSpace(strings.TrimPrefix(line, "State:"))
			switch {
			case strings.Contains(state, "running"):
				return "running"
			case strings.Contains(state, "shut off"), strings.Contains(state, "shutdown"):
				return "stopped"
			case strings.Contains(state, "paused"):
				return "paused"
			default:
				return "unknown"
			}
		}
	}
	return "unknown"
}
