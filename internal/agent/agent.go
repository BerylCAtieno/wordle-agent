package agent

import (
	"encoding/json"
	"fmt"
	"os"
)

// Global variable to hold the loaded Agent Card JSON content
var AgentCardData []byte

// loadAgentCard reads the agent_card.json file into the global variable.
func LoadAgentCard() error {
	const filename = "agent_card.json"
	var err error

	// Read the file content
	AgentCardData, err = os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read agent card file '%s': %w", filename, err)
	}

	// Optional: Validate the JSON structure before serving
	var temp map[string]interface{}
	if err := json.Unmarshal(AgentCardData, &temp); err != nil {
		return fmt.Errorf("agent card file '%s' contains invalid JSON: %w", filename, err)
	}

	return nil
}
