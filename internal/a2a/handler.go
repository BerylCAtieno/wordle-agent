package a2a

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/BerylCAtieno/wordle-agent/internal/dictionary"
	"github.com/BerylCAtieno/wordle-agent/internal/game"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GameSession stores the state of an active Wordle game
type GameSession struct {
	GameMaster  *game.GameMaster
	Attempts    int
	MaxAttempts int
	History     []A2AMessage
	IsComplete  bool
}

// WordleHandler handles A2A requests for Wordle
type WordleHandler struct {
	dictionary *dictionary.Dictionary
	sessions   map[string]*GameSession // contextID -> session
	mu         sync.RWMutex
}

func NewWordleHandler(dict *dictionary.Dictionary) *WordleHandler {
	return &WordleHandler{
		dictionary: dict,
		sessions:   make(map[string]*GameSession),
	}
}

// HandleA2ARequest processes incoming A2A requests
func (h *WordleHandler) HandleA2ARequest(c *gin.Context) {
	var req JSONRPCRequest

	if err := c.BindJSON(&req); err != nil {
		h.sendError(c, "", -32700, "Parse error", err.Error())
		return
	}

	// Validate JSON-RPC format
	if req.JSONRPC != "2.0" || req.ID == "" {
		h.sendError(c, req.ID, -32600, "Invalid Request", "jsonrpc must be '2.0' and id is required")
		return
	}

	// Process based on method
	switch req.Method {
	case "message/send":
		h.handleMessageSend(c, req)
	case "execute":
		h.handleExecute(c, req)
	default:
		h.sendError(c, req.ID, -32601, "Method not found", fmt.Sprintf("Unknown method: %s", req.Method))
	}
}

// handleMessageSend processes message/send requests
func (h *WordleHandler) handleMessageSend(c *gin.Context, req JSONRPCRequest) {
	// Parse params
	paramsJSON, _ := json.Marshal(req.Params)
	var params MessageParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		h.sendError(c, req.ID, -32602, "Invalid params", err.Error())
		return
	}

	// Process with a single message
	result, err := h.processMessages(
		[]A2AMessage{params.Message},
		nil,
		params.Message.TaskID,
		&params.Configuration,
	)

	if err != nil {
		h.sendError(c, req.ID, -32603, "Internal error", err.Error())
		return
	}

	h.sendResult(c, req.ID, result)
}

// handleExecute processes execute requests
func (h *WordleHandler) handleExecute(c *gin.Context, req JSONRPCRequest) {
	// Parse params
	paramsJSON, _ := json.Marshal(req.Params)
	var params ExecuteParams
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		h.sendError(c, req.ID, -32602, "Invalid params", err.Error())
		return
	}

	result, err := h.processMessages(
		params.Messages,
		params.ContextID,
		params.TaskID,
		nil,
	)

	if err != nil {
		h.sendError(c, req.ID, -32603, "Internal error", err.Error())
		return
	}

	h.sendResult(c, req.ID, result)
}

// processMessages handles the core game logic
func (h *WordleHandler) processMessages(
	messages []A2AMessage,
	contextID *string,
	taskID *string,
	config *MessageConfiguration,
) (*TaskResult, error) {
	// Generate IDs if not provided
	ctx := h.getOrCreateContextID(contextID)
	task := h.getOrCreateTaskID(taskID)

	// Get or create game session
	session := h.getOrCreateSession(ctx)

	// Extract user message
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	userMessage := messages[len(messages)-1]
	session.History = append(session.History, userMessage)

	// Extract guess from message
	guess := h.extractGuess(userMessage)

	// Default feedback for initialization/errors
	currentFeedback := ""

	// Check if game is already complete
	if session.IsComplete {
		responseText := fmt.Sprintf("The game is complete! The word was **%s**. Type 'new game' to play again.", session.GameMaster.Secret)
		return h.buildResult(
			task,
			ctx,
			session,
			responseText,
			StateCompleted,
			currentFeedback,
		), nil
	}

	// Handle special commands
	if strings.ToLower(guess) == "new game" || strings.ToLower(guess) == "restart" {
		h.resetSession(ctx)
		session = h.getOrCreateSession(ctx)
		responseText := fmt.Sprintf("🎮 New game started! I'm thinking of a 5-letter word. You have %d attempts. Make your first guess!", session.MaxAttempts)
		return h.buildResult(
			task,
			ctx,
			session,
			responseText,
			StateInputRequired,
			currentFeedback,
		), nil
	}

	// Validate guess length
	if len(guess) != 5 {
		responseText := "❌ Please enter exactly a 5-letter word."
		return h.buildResult(
			task,
			ctx,
			session,
			responseText,
			StateInputRequired,
			currentFeedback,
		), nil
	}

	// Process the guess
	// NOTE: GameMaster.EvaluateGuess is assumed to be updated to accept a string
	feedback, isValid := session.GameMaster.EvaluateGuess(guess)
	currentFeedback = feedback

	// Validate word is in dictionary
	if !isValid {
		responseText := "❌ Not a valid word in the dictionary. Try another word!"
		return h.buildResult(
			task,
			ctx,
			session,
			responseText,
			StateInputRequired,
			currentFeedback, // Empty feedback is fine here
		), nil
	}

	// Only increment attempts for valid words
	session.Attempts++

	// Check win condition
	if feedback == "🟩🟩🟩🟩🟩" {
		session.IsComplete = true
		responseText := fmt.Sprintf("🎉 Congratulations! You guessed the word in %d attempt(s)!\n\nFeedback: %s\n\nType 'new game' to play again!",
			session.Attempts, feedback)
		return h.buildResult(task, ctx, session, responseText, StateCompleted, currentFeedback), nil
	}

	// Check lose condition
	if session.Attempts >= session.MaxAttempts {
		session.IsComplete = true
		responseText := fmt.Sprintf("😞 Game over! You've used all %d attempts.\n\nThe word was: **%s**\n\nType 'new game' to play again!",
			session.MaxAttempts, session.GameMaster.Secret)
		return h.buildResult(task, ctx, session, responseText, StateCompleted, currentFeedback), nil
	}

	// Continue game
	responseText := fmt.Sprintf("Attempt %d/%d\n\nFeedback: %s\n\n🟩 = correct position\n🟨 = wrong position\n⬛ = not in word\n\nMake your next guess!",
		session.Attempts, session.MaxAttempts, feedback)
	return h.buildResult(task, ctx, session, responseText, StateInputRequired, currentFeedback), nil
}

