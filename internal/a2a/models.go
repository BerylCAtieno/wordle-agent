package a2a

import "time"

type MessageKind string
type MessageRole string
type PartKind string
type TaskState string

const (
	MessageKindMessage MessageKind = "message"

	RoleUser   MessageRole = "user"
	RoleAgent  MessageRole = "agent"
	RoleSystem MessageRole = "system"

	PartKindText PartKind = "text"
	PartKindData PartKind = "data"
	PartKindFile PartKind = "file"

	StateWorking       TaskState = "working"
	StateCompleted     TaskState = "completed"
	StateInputRequired TaskState = "input-required"
	StateFailed        TaskState = "failed"
)

// MessagePart represents a part of an A2A message
type MessagePart struct {
	Kind    PartKind               `json:"kind"`
	Text    *string                `json:"text,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	FileURL *string                `json:"file_url,omitempty"`
}

// A2AMessage represents a message in the A2A protocol
type A2AMessage struct {
	Kind      MessageKind            `json:"kind"`
	Role      MessageRole            `json:"role"`
	Parts     []MessagePart          `json:"parts"`
	MessageID string                 `json:"messageId"`
	TaskID    *string                `json:"taskId,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PushNotificationConfig for webhook notifications
type PushNotificationConfig struct {
	URL            string                 `json:"url"`
	Token          *string                `json:"token,omitempty"`
	Authentication map[string]interface{} `json:"authentication,omitempty"`
}

// MessageConfiguration for message handling
type MessageConfiguration struct {
	Blocking               bool                    `json:"blocking"`
	AcceptedOutputModes    []string                `json:"acceptedOutputModes"`
	PushNotificationConfig *PushNotificationConfig `json:"pushNotificationConfig,omitempty"`
}

// MessageParams for message/send method
type MessageParams struct {
	Message       A2AMessage           `json:"message"`
	Configuration MessageConfiguration `json:"configuration"`
}

// ExecuteParams for execute method
type ExecuteParams struct {
	ContextID *string      `json:"contextId,omitempty"`
	TaskID    *string      `json:"taskId,omitempty"`
	Messages  []A2AMessage `json:"messages"`
}

// JSONRPCRequest represents the incoming request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// TaskStatus represents the current state of a task
type TaskStatus struct {
	State     TaskState   `json:"state"`
	Timestamp string      `json:"timestamp"`
	Message   *A2AMessage `json:"message,omitempty"`
}

// Artifact represents a generated artifact
type Artifact struct {
	ArtifactID string        `json:"artifactId"`
	Name       string        `json:"name"`
	Parts      []MessagePart `json:"parts"`
}

// TaskResult represents the result of processing
type TaskResult struct {
	ID        string       `json:"id"`
	ContextID string       `json:"contextId"`
	Status    TaskStatus   `json:"status"`
	Artifacts []Artifact   `json:"artifacts"`
	History   []A2AMessage `json:"history"`
	Kind      string       `json:"kind"`
}

// JSONRPCResponse represents the response
type JSONRPCResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      string                 `json:"id"`
	Result  *TaskResult            `json:"result,omitempty"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

// Helper function to create timestamp
func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Helper function to create text message part
func TextPart(text string) MessagePart {
	return MessagePart{
		Kind: PartKindText,
		Text: &text,
	}
}

// Helper function to create data message part
func DataPart(data map[string]interface{}) MessagePart {
	return MessagePart{
		Kind: PartKindData,
		Data: data,
	}
}
