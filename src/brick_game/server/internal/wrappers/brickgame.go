package wrappers

import "server/internal/model"

// BrickGame defines the common interface that all games must implement.
type BrickGame interface {
	// UserInput handles a user action in the game.
	UserInput(action int, hold bool)
	// UpdateCurrentState returns the current state of the game as a GameState struct.
	UpdateCurrentState() model.GameState
}
