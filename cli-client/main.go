package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli-client <server-url> [player-name] [room-id]")
		fmt.Println("Example: cli-client ws://localhost:8080/ws Alice room123")
		os.Exit(1)
	}

	serverURL := os.Args[1]
	playerName := "Guest"
	roomID := ""

	if len(os.Args) > 2 {
		playerName = os.Args[2]
	}
	if len(os.Args) > 3 {
		roomID = os.Args[3]
	}

	// Build query params
	query := "?name=" + playerName
	if roomID != "" {
		query += "&room=" + roomID
	}

	fullURL := serverURL + query

	fmt.Printf("🎮 Connecting to %s as %s...\n", fullURL, playerName)

	ws, _, err := websocket.DefaultDialer.Dial(fullURL, nil)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer ws.Close()

	fmt.Println("✅ Connected!")

	// Read welcome message
	var welcome Message
	if err := ws.ReadJSON(&welcome); err != nil {
		log.Fatalf("Failed to read welcome: %v", err)
	}

	payload := welcome.Payload.(map[string]interface{})
	playerID := payload["id"].(string)
	joinedRoomID := payload["roomId"].(string)

	fmt.Printf("\n🎯 Player ID: %s\n", playerID)
	fmt.Printf("🏠 Room ID: %s\n\n", joinedRoomID)

	// Read responses in background
	go func() {
		for {
			var msg Message
			if err := ws.ReadJSON(&msg); err != nil {
				fmt.Printf("\n❌ Connection closed: %v\n", err)
				os.Exit(0)
			}

			fmt.Printf("\n📨 [%s]\n", msg.Type)
			if payloadBytes, err := json.MarshalIndent(msg.Payload, "", "  "); err == nil {
				fmt.Printf("%s\n", string(payloadBytes))
			}
			fmt.Print("\n> ")
		}
	}()

	// Interactive prompt
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Commands:")
	fmt.Println("  /join-room <room-id>    - Join a game room")
	fmt.Println("  /guess <answer>         - Make a guess")
	fmt.Println("  /start                  - Start game")
	fmt.Println("  /status                 - Get game status")
	fmt.Println("  /quit                   - Exit\n")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "/quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Parse command
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		var action Message

		switch cmd {
		case "/join-room":
			if len(parts) < 2 {
				fmt.Println("Usage: /join-room <room-id>")
				continue
			}
			action = Message{
				Type: "JOIN_ROOM",
				Payload: map[string]string{
					"roomId": parts[1],
				},
			}

		case "/guess":
			if len(parts) < 2 {
				fmt.Println("Usage: /guess <answer>")
				continue
			}
			action = Message{
				Type: "GUESS",
				Payload: map[string]string{
					"answer": parts[1],
				},
			}

		case "/start":
			action = Message{
				Type: "START_GAME",
				Payload: map[string]interface{}{},
			}

		case "/status":
			action = Message{
				Type: "GET_STATUS",
				Payload: map[string]interface{}{},
			}

		default:
			fmt.Println("Unknown command. Use /join-room, /guess, /start, /status, or /quit")
			continue
		}

		if err := ws.WriteJSON(action); err != nil {
			fmt.Printf("Send error: %v\n", err)
			break
		}
	}
}
