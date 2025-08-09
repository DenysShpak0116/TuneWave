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

func (m *HubManager) GetHub(chatID string) *Hub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, exists := m.hubs[chatID]; exists {
		return hub
	}

	hub := NewHub()
	m.hubs[chatID] = hub
	go hub.Run()
	return hub
}
