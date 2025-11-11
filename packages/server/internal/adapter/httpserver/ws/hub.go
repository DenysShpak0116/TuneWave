package ws

import (
	"encoding/base64"
	"encoding/json"
	"math/big"

	diffiehelman "github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers/math"
)

type MyBigInt struct {
	*big.Int
}

func (m *MyBigInt) MarshalJSON() ([]byte, error) {
	if m.Int == nil {
		return []byte(`""`), nil
	}

	bytes := m.Int.Bytes()
	encoded := base64.URLEncoding.EncodeToString(bytes)
	return json.Marshal(encoded)
}

func (m *MyBigInt) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	if encoded == "" {
		m.Int = new(big.Int)
		return nil
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	m.Int = new(big.Int)
	m.Int.SetBytes(bytes)

	return nil
}

type DhKeyData struct {
	Prime     *MyBigInt `json:"q"`
	Generator *MyBigInt `json:"p"`
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	DhKeys *DhKeyData
}

func NewHub() (*Hub, error) {
	q, p, err := diffiehelman.GetRfc3526Group2048()
	if err != nil {
		return nil, err
	}

	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		DhKeys: &DhKeyData{
			Prime:     &MyBigInt{q},
			Generator: &MyBigInt{p},
		},
	}, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				client.Send <- message
			}
		}
	}
}
