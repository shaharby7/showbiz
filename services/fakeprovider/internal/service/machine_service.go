package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shaharby7/showbiz/services/fakeprovider/internal/kubevirt"
	"github.com/shaharby7/showbiz/services/fakeprovider/internal/model"
)

type MachineService struct {
	mu       sync.RWMutex
	machines map[string]*model.Machine
	kv       *kubevirt.Client
}

func NewMachineService(kv *kubevirt.Client) *MachineService {
	return &MachineService{
		machines: make(map[string]*model.Machine),
		kv:       kv,
	}
}

type CreateMachineInput struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	CPU       int    `json:"cpu"`
	MemoryMB  int    `json:"memoryMB"`
	Image     string `json:"image"`
}

type UpdateMachineInput struct {
	CPU      *int `json:"cpu,omitempty"`
	MemoryMB *int `json:"memoryMB,omitempty"`
}

func (s *MachineService) Create(input CreateMachineInput) *model.Machine {
	now := time.Now().UTC()
	m := &model.Machine{
		ID:        uuid.New().String(),
		Name:      input.Name,
		Namespace: input.Namespace,
		CPU:       input.CPU,
		MemoryMB:  input.MemoryMB,
		Image:     input.Image,
		Status:    model.StatusInitialized,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mu.Lock()
	s.machines[m.ID] = m
	s.mu.Unlock()

	go s.provision(m.ID)

	return m
}

func (s *MachineService) provision(id string) {
	s.mu.Lock()
	m, ok := s.machines[id]
	if !ok {
		s.mu.Unlock()
		return
	}
	m.Status = model.StatusProvisioning
	m.UpdatedAt = time.Now().UTC()
	name := m.Name
	namespace := m.Namespace
	cpu := m.CPU
	memoryMB := m.MemoryMB
	image := m.Image
	s.mu.Unlock()

	ctx := context.Background()

	if err := s.kv.CreateVMI(ctx, name, namespace, cpu, memoryMB, image); err != nil {
		slog.Error("failed to create VMI", "id", id, "error", err)
		s.setStatus(id, model.StatusFailed, "")
		return
	}

	// Poll for VMI to become Running with an IP
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-ticker.C:
			status, err := s.kv.GetVMIStatus(ctx, name, namespace)
			if err != nil {
				slog.Error("failed to get VMI status", "id", id, "error", err)
				continue
			}
			if status.Phase == "Running" && status.IP != "" {
				s.setStatus(id, model.StatusReady, status.IP)
				return
			}
		case <-timeout:
			slog.Error("VMI provision timed out", "id", id)
			s.setStatus(id, model.StatusFailed, "")
			return
		}
	}
}

func (s *MachineService) setStatus(id, status, ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.machines[id]; ok {
		m.Status = status
		if ip != "" {
			m.IP = ip
		}
		m.UpdatedAt = time.Now().UTC()
	}
}

func (s *MachineService) Get(id string) (*model.Machine, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.machines[id]
	if !ok {
		return nil, fmt.Errorf("machine not found")
	}
	copy := *m
	return &copy, nil
}

func (s *MachineService) List() []*model.Machine {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*model.Machine, 0, len(s.machines))
	for _, m := range s.machines {
		copy := *m
		result = append(result, &copy)
	}
	return result
}

func (s *MachineService) Update(id string, input UpdateMachineInput) (*model.Machine, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.machines[id]
	if !ok {
		return nil, fmt.Errorf("machine not found")
	}
	if input.CPU != nil {
		m.CPU = *input.CPU
	}
	if input.MemoryMB != nil {
		m.MemoryMB = *input.MemoryMB
	}
	m.UpdatedAt = time.Now().UTC()
	copy := *m
	return &copy, nil
}

func (s *MachineService) Delete(id string) error {
	s.mu.Lock()
	m, ok := s.machines[id]
	if !ok {
		s.mu.Unlock()
		return fmt.Errorf("machine not found")
	}
	name := m.Name
	namespace := m.Namespace
	m.Status = model.StatusDeleted
	m.UpdatedAt = time.Now().UTC()
	delete(s.machines, id)
	s.mu.Unlock()

	ctx := context.Background()
	if err := s.kv.DeleteVMI(ctx, name, namespace); err != nil {
		slog.Error("failed to delete VMI", "id", id, "error", err)
		return fmt.Errorf("delete VMI: %w", err)
	}
	return nil
}
