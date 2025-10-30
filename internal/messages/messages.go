package messages

type MessageType string

const (
	GuessMessage    MessageType = "GUESS"
	FeedbackMessage MessageType = "FEEDBACK"
)

type Message struct {
	From    string
	To      string
	Type    MessageType
	Content string
}
