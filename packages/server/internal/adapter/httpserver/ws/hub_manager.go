package ws

import (
	"log"
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

	log.Printf("Getting hub for chatID: %s", chatID)
	if hub, exists := m.hubs[chatID]; exists {
		log.Println("Found existing hub")
		return hub
	}

	log.Println("Creating new hub")

	hub := NewHub()
	m.hubs[chatID] = hub
	go hub.Run()
	return hub
}