// Helper methods

func (h *WordleHandler) getOrCreateContextID(id *string) string {
	if id != nil && *id != "" {
		return *id
	}
	return uuid.New().String()
}

func (h *WordleHandler) getOrCreateTaskID(id *string) string {
	if id != nil && *id != "" {
		return *id
	}
	return uuid.New().String()
}

func (h *WordleHandler) getOrCreateSession(contextID string) *GameSession {
	h.mu.Lock()
	defer h.mu.Unlock()

	if session, exists := h.sessions[contextID]; exists {
		return session
	}

	// Create new session with random word
	secret := h.dictionary.RandomWord()
	gm := game.NewGameMaster(secret, h.dictionary)

	session := &GameSession{
		GameMaster:  gm,
		Attempts:    0,
		MaxAttempts: 6,
		History:     []A2AMessage{},
		IsComplete:  false,
	}

	h.sessions[contextID] = session
	return session
}

func (h *WordleHandler) resetSession(contextID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.sessions, contextID)
}

func (h *WordleHandler) extractGuess(msg A2AMessage) string {
	for _, part := range msg.Parts {
		if part.Kind == PartKindText && part.Text != nil {
			return strings.TrimSpace(strings.ToUpper(*part.Text))
		}
	}
	return ""
}

// buildResult now accepts the 'feedback' string to be included as a dedicated artifact.
func (h *WordleHandler) buildResult(
	taskID string,
	contextID string,
	session *GameSession,
	responseText string,
	state TaskState,
	feedback string,
) *TaskResult {
	// Create response message
	responseMsg := A2AMessage{
		Kind:      MessageKindMessage,
		Role:      RoleAgent,
		Parts:     []MessagePart{TextPart(responseText)},
		MessageID: uuid.New().String(),
		TaskID:    &taskID,
	}

	session.History = append(session.History, responseMsg)

	// Create artifacts
	artifacts := []Artifact{
		// 1. Game State Artifact
		{
			ArtifactID: uuid.New().String(),
			Name:       "game_state",
			Parts: []MessagePart{
				DataPart(map[string]interface{}{
					"attempts":     session.Attempts,
					"max_attempts": session.MaxAttempts,
					"is_complete":  session.IsComplete,
				}),
			},
		},
		// 2. Wordle Feedback Artifact (Essential for client rendering)
		{
			ArtifactID: uuid.New().String(),
			Name:       "wordle_feedback",
			Parts: []MessagePart{
				DataPart(map[string]interface{}{
					"feedback_emojis": feedback, // e.g., "🟩🟨⬛🟩⬛"
				}),
			},
		},
	}

	return &TaskResult{
		ID:        taskID,
		ContextID: contextID,
		Status: TaskStatus{
			State:     state,
			Timestamp: Timestamp(),
			Message:   &responseMsg,
		},
		Artifacts: artifacts,
		History:   session.History,
		Kind:      "task",
	}
}

func (h *WordleHandler) sendResult(c *gin.Context, requestID string, result *TaskResult) {
	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      requestID,
		Result:  result,
	})
}

func (h *WordleHandler) sendError(c *gin.Context, requestID string, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      requestID,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	})
}
