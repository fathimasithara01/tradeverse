package ws

// import (
// 	"sync"

// 	"github.com/gorilla/websocket"
// )

// type Client struct {
// 	ID   uint
// 	Conn *websocket.Conn
// }

// type Manager struct {
// 	Clients map[uint]*Client
// 	mu      sync.RWMutex
// }

// var WSManager = Manager{
// 	Clients: make(map[uint]*Client),
// }

// func (m *Manager) AddClient(userID uint, conn *websocket.Conn) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	m.Clients[userID] = &Client{
// 		ID:   userID,
// 		Conn: conn,
// 	}
// }

// func (m *Manager) RemoveClient(userID uint) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()
// 	delete(m.Clients, userID)
// }

// func (m *Manager) SendTo(userID uint, message string) error {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()

// 	client, exists := m.Clients[userID]
// 	if !exists {
// 		return nil // user not connected
// 	}
// 	return client.Conn.WriteMessage(websocket.TextMessage, []byte(message))
// }
