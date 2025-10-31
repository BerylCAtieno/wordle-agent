package agent

import (
	_ "embed"
	"fmt"
)

var AgentCardData []byte

func LoadAgentCard() error {
	if len(AgentCardData) == 0 {
		return fmt.Errorf("agent card is empty")
	}
	return nil
}
