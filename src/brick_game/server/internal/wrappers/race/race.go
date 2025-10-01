package race

import (
	"race/game"
	"server/internal/model"
	"server/internal/wrappers"
)

// RaceGame implements the BrickGame interface for the "Race" game.
type RaceGame struct{}

// NewRace creates a new RaceGame instance and returns it as a BrickGame.
func NewRace() wrappers.BrickGame {
	return &RaceGame{}
}

// UserInput processes a user action in the Race game.
func (r *RaceGame) UserInput(action int, hold bool) {
	game.Instance().UserInput = game.UserAction(action)
}

// UpdateCurrentState updates the game state and returns it as a model.GameState.
func (r *RaceGame) UpdateCurrentState() model.GameState {
	game.FiniteStateMachine()
	g := *game.Instance()
	game.ApplyCarPosition(&g)

	return r.convertState(g)
}

// convertState converts the internal game.Game struct to a model.GameState
// suitable for JSON serialization and frontend rendering.
func (r *RaceGame) convertState(g game.Game) model.GameState {
	field := make([][]int, game.HEIGHT)
	for i := 0; i < game.HEIGHT; i++ {
		field[i] = make([]int, game.WIDTH)
		for j := 0; j < game.WIDTH; j++ {
			field[i][j] = int(g.Field[i][j])
		}
	}

	next := make([][]int, game.NEXT_SIZE)
	for i := 0; i < game.NEXT_SIZE; i++ {
		next[i] = make([]int, game.NEXT_SIZE)
		for j := 0; j < game.NEXT_SIZE; j++ {
			next[i][j] = int(g.Next[i][j])
		}
	}

	return model.GameState{
		Field:     field,
		Next:      next,
		Score:     g.Score,
		HighScore: g.HighScore,
		Level:     g.Level,
		Speed:     int(g.Speed),
		Pause:     int(g.Pause),
	}
}
