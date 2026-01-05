package game

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for dev
	},
}

type Manager struct {
	rooms map[string]*Room
	dict  *Dictionary
	um    *UserManager
}

func NewManager(dict *Dictionary, um *UserManager) *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
		dict:  dict,
		um:    um,
	}
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
	// Parse Query Params
	query := r.URL.Query()
	name := query.Get("name")
	roomID := query.Get("room")
	
	if name == "" {
		name = "Guest"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create or Join Room
	if roomID == "" {
		roomID = uuid.New().String()[:6] // Short ID
	}
	
	room, ok := m.rooms[roomID]
	if !ok {
		room = NewRoom(roomID, m.dict, m.um)
		m.rooms[roomID] = room
		go room.Run()
	}

	playerID := uuid.New().String()
	player := NewPlayer(playerID, name, PlayerHuman, conn)
	
	room.Register <- player

	// Send ID and RoomID to client
	welcomeMsg := map[string]interface{}{
		"type": "WELCOME",
		"payload": map[string]string{
			"id": player.ID,
			"roomId": roomID,
		},
	}
	if err := player.Conn.WriteJSON(welcomeMsg); err != nil {
		log.Println("Error sending welcome:", err)
		return
	}

	go m.writePump(player)
	go m.readPump(player, room)
}

func (m *Manager) readPump(p *Player, r *Room) {
	defer func() {
		r.Unregister <- p
		p.Conn.Close()
	}()
	
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			break
		}
		
		var action ActionMessage
		if err := json.Unmarshal(message, &action); err == nil {
			action.PlayerID = p.ID
			r.Action <- &action
		}
	}
}

func (m *Manager) writePump(p *Player) {
	defer p.Conn.Close()
	for {
		message, ok := <-p.Send
		if !ok {
			p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		
		w, err := p.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)
		
		if err := w.Close(); err != nil {
			return
		}
	}
}
