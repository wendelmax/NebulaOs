package services

import (
	"context"
	"fmt"
	"time"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type BareMetalManager struct {
	repo domain.BareMetalRepository
}

func NewBareMetalManager(repo domain.BareMetalRepository) *BareMetalManager {
	return &BareMetalManager{repo: repo}
}

func (m *BareMetalManager) RegisterNode(ctx context.Context, node *domain.BareMetalNode) error {
	node.CreatedAt = time.Now()
	node.State = domain.NodeStateAvailable
	return m.repo.Create(ctx, node)
}

func (m *BareMetalManager) ProvisionNode(ctx context.Context, nodeID string, blueprintID string) error {
	node, err := m.repo.GetByID(ctx, nodeID)
	if err != nil {
		return err
	}

	if node.State != domain.NodeStateAvailable {
		return fmt.Errorf("node is not available for provisioning (current state: %s)", node.State)
	}

	node.State = domain.NodeStateProvisioning
	m.repo.Update(ctx, node)

	m.repo.AddLog(ctx, &domain.ProvisioningLog{
		ID:        domain.NewID(),
		NodeID:    nodeID,
		Message:   "Starting iPXE provisioning sequence",
		Level:     "info",
		Timestamp: time.Now(),
	})

	// Async provisioning simulation
	go m.simulateProvisioning(nodeID)

	return nil
}

func (m *BareMetalManager) simulateProvisioning(nodeID string) {
	ctx := context.Background()
	time.Sleep(5 * time.Second)
	m.repo.AddLog(ctx, &domain.ProvisioningLog{
		ID:        domain.NewID(),
		NodeID:    nodeID,
		Message:   "Powering on via Redfish/IPMI",
		Level:     "info",
		Timestamp: time.Now(),
	})

	time.Sleep(10 * time.Second)
	m.repo.AddLog(ctx, &domain.ProvisioningLog{
		ID:        domain.NewID(),
		NodeID:    nodeID,
		Message:   "PXE Boot: Loading kernel and initrd",
		Level:     "info",
		Timestamp: time.Now(),
	})

	time.Sleep(15 * time.Second)
	node, _ := m.repo.GetByID(ctx, nodeID)
	node.State = domain.NodeStateActive
	m.repo.Update(ctx, node)

	m.repo.AddLog(ctx, &domain.ProvisioningLog{
		ID:        domain.NewID(),
		NodeID:    nodeID,
		Message:   "Provisioning completed. Node is Active.",
		Level:     "info",
		Timestamp: time.Now(),
	})
}

func (m *BareMetalManager) ListNodes(ctx context.Context) ([]*domain.BareMetalNode, error) {
	return m.repo.List(ctx)
}

func (m *BareMetalManager) GetNodeStatus(ctx context.Context, nodeID string) (*domain.BareMetalNode, []*domain.ProvisioningLog, error) {
	node, err := m.repo.GetByID(ctx, nodeID)
	if err != nil {
		return nil, nil, err
	}
	logs, _ := m.repo.GetLogs(ctx, nodeID)
	return node, logs, nil
}
