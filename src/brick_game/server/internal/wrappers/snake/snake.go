package snake

/*
#cgo LDFLAGS: ${SRCDIR}/../../../../../build/libs/snake_wrap.a -lstdc++
#include "../../../../brickgame_api.h"

extern void snake_userInput(UserAction_t action, bool hold);
extern GameInfo_t snake_updateCurrentState();
*/
import "C"
import (
	"server/internal/model"
	"server/internal/wrappers"
	"unsafe"
)

// SnakeGame implements the BrickGame interface for the "Snake" game.
type SnakeGame struct{}

// NewSnake creates a new SnakeGame instance and returns it as a BrickGame.
func NewSnake() wrappers.BrickGame {
	return &SnakeGame{}
}

// UserInput processes a user action in the Race game.
func (s *SnakeGame) UserInput(action int, hold bool) {
	C.snake_userInput(C.UserAction_t(action), C.bool(hold))
}

// UpdateCurrentState updates the game state and returns it as a model.GameState.
func (s *SnakeGame) UpdateCurrentState() model.GameState {
	cstate := C.snake_updateCurrentState()
	if cstate.field == nil || cstate.next == nil {
		return model.GameState{}
	}
	return s.convertState(C.GameInfo_t(cstate))
}

// convertState converts the internal game.Game struct to a model.GameState
// suitable for JSON serialization and frontend rendering.
func (s *SnakeGame) convertState(cstate C.GameInfo_t) model.GameState {
	field := make([][]int, C.HEIGHT)
	ptrField := (*[1 << 30]*C.int)(unsafe.Pointer(cstate.field))[:C.HEIGHT:C.HEIGHT]
	for i := 0; i < int(C.HEIGHT); i++ {
		row := make([]int, C.WIDTH)
		rowData := (*[1 << 30]C.int)(unsafe.Pointer(ptrField[i]))[:C.WIDTH:C.WIDTH]
		for j := 0; j < int(C.WIDTH); j++ {
			row[j] = int(rowData[j])
		}
		field[i] = row
	}

	next := make([][]int, C.NEXT_SIZE)
	ptrNext := (*[1 << 30]*C.int)(unsafe.Pointer(cstate.next))[:C.NEXT_SIZE:C.NEXT_SIZE]
	for i := 0; i < int(C.NEXT_SIZE); i++ {
		row := make([]int, C.NEXT_SIZE)
		rowData := (*[1 << 30]C.int)(unsafe.Pointer(ptrNext[i]))[:C.NEXT_SIZE:C.NEXT_SIZE]
		for j := 0; j < int(C.NEXT_SIZE); j++ {
			row[j] = int(rowData[j])
		}
		next[i] = row
	}

	return model.GameState{
		Field:     field,
		Next:      next,
		Score:     int(cstate.score),
		HighScore: int(cstate.high_score),
		Level:     int(cstate.level),
		Speed:     int(cstate.speed),
		Pause:     int(cstate.pause),
	}
}
