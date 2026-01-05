package game

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type GameState string

const (
	StateWaiting GameState = "WAITING"
	StatePlaying GameState = "PLAYING"
	StateEnded   GameState = "ENDED"
)

type Move struct {
	PlayerID   string `json:"playerId"`
	PlayerName string `json:"playerName"`
	Word       string `json:"word"`
	Type       string `json:"type"` // City, Country, etc.
	Timestamp  int64  `json:"timestamp"`
}

type Room struct {
	ID               string
	Players          map[string]*Player
	TurnOrder        []string
	CurrentTurnIndex int
	State            GameState
	UsedWords        map[string]bool
	LastWord         string
	History          []Move `json:"history"`
	Round            int    `json:"round"`
	
	Mode     string         `json:"mode"`     // CLASSIC, POINT_RUSH, SUDDEN_DEATH
	Settings map[string]int `json:"settings"` // e.g., "timeLimit": 300

	Dict        *Dictionary
	BotBrain    *Bot
	UserManager *UserManager
	TurnStartTime time.Time

	Register   chan *Player
	Unregister chan *Player
	Broadcast  chan []byte
	Action     chan *ActionMessage

	mu sync.RWMutex
}

type ActionMessage struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	PlayerID string          `json:"-"`
}

func NewRoom(id string, dict *Dictionary, um *UserManager) *Room {
	return &Room{
		ID:          id,
		Players:     make(map[string]*Player),
		UsedWords:   make(map[string]bool),
		State:       StateWaiting,
		Mode:        "CLASSIC",
		Settings:    make(map[string]int),
		Dict:        dict,
		BotBrain:    NewBot(dict),
		UserManager: um,
		Register:    make(chan *Player),
		Unregister:  make(chan *Player),
		Broadcast:   make(chan []byte),
		Action:      make(chan *ActionMessage),
		History:     []Move{},
		Round:       1,
	}
}

func (r *Room) Run() {
	for {
		select {
		case player := <-r.Register:
			r.mu.Lock()
			r.Players[player.ID] = player
			r.TurnOrder = append(r.TurnOrder, player.ID)
			r.mu.Unlock()
			r.broadcastState()

		case player := <-r.Unregister:
			r.mu.Lock()
			if _, ok := r.Players[player.ID]; ok {
				delete(r.Players, player.ID)
				close(player.Send)

				removedIndex := -1
				for i, pid := range r.TurnOrder {
					if pid == player.ID {
						removedIndex = i
						break
					}
				}

				if removedIndex != -1 {
					r.TurnOrder = append(r.TurnOrder[:removedIndex], r.TurnOrder[removedIndex+1:]...)
					
					if len(r.TurnOrder) == 0 {
						r.CurrentTurnIndex = 0
						r.State = StateWaiting
					} else {
						if r.CurrentTurnIndex > removedIndex {
							r.CurrentTurnIndex--
						} else if r.CurrentTurnIndex == removedIndex {
							r.CurrentTurnIndex = r.CurrentTurnIndex % len(r.TurnOrder)
						}
					}
				}

				if r.State == StatePlaying {
					r.checkGameOver()
				}
			}
			r.mu.Unlock()
			r.broadcastState()

		case message := <-r.Broadcast:
			r.mu.RLock()
			for _, player := range r.Players {
				if player.Type == PlayerHuman {
					select {
					case player.Send <- message:
					default:
						close(player.Send)
						delete(r.Players, player.ID)
					}
				}
			}
			r.mu.RUnlock()

		case action := <-r.Action:
			r.handleAction(action)
		}
	}
}

func (r *Room) handleAction(action *ActionMessage) {
	switch action.Type {
	case "START_GAME":
		var p struct {
			Mode     string         `json:"mode"`
			Settings map[string]int `json:"settings"`
		}
		if err := json.Unmarshal(action.Payload, &p); err == nil {
			r.startGame(p.Mode, p.Settings)
		} else {
			r.startGame("CLASSIC", nil)
		}
	case "ADD_BOT":
		r.addBot()
	case "SUBMIT_WORD":
		var p struct{ Word string `json:"word"` }
		if err := json.Unmarshal(action.Payload, &p); err == nil {
			r.processTurn(action.PlayerID, p.Word)
		}
	case "BOT_MOVE":
		r.processBotTurn(action.PlayerID)
	}
}

