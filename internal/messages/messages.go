package messages

type MessageType string

const (
	GuessMessage    MessageType = "GUESS"
	FeedbackMessage MessageType = "FEEDBACK"
	ErrorMessage    MessageType = "ERROR"
)

type Message struct {
	From    string
	To      string
	Type    MessageType
	Content string
}
