package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BerylCAtieno/wordle-agent/internal/a2a"
	"github.com/BerylCAtieno/wordle-agent/internal/agent"
	"github.com/BerylCAtieno/wordle-agent/internal/dictionary"
	"github.com/gin-gonic/gin"
)

func main() {

	if err := agent.LoadAgentCard(); err != nil {
		log.Fatal(err)
	}

	// 2. Initialize the Wordle dictionary and handler
	dict, err := dictionary.LoadDictionary("internal/dictionary/words.txt")
	if err != nil {
		log.Fatalf("Failed to load dictionary: %v", err)
	}

	wordleHandler := a2a.NewWordleHandler(dict)

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"agent":  "wordle",
		})
	})

	// --- A2A/Telex Integration Points ---

	// 1. Agent Card Discovery Endpoint (GET /.well-known/agent.json)
	router.GET("/.well-known/agent.json", func(c *gin.Context) {
		// Serve the pre-loaded JSON content
		c.Data(http.StatusOK, "application/json", agent.AgentCardData)
	})

	// 2. JSON-RPC Communication Endpoint (POST /)
	// All A2A methods (message/send, execute) are routed here.
	router.POST("/a2a/wordle", wordleHandler.HandleA2ARequest)

	// --- End A2A Integration Points ---

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001" // fallback for local dev
	}

	log.Printf("🎮 Wordle A2A Agent starting on 0.0.0.0:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(fmt.Errorf("server failed to start: %w", err))
	}
}