func (r *Room) addBot() {
	r.mu.Lock()
	defer r.mu.Unlock()

	botID := uuid.New().String()
	name := fmt.Sprintf("Bot-%s", botID[:4])

	botPlayer := NewPlayer(botID, name, PlayerBot, nil)
	r.Players[botID] = botPlayer
	r.TurnOrder = append(r.TurnOrder, botID)

	go r.broadcastState()
}

func (r *Room) startGame(mode string, settings map[string]int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.Players) < 1 { // Allow solo play for testing?
		return
	}
	
	if mode != "" {
		r.Mode = mode
	} else {
		r.Mode = "CLASSIC"
	}
	
	if settings != nil {
		r.Settings = settings
	} else {
		r.Settings = make(map[string]int)
	}

	r.State = StatePlaying
	r.CurrentTurnIndex = 0
	r.UsedWords = make(map[string]bool)
	r.LastWord = ""
	r.History = []Move{}
	r.Round = 1
	r.TurnStartTime = time.Now()

	for _, p := range r.Players {
		p.Lives = 3
		if r.Mode == "SUDDEN_DEATH" {
			p.Lives = 1
		}
		p.IsTurn = false
		p.Score = 0
	}

	firstPlayerID := r.TurnOrder[0]
	r.Players[firstPlayerID].IsTurn = true

	r.broadcastStateInternal()

	if r.Players[firstPlayerID].Type == PlayerBot {
		go func() {
			time.Sleep(1 * time.Second)
			r.Action <- &ActionMessage{Type: "BOT_MOVE", PlayerID: firstPlayerID}
		}()
	}
}

func (r *Room) processBotTurn(playerID string) {
	r.mu.Lock()
	if r.State != StatePlaying {
		r.mu.Unlock()
		return
	}
	currentPlayerID := r.TurnOrder[r.CurrentTurnIndex]
	if currentPlayerID != playerID {
		r.mu.Unlock()
		return
	}

	lastWord := r.LastWord
	used := make(map[string]bool)
	for k, v := range r.UsedWords {
		used[k] = v
	}
	r.mu.Unlock()

	log.Printf("[Bot] Thinking for last word: %s", lastWord)
	move := r.BotBrain.GetMove(lastWord, used)
	log.Printf("[Bot] Decided: %s", move.Name)

	if move.Name == "" {
		// Bot gives up or failed
		r.mu.Lock()
		player := r.Players[playerID]
		player.Lives--
		log.Printf("[Bot] Failed/Gave up, lives left: %d", player.Lives)
		r.nextTurn()
		
		if r.checkGameOver() {
			r.broadcastStateInternal()
			r.mu.Unlock()
			return
		}

		r.broadcastStateInternal()

		// Trigger NEXT bot if applicable
		nextPlayerID := r.TurnOrder[r.CurrentTurnIndex]
		if r.Players[nextPlayerID].Type == PlayerBot {
			go func() {
				time.Sleep(2 * time.Second)
				r.Action <- &ActionMessage{Type: "BOT_MOVE", PlayerID: nextPlayerID}
			}()
		}
		r.mu.Unlock()
	} else {
		r.processTurn(playerID, move.Name)
	}
}

func (r *Room) processTurn(playerID string, word string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.State != StatePlaying {
		return
	}

	currentPlayerID := r.TurnOrder[r.CurrentTurnIndex]
	if currentPlayerID != playerID {
		return
	}

	word = strings.TrimSpace(word)
	lowerWord := strings.ToLower(word)
	player := r.Players[playerID]

	handleFailure := func(msg string) {
		player.Lives--
		if r.Mode == "SUDDEN_DEATH" {
			player.Lives = 0
		}
		r.sendError(playerID, msg)
		r.nextTurn()
		
		if r.checkGameOver() {
			r.broadcastStateInternal()
			return
		}

		r.broadcastStateInternal()

		// Trigger next bot if it's their turn now
		nextPlayerID := r.TurnOrder[r.CurrentTurnIndex]
		if r.Players[nextPlayerID].Type == PlayerBot {
			go func() {
				time.Sleep(2 * time.Second)
				r.Action <- &ActionMessage{Type: "BOT_MOVE", PlayerID: nextPlayerID}
			}()
		}
	}

	// Validate
	isValid, pType, canonicalName := r.Dict.IsValid(word)
	if !isValid {
		handleFailure("Invalid place name!")
		return
	}
	if r.UsedWords[lowerWord] {
		handleFailure("Place already used!")
		return
	}
	if r.LastWord != "" {
		lastChar := strings.ToLower(string(r.LastWord[len(r.LastWord)-1]))
		firstChar := strings.ToLower(string(word[0]))
		if lastChar != firstChar {
			handleFailure(fmt.Sprintf("Must start with '%s'!", strings.ToUpper(lastChar)))
			return
		}
	}

	// Record Move
	r.UsedWords[lowerWord] = true
	r.LastWord = canonicalName
	player.IsTurn = false
	player.MostUsedPlaces[lowerWord]++
	
	// Scoring Logic
	points := 10
	if r.Mode == "POINT_RUSH" {
		seconds := time.Since(r.TurnStartTime).Seconds()
		if seconds < 1 {
			seconds = 1
		}
		points = int(100.0/seconds) + (len(word) * 5)
	}
	player.Score += points

	r.History = append(r.History, Move{
		PlayerID:   playerID,
		PlayerName: player.Name,
		Word:       canonicalName,
		Type:       pType,
		Timestamp:  time.Now().Unix(),
	})

	r.nextTurn()
	r.broadcastStateInternal()

	if r.State == StateEnded {
		return
	}

	nextPlayerID := r.TurnOrder[r.CurrentTurnIndex]
	if r.Players[nextPlayerID].Type == PlayerBot {
		go func() {
			time.Sleep(2 * time.Second)
			r.Action <- &ActionMessage{Type: "BOT_MOVE", PlayerID: nextPlayerID}
		}()
	}
}

