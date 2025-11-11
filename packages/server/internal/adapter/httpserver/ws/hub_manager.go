package ws

import (
	"sync"
)

type HubManager struct {
	hubs map[string]*Hub
	mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (m *HubManager) GetHub(chatID string) (*Hub, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, exists := m.hubs[chatID]; exists {
		return hub, nil
	}

	hub, err := NewHub()
	if err != nil {
		return nil, err
	}

	m.hubs[chatID] = hub
	go hub.Run()
	return hub, nil
}
