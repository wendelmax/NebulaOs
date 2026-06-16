package proxmox

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type ProxmoxProvider struct {
	BaseURL string
	Token   string
	client  *Client
}

func NewProxmoxProvider(baseURL, token string) *ProxmoxProvider {
	client := NewClient(baseURL)
	if strings.Contains(token, "!") {
		client.WithAPIToken(token)
	} else {
		client.WithCredentials("root@pam", token)
	}
	return &ProxmoxProvider{
		BaseURL: baseURL,
		Token:   token,
		client:  client,
	}
}

func (p *ProxmoxProvider) Ping(ctx context.Context) error {
	ver, err := p.client.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("proxmox health check: %w", err)
	}
	log.Printf("[proxmox] health check OK — version %s", ver)
	return nil
}

func (p *ProxmoxProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	vmid, err := p.client.GetNextID(ctx)
	if err != nil {
		return fmt.Errorf("proxmox get nextid: %w", err)
	}

	node, err := p.resolveNode(ctx, resource.Metadata)
	if err != nil {
		return fmt.Errorf("proxmox resolve node: %w", err)
	}

	meta := resource.Metadata
	if meta == nil {
		meta = make(map[string]interface{})
	}

	params := map[string]interface{}{
		"vmid":   vmid,
		"name":   resource.Name,
		"cores":  getInt(meta, "cores", 2),
		"memory": getInt(meta, "memory", 2048),
		"ostype": getString(meta, "ostype", "l26"),
		"net0":  getString(meta, "net", "virtio,bridge=vmbr0"),
		"scsihw": "virtio-scsi-single",
	}

	if disk := getInt(meta, "disk", 0); disk > 0 {
		storage := getString(meta, "storage", "local-lvm")
		params["scsi0"] = fmt.Sprintf("%s:%d,format=qcow2", storage, disk)
	}

	if iso := getString(meta, "iso", ""); iso != "" {
		params["ide2"] = iso + ",media=cdrom"
		params["boot"] = "order=ide2;scsi0;net0"
	} else {
		params["boot"] = "order=scsi0;net0"
	}

	if ciuser := getString(meta, "ciuser", ""); ciuser != "" {
		params["ciuser"] = ciuser
	}
	if cipass := getString(meta, "cipassword", ""); cipass != "" {
		params["cipassword"] = cipass
	}
	if sshkey := getString(meta, "sshkey", ""); sshkey != "" {
		params["sshkeys"] = sshkey
	}
	if searchdomain := getString(meta, "searchdomain", ""); searchdomain != "" {
		params["searchdomain"] = searchdomain
	}
	if nameserver := getString(meta, "nameserver", ""); nameserver != "" {
		params["nameserver"] = nameserver
	}
	if ipconfig := getString(meta, "ipconfig0", ""); ipconfig != "" {
		params["ipconfig0"] = ipconfig
	}

	if err := p.client.CreateVM(ctx, node, params); err != nil {
		return fmt.Errorf("proxmox create vm: %w", err)
	}

	resource.ProviderID = fmt.Sprintf("pve-%d", vmid)
	if resource.Metadata == nil {
		resource.Metadata = make(map[string]interface{})
	}
	resource.Metadata["node"] = node
	resource.Metadata["vmid"] = vmid
	resource.State = "running"
	log.Printf("[proxmox] provisioned VM %d on node %s", vmid, node)
	return nil
}

func (p *ProxmoxProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	vmid, node, err := p.extractVMID(ctx, resource)
	if err != nil {
		return fmt.Errorf("proxmox decommission: %w", err)
	}

	if err := p.client.DeleteVM(ctx, node, vmid); err != nil {
		log.Printf("[proxmox] DeleteVM failed, trying DeleteLXC: %v", err)
		if err2 := p.client.DeleteLXC(ctx, node, vmid); err2 != nil {
			return fmt.Errorf("proxmox delete failed (tried qemu and lxc): %w", err2)
		}
	}

	resource.State = "deleted"
	log.Printf("[proxmox] decommissioned VM/LXC %d on node %s", vmid, node)
	return nil
}

func (p *ProxmoxProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	vmid, node, err := p.parseProviderID(ctx, resourceID)
	if err != nil {
		return "", fmt.Errorf("proxmox get status: %w", err)
	}

	status, err := p.client.GetVMStatus(ctx, node, vmid)
	if err == nil {
		return status, nil
	}

	status, err = p.client.GetLXCStatus(ctx, node, vmid)
	if err != nil {
		return "", fmt.Errorf("proxmox get status (tried qemu and lxc): %w", err)
	}
	return status, nil
}

func (p *ProxmoxProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	log.Printf("[proxmox] AttachSecurityGroup not yet implemented for resource %s, sg %s", resourceID, sgID)
	return nil
}

