package main

/*
#include "../brickgame_api.h"
#include <stdlib.h>
*/
import "C"
import (
	"race/game"
	"runtime"
	"unsafe"
)

// Main game structure for binding with BrickGame API
var info C.GameInfo_t
var isAllocated bool = false
var firstInit bool = true

// Variables for allocating and storing memory for main game struct
var (
	fieldData     []C.int
	fieldPointers []*C.int
	nextData      []C.int
	nextPointers  []*C.int
	fieldPinner   runtime.Pinner
)

//export updateCurrentState
func updateCurrentState() C.GameInfo_t {
	if firstInit {
		initGameInfo()
		firstInit = false
	}

	game.FiniteStateMachine()

	info.score = C.int(game.Instance().Score)
	info.high_score = C.int(game.Instance().HighScore)
	info.level = C.int(game.Instance().Level)
	info.speed = C.int(game.Instance().Speed)
	info.pause = C.int(game.Instance().Pause)

	if game.Instance().State == game.EXIT {
		info.level = C.int(-1)
		freeFields()
		return info
	}

	for i := 0; i < game.HEIGHT; i++ {
		for j := 0; j < game.WIDTH; j++ {
			fieldData[i*game.WIDTH+j] = C.int(game.Instance().Field[i][j])
		}
	}

	applyCarPosition()

	for i := 0; i < game.NEXT_SIZE; i++ {
		for j := 0; j < game.NEXT_SIZE; j++ {
			nextData[i*game.NEXT_SIZE+j] = C.int(game.Instance().Next[i][j])
		}
	}

	return info
}

//export userInput
func userInput(action C.UserAction_t, hold C.bool) {
	game.Instance().UserInput = game.UserAction(action)
	if game.UserAction(action) == game.Terminate {
		game.Instance().State = game.EXIT
	}
}

// initGameInfo allocates memory for main game struct and fills it with zeroes
func initGameInfo() {
	if !isAllocated {
		allocateFields()
		isAllocated = true
	}
	setEmptyFields()
}

// allocateFields initializes memory for both arrays used by C code.
// It creates flat data arrays and corresponding pointer arrays for each row to
// simulate 2D arrays compatible with C's int** type. Runtime.Pinner is using to prevent
// garbage collector from moving it during operations in C code.
func allocateFields() {
	height := int(C.HEIGHT)
	width := int(C.WIDTH)
	totalFieldSize := height * width

	fieldData = make([]C.int, totalFieldSize)
	fieldPointers = make([]*C.int, height)
	fieldPinner.Pin(unsafe.SliceData(fieldData))

	for i := 0; i < height; i++ {
		fieldPointers[i] = &fieldData[i*width]
	}

	nextTotalSize := game.NEXT_SIZE * game.NEXT_SIZE
	nextData = make([]C.int, nextTotalSize)
	nextPointers = make([]*C.int, game.NEXT_SIZE)
	fieldPinner.Pin(unsafe.SliceData(nextData))

	for i := 0; i < game.NEXT_SIZE; i++ {
		nextPointers[i] = &nextData[i*game.NEXT_SIZE]
	}

	fieldPinner.Pin(unsafe.SliceData(fieldPointers))
	fieldPinner.Pin(unsafe.SliceData(nextPointers))

	info.field = (**C.int)(unsafe.Pointer(unsafe.SliceData(fieldPointers)))
	info.next = (**C.int)(unsafe.Pointer(unsafe.SliceData(nextPointers)))
}

// freeFields releases all allocated memory and unpins the arrays.
func freeFields() {
	if isAllocated {
		fieldPinner.Unpin()
		fieldData = nil
		nextData = nil
		fieldPointers = nil
		nextPointers = nil
		info.field = nil
		info.next = nil

		isAllocated = false
	}
}

// setEmptyFields is filling both fields with zeroes - Empty Cell
func setEmptyFields() {
	if !isAllocated {
		return
	}

	for i := 0; i < game.HEIGHT; i++ {
		for j := 0; j < game.WIDTH; j++ {
			fieldData[i*game.WIDTH+j] = C.int(game.Empty)
		}
	}

	for i := 0; i < game.NEXT_SIZE; i++ {
		for j := 0; j < game.NEXT_SIZE; j++ {
			nextData[i*game.NEXT_SIZE+j] = C.int(game.Empty)
		}
	}
}

// applyCarPosition is appends cells with player's car on game field
func applyCarPosition() {
	if !isAllocated {
		return
	}

	x := game.Instance().PlayerCar.X
	y := game.Instance().PlayerCar.Y

	for i := 0; i < game.CAR_HEIGHT; i++ {
		for j := 0; j < game.CAR_WIDTH; j++ {
			n := (y * game.WIDTH) + x
			cell := game.Instance().PlayerCar.Data[i][j]
			if cell != game.Empty {
				fieldData[n] = C.int(cell)
			}
			x++
		}
		x = game.Instance().PlayerCar.X
		y++
	}
}

func main() {}
