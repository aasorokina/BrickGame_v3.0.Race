package tetris

/*
#cgo LDFLAGS: ${SRCDIR}/../../../../../build/libs/tetris_wrap.a
#include "../../../../brickgame_api.h"

extern void tetris_userInput(UserAction_t action, bool hold);
extern GameInfo_t tetris_updateCurrentState();
*/
import "C"
import (
	"server/internal/model"
	"server/internal/wrappers"
	"unsafe"
)

// TetrisGame implements the BrickGame interface for the "Tetris" game.
type TetrisGame struct{}

// NewTetris creates a new TetrisGame instance and returns it as a BrickGame.
func NewTetris() wrappers.BrickGame {
	return &TetrisGame{}
}

// UserInput processes a user action in the Race game.
func (t *TetrisGame) UserInput(action int, hold bool) {
	C.tetris_userInput(C.UserAction_t(action), C.bool(hold))
}

// UpdateCurrentState updates the game state and returns it as a model.GameState.
func (t *TetrisGame) UpdateCurrentState() model.GameState {
	cstate := C.tetris_updateCurrentState()

	if cstate.field == nil || cstate.next == nil {
		return model.GameState{}
	}

	return t.convertState(C.GameInfo_t(cstate))
}

// convertState converts the internal game.Game struct to a model.GameState
// suitable for JSON serialization and frontend rendering.
func (t *TetrisGame) convertState(cstate C.GameInfo_t) model.GameState {
	field := make([][]int, C.HEIGHT)
	ptrField := (*[1 << 30]*C.int)(unsafe.Pointer(cstate.field))[:C.HEIGHT:C.HEIGHT]
	for i := 0; i < int(C.HEIGHT); i++ {
		rowData := (*[1 << 30]C.int)(unsafe.Pointer(ptrField[i]))[:C.WIDTH:C.WIDTH]
		row := make([]int, C.WIDTH)
		for j := 0; j < int(C.WIDTH); j++ {
			row[j] = int(rowData[j])
		}
		field[i] = row
	}

	next := make([][]int, C.NEXT_SIZE)
	ptrNext := (*[1 << 30]*C.int)(unsafe.Pointer(cstate.next))[:C.NEXT_SIZE:C.NEXT_SIZE]
	for i := 0; i < int(C.NEXT_SIZE); i++ {
		rowData := (*[1 << 30]C.int)(unsafe.Pointer(ptrNext[i]))[:C.NEXT_SIZE:C.NEXT_SIZE]
		row := make([]int, C.NEXT_SIZE)
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