func (p *ProxmoxProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	nodes, err := p.client.GetNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("proxmox list images: %w", err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("proxmox list images: no nodes available")
	}

	var images []domain.Image
	seen := make(map[string]bool)

	for _, node := range nodes {
		storages, err := p.client.GetNodeStorages(ctx, node)
		if err != nil {
			log.Printf("[proxmox] list storages on node %s: %v", node, err)
			continue
		}
		for _, storage := range storages {
			contents, err := p.client.GetStorageContent(ctx, node, storage)
			if err != nil {
				log.Printf("[proxmox] list content on %s/%s: %v", node, storage, err)
				continue
			}
			for _, c := range contents {
				if c.Content != "iso" && c.Content != "vztmpl" {
					continue
				}
				if seen[c.Volume] {
					continue
				}
				seen[c.Volume] = true
				imgType := "iso"
				if c.Content == "vztmpl" {
					imgType = "template"
				}
				name := c.Volume
				if c.Name != "" {
					name = c.Name
				}
				images = append(images, domain.Image{
					ID:       c.Volume,
					Name:     name,
					Provider: "proxmox",
					Type:     imgType,
				})
			}
		}
	}

	if len(images) == 0 {
		images = append(images, domain.Image{
			ID: "debian-12-iso", Name: "Debian 12 ISO", Provider: "proxmox", Type: "iso",
		})
	}

	return images, nil
}

func (p *ProxmoxProvider) DeployContainer(ctx context.Context, container *domain.Container) error {
	nodes, err := p.client.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("proxmox deploy container: %w", err)
	}
	if len(nodes) == 0 {
		return fmt.Errorf("proxmox deploy container: no nodes available")
	}
	node := nodes[0]

	vmid, err := p.client.GetNextID(ctx)
	if err != nil {
		return fmt.Errorf("proxmox deploy container nextid: %w", err)
	}

	params := map[string]interface{}{
		"vmid":     vmid,
		"hostname": container.Name,
		"ostemplate": container.Image,
		"cores":   int(container.CPU),
		"memory":  container.MemoryMB,
		"net0":    "name=eth0,bridge=vmbr0,ip=dhcp",
		"storage": "local-lvm",
		"password": "nebulaos",
		"swap":     256,
	}

	if err := p.client.CreateLXC(ctx, node, params); err != nil {
		return fmt.Errorf("proxmox create lxc: %w", err)
	}

	container.State = "running"
	log.Printf("[proxmox] deployed LXC container %d on node %s", vmid, node)
	return nil
}

func (p *ProxmoxProvider) ConfigureNetwork(ctx context.Context, network *domain.Network) error {
	nodes, err := p.client.GetNodes(ctx)
	if err != nil {
		return fmt.Errorf("proxmox configure network: %w", err)
	}
	if len(nodes) == 0 {
		return fmt.Errorf("proxmox configure network: no nodes available")
	}
	node := nodes[0]

	params := map[string]interface{}{
		"type":   "bridge",
		"name":   network.Name,
		"cidr":   network.CIDR,
		"gateway": network.Gateway,
		"autostart": true,
		"comments":  "Created by NebulaOS",
	}

	if err := p.client.CreateNetwork(ctx, node, params); err != nil {
		return fmt.Errorf("proxmox create network: %w", err)
	}

	network.State = "active"
	log.Printf("[proxmox] configured network bridge %s on node %s", network.Name, node)
	return nil
}

func (p *ProxmoxProvider) AddRoute(ctx context.Context, route *domain.Route) error {
	log.Printf("[proxmox] AddRoute not yet implemented — route %s -> %s", route.Destination, route.NextHop)
	return nil
}

func (p *ProxmoxProvider) resolveNode(ctx context.Context, meta map[string]interface{}) (string, error) {
	if node := getString(meta, "node", ""); node != "" {
		return node, nil
	}
	nodes, err := p.client.GetNodes(ctx)
	if err != nil {
		return "", err
	}
	if len(nodes) == 0 {
		return "", fmt.Errorf("no proxmox nodes available")
	}
	return nodes[0], nil
}

func (p *ProxmoxProvider) extractVMID(ctx context.Context, resource *domain.Resource) (int, string, error) {
	if resource.Metadata != nil {
		if v, ok := resource.Metadata["vmid"]; ok {
			switch vv := v.(type) {
			case float64:
				return int(vv), getString(resource.Metadata, "node", ""), nil
			case int:
				return vv, getString(resource.Metadata, "node", ""), nil
			}
		}
	}
	vmid, node, err := p.parseProviderID(ctx, resource.ProviderID)
	if err != nil {
		return 0, "", err
	}
	if node == "" {
		return 0, "", fmt.Errorf("no node found for resource %s", resource.ID)
	}
	return vmid, node, nil
}

func (p *ProxmoxProvider) parseProviderID(ctx context.Context, providerID string) (int, string, error) {
	if providerID == "" {
		return 0, "", fmt.Errorf("empty provider ID")
	}

	id := strings.TrimPrefix(providerID, "pve-")
	vmid, err := strconv.Atoi(id)
	if err != nil {
		return 0, "", fmt.Errorf("parse vmid from %q: %w", providerID, err)
	}

	nodes, err := p.client.GetNodes(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("lookup nodes: %w", err)
	}
	if len(nodes) == 0 {
		return 0, "", fmt.Errorf("no proxmox nodes available")
	}

	return vmid, nodes[0], nil
}

func getString(m map[string]interface{}, key, def string) string {
	if m == nil {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	s, ok := v.(string)
	if !ok {
		return def
	}
	return s
}

func getInt(m map[string]interface{}, key string, def int) int {
	if m == nil {
		return def
	}
	v, ok := m[key]
	if !ok {
		return def
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case string:
		i, err := strconv.Atoi(n)
		if err != nil {
			return def
		}
		return i
	}
	return def
}
