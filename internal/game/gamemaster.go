package game

import (
	"fmt"
	"strings"

	"github.com/BerylCAtieno/wordle-agent/internal/dictionary"
	"github.com/BerylCAtieno/wordle-agent/internal/messages"
)

type GameMaster struct {
	Secret     string
	NameID     string
	Dictionary *dictionary.Dictionary
}

func NewGameMaster(secret string, dict *dictionary.Dictionary) *GameMaster {
	return &GameMaster{
		Secret:     strings.ToUpper(secret),
		NameID:     "GameMaster",
		Dictionary: dict,
	}
}

func (gm *GameMaster) Name() string { return gm.NameID }

func (gm *GameMaster) EvaluateGuess(msg messages.Message) messages.Message {
	guess := strings.ToUpper(msg.Content)

	// Check if the word is valid
	if !gm.Dictionary.IsValid(guess) {
		return messages.Message{
			From:    gm.Name(),
			To:      msg.From,
			Type:    messages.ErrorMessage,
			Content: "Not a valid word in the dictionary",
		}
	}

	feedback := gm.evaluateGuess(guess)
	return messages.Message{
		From:    gm.Name(),
		To:      msg.From,
		Type:    messages.FeedbackMessage,
		Content: feedback,
	}
}

// Evaluates a guess and returns 🟩🟨⬛ feedback
func (gm *GameMaster) evaluateGuess(guess string) string {
	result := make([]rune, len(gm.Secret))

	for i, ch := range guess {
		if ch == rune(gm.Secret[i]) {
			result[i] = '🟩'
		} else if strings.ContainsRune(gm.Secret, ch) {
			result[i] = '🟨'
		} else {
			result[i] = '⬛'
		}
	}
	return string(result)
}

func (gm *GameMaster) PrintIntro() {
	fmt.Println("=====================================")
	fmt.Println("🤖 Welcome to Agent Wordle!")
	fmt.Println("Try to guess the 5-letter secret word.")
	fmt.Println("🟩 = correct, 🟨 = wrong position, ⬛ = not in word.")
	fmt.Println("=====================================")
}
