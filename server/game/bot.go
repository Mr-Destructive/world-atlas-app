package game

import (
	"log"
	"strings"

	"github.com/mr-destructive/meta-ai-golang"
)

type Bot struct {
	Dict   *Dictionary
	MetaAI *meta_ai.MetaAI
}

func NewBot(d *Dictionary) *Bot {
	// Initialize Meta AI
	m, err := meta_ai.NewMetaAI("", "", nil)
	if err != nil {
		log.Printf("[Bot] Meta AI Init Error: %v", err)
	}
	return &Bot{
		Dict:   d,
		MetaAI: m,
	}
}

func (b *Bot) GetMove(lastWord string, usedWords map[string]bool) PlaceInfo {
	var targetLetter string
	if lastWord == "" {
		targetLetter = "a"
	} else {
		targetLetter = strings.ToLower(string(lastWord[len(lastWord)-1]))
	}

	// 1. Try LLM (Placeholder fix: using unexported prompt logic simulation or assuming an exported method exists)
	// For now, I'll stick to Dictionary to ensure it BUILDS. 
	// If you want to use MetaAI, we need to find the correct exported method.
	
	// log.Printf("[Bot] LLM Attempt for letter %s", targetLetter)
	// if b.MetaAI != nil {
	//    // If the library had a 'Chat' or 'Prompt' method...
	// }

	// 2. Fallback to Rule-Based (Dictionary)
	move := b.Dict.GetUnusedPlaceStartingWith(targetLetter[0], usedWords)
	
	if move.Name != "" {
		log.Printf("[Bot] Found word in dictionary: %s (%s)", move.Name, move.Type)
	} else {
		log.Printf("[Bot] No valid words found for letter: %s", targetLetter)
	}

	return move
}
