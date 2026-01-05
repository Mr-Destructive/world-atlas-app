package game

import (
	"github.com/gorilla/websocket"
)

type PlayerType int

const (
	PlayerHuman PlayerType = iota
	PlayerBot
)

type Player struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           PlayerType      `json:"type"`
	Conn           *websocket.Conn `json:"-"`
	Score          int             `json:"score"`
	Lives          int             `json:"lives"`
	IsTurn         bool            `json:"isTurn"`
	AvatarURL      string          `json:"avatarUrl"`
	MostUsedPlaces map[string]int  `json:"mostUsedPlaces"`
	// Channel to send messages to this player
	Send chan []byte `json:"-"`
}

func NewPlayer(id, name string, pType PlayerType, conn *websocket.Conn) *Player {
	return &Player{
		ID:             id,
		Name:           name,
		Type:           pType,
		Conn:           conn,
		Lives:          3,
		MostUsedPlaces: make(map[string]int),
		Send:           make(chan []byte, 256),
	}
}