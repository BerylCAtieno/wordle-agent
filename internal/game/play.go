package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BerylCAtieno/wordle-agent/internal/dictionary"
)

func PlayGame() {
	dict, err := dictionary.LoadDictionary("internal/dictionary/words.txt")
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

		// GameMaster evaluates
		feedback, _ := gm.EvaluateGuess(guess)

		fmt.Printf("Feedback: %s\n", feedback)

		if feedback == "🟩🟩🟩🟩🟩" {
			fmt.Println("🎉 Congratulations! You guessed the word!")
			return
		}
	}

	fmt.Printf("\n😞 Game over! The word was %s.\n", gm.Secret)
}
