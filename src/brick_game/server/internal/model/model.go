package model

// GameInfo represents basic information about a game.
type GameInfo struct {
	ID   int    `json:"id"`   // Unique identifier of the game
	Name string `json:"name"` // Name of the game
}

// GamesList holds a list of available games.
type GamesList struct {
	Games []GameInfo `json:"games"` // Slice of GameInfo objects
}

// UserAction represents an action sent by the user to the server.
type UserAction struct {
	ActionID int  `json:"action_id"` // Identifier of the action performed by the user
	Hold     bool `json:"hold"`      // Whether the action is a hold/continuous action
}

// GameState represents the current state of a game, including the board and stats.
type GameState struct {
	Field     [][]int `json:"field"`      // 2D array representing the current game board
	Next      [][]int `json:"next"`       // 2D array representing the next piece/board
	Score     int     `json:"score"`      // Current score
	HighScore int     `json:"high_score"` // High score achieved
	Level     int     `json:"level"`      // Current level of the game
	Speed     int     `json:"speed"`      // Current speed of the game
	Pause     int     `json:"pause"`      // Pause state (0 = running, 1 = paused)
}

// ErrorMessage represents a JSON error message returned by the server.
type ErrorMessage struct {
	Message string `json:"message"` // Description of the error
}
