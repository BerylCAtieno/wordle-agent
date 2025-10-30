package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BerylCAtieno/wordle-agent/internal/dictionary"
	"github.com/BerylCAtieno/wordle-agent/internal/messages"
)

func PlayGame() {
	dict, err := dictionary.LoadFromFile("internal/dictionary/words.txt")
	if err != nil {
		fmt.Printf("Error loading dictionary: %v\n", err)
		os.Exit(1)
	}

	randomWord := dict.RandomWord()

	gm := NewGameMaster(randomWord, dict)
	gm.PrintIntro()

	reader := bufio.NewReader(os.Stdin)
	attempts := 6

	for i := 1; i <= attempts; i++ {
		fmt.Printf("\nAttempt %d/%d — Enter your guess: ", i, attempts)
		input, _ := reader.ReadString('\n')
		guess := strings.TrimSpace(input)

		if len(guess) != len(gm.Secret) {
			fmt.Println("❌ Please enter a 5-letter word.")
			i--
			continue
		}

		msg := messages.Message{
			From:    "HumanPlayer",
			To:      gm.Name(),
			Type:    messages.GuessMessage,
			Content: guess,
		}

		// GameMaster evaluates
		feedback := gm.EvaluateGuess(msg)

		// Check if it's an error message (invalid word)
		if feedback.Type == messages.ErrorMessage {
			fmt.Printf("❌ %s\n", feedback.Content)
			i--
			continue
		}

		fmt.Printf("Feedback: %s\n", feedback.Content)

		if feedback.Content == "🟩🟩🟩🟩🟩" {
			fmt.Println("🎉 Congratulations! You guessed the word!")
			return
		}
	}

	fmt.Printf("\n😞 Game over! The word was %s.\n", gm.Secret)
}