func (r *Room) nextTurn() {
	player := r.Players[r.TurnOrder[r.CurrentTurnIndex]]
	player.IsTurn = false
	r.TurnStartTime = time.Now()

	// Find next player with lives
	for i := 0; i < len(r.TurnOrder); i++ {
		r.CurrentTurnIndex = (r.CurrentTurnIndex + 1) % len(r.TurnOrder)
		if r.CurrentTurnIndex == 0 {
			r.Round++
		}
		nextPlayerID := r.TurnOrder[r.CurrentTurnIndex]
		if r.Players[nextPlayerID].Lives > 0 {
			r.Players[nextPlayerID].IsTurn = true
			return
		}
	}
	// If we get here, no one has lives
	r.State = StateEnded
	r.checkGameOver() // Ensure stats are saved
}

func (r *Room) checkGameOver() bool {
	alivePlayers := 0
	var winnerName string
	var winnerID string
	
	for _, p := range r.Players {
		if p.Lives > 0 {
			alivePlayers++
			winnerName = p.Name
			winnerID = p.ID
		}
	}

	isGameOver := false
	if alivePlayers <= 1 && len(r.Players) > 1 {
		r.State = StateEnded
		log.Printf("Game Over! Winner: %s", winnerName)
		isGameOver = true
	} else if alivePlayers == 0 && len(r.Players) == 1 {
		r.State = StateEnded
		log.Printf("Game Over! Solo player out of lives.")
		winnerID = "" // No winner
		isGameOver = true
	} else if len(r.Players) == 0 {
		r.State = StateEnded
		return true
	}

	if isGameOver && r.UserManager != nil {
		// Save Stats
		for _, p := range r.Players {
			if p.Type == PlayerHuman {
				isWin := (p.ID == winnerID)
				r.UserManager.UpdateStats(p.Name, p.Score, isWin)
			}
		}
	}
	
	return isGameOver
}

func (r *Room) broadcastState() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.broadcastStateInternal()
}

func (r *Room) broadcastStateInternal() {
	currentTurn := ""
	if len(r.TurnOrder) > 0 && r.CurrentTurnIndex < len(r.TurnOrder) {
		currentTurn = r.TurnOrder[r.CurrentTurnIndex]
	}

	state := map[string]interface{}{
		"type": "GAME_STATE",
		"payload": map[string]interface{}{
			"players":     r.Players,
			"state":       r.State,
			"lastWord":    r.LastWord,
			"turnOrder":   r.TurnOrder,
			"currentTurn": currentTurn,
			"history":     r.History,
			"round":       r.Round,
		},
	}
	bytes, _ := json.Marshal(state)

	for _, player := range r.Players {
		if player.Type == PlayerHuman {
			select {
			case player.Send <- bytes:
			default:
				// Avoid blocking if channel is full
			}
		}
	}
}

func (r *Room) sendError(playerID string, msg string) {
	p, ok := r.Players[playerID]
	if !ok || p.Type != PlayerHuman {
		return
	}

	errData := map[string]interface{}{
		"type":    "ERROR",
		"payload": map[string]string{"message": msg},
	}
	bytes, _ := json.Marshal(errData)
	p.Send <- bytes
}