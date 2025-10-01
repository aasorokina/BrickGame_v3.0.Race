package game

/*
#include "../../brickgame_api.h"
*/
import "C"
import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// instance is the global singleton instance of the Game structure.
// It is initialized only once using the sync.Once
var instance *Game
var once sync.Once

// SpawnedCoords is slice which contains Y coords of already spawned cars for score counting
var SpawnedCoords []int

// Constants from BrickGame API
const (
	HEIGHT     = C.HEIGHT
	WIDTH      = C.WIDTH
	NEXT_SIZE  = C.NEXT_SIZE
	LEVEL_MAX  = C.LEVEL_MAX
	LEVEL_MIN  = C.LEVEL_MIN
	SPEED_MIN  = 700
	CAR_HEIGHT = 4
	CAR_WIDTH  = 3
)

// UserAction is enum from BrickGame API
type UserAction int

// All UserAction variants
const (
	Start     UserAction = UserAction(C.Start)
	Pause     UserAction = UserAction(C.Pause)
	Terminate UserAction = UserAction(C.Terminate)
	Left      UserAction = UserAction(C.Left)
	Right     UserAction = UserAction(C.Right)
	Up        UserAction = UserAction(C.Up)
	Down      UserAction = UserAction(C.Down)
	Action    UserAction = UserAction(C.Action)
)

// ScreenState is using for rendering user interface depending on the game state
type ScreenState int

// All possible screen states
const (
	StartScreen    ScreenState = ScreenState(C.START_)
	PlayScreen     ScreenState = ScreenState(C.PLAYING_)
	PauseScreen    ScreenState = ScreenState(C.PAUSE_)
	GameOverScreen ScreenState = ScreenState(C.GAMEOVER_)
	WinScreen      ScreenState = ScreenState(C.WIN_)
)

// file for storing high score
const highscoreFile = "highscore_race.txt"

// GameState represents the current state of the game finite state machine
type GameState int

// Possible game states
const (
	START GameState = iota
	SPAWN
	MOVING
	SHIFTING
	COLLIDE
	GAMEOVER
	PAUSE
	EXIT
)

// Cell represents the color of a single grid cell in the game field.
type Cell int

// Possible cell colors
const (
	Empty      Cell = 0
	Red        Cell = 1
	Green      Cell = 2
	Blue       Cell = 4
	Cyan       Cell = 6
	Orange     Cell = 8
	Yellow     Cell = 9
	Violet     Cell = 10
	LightGreen Cell = 11
)

// Game element colors define the visual appearance of different game objects
const (
	PlayerColor   = Cyan
	RivalCarColor = Blue
	StripeColor   = Green
)

// Track configuration constants define the road structure and layout
const (
	TrackCount int = 3
	TrackWidth int = 3
)

// Car represents a game entity with visual data and position.
// The Data array stores the visual representation of the car in a grid pattern,
// X and Y is current position on the game field.
type Car struct {
	Data [CAR_HEIGHT][CAR_WIDTH]Cell
	X    int
	Y    int
}

// IsEmpty checks if the car has no visible cells in its data pattern.
func (c Car) IsEmpty() bool {
	for i := 0; i < CAR_HEIGHT; i++ {
		for j := 0; j < CAR_WIDTH; j++ {
			if c.Data[i][j] != Empty {
				return false
			}
		}
	}
	return true
}

// carTemplate defines the default visual pattern for cars.
var carTemplate = [CAR_HEIGHT][CAR_WIDTH]int{
	{0, 15, 0},
	{15, 15, 15},
	{0, 15, 0},
	{15, 0, 15},
}

// burstTemplate defines the default visual pattern for car crash.
var burstTemplate = [4][4]int{
	{15, 0, 0, 15},
	{0, 15, 15, 0},
	{0, 15, 15, 0},
	{15, 0, 0, 15},
}

// Constants for road borders
const (
	LeftBorder   = 1
	TopBorder    = 0
	RightBorder  = WIDTH - 1
	BottomBorder = HEIGHT - 1
)

// Game represents the complete state of the racing game instance.
type Game struct {
	Field       [HEIGHT][WIDTH]Cell
	Next        [NEXT_SIZE][NEXT_SIZE]Cell
	PlayerCar   Car
	NextCar     Car
	Time        int64
	UserInput   UserAction
	State       GameState
	Score       int
	HighScore   int
	Level       int
	Speed       int64
	Pause       ScreenState
	Initialized bool
}

// Instance returns singleton instance of the Game structure
func Instance() *Game {
	once.Do(func() {
		instance = NewGame()
	})

	if instance.Initialized == false {
		instance = NewGame()
	}

	return instance
}

// getTime returns Unix time in milliseconds
func getTime() int64 {
	return time.Now().UnixMilli()
}

// loadHighScore is getting high score from file to Game structure
func loadHighScore() int {
	file, err := os.Open(highscoreFile)
	if err != nil {
		return 0
	}

	scanner := bufio.NewScanner(file)
	score := 0
	for scanner.Scan() {
		score, err = strconv.Atoi(scanner.Text())
		if err != nil {
			return 0
		}
	}
	file.Close()
	return score
}

// loadHighScore is saving high score to file
func saveHighScore() {
	file, err := os.OpenFile(highscoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	file.WriteString(strconv.Itoa(Instance().HighScore))
	file.Close()
}

// NewGame is creates a new Game structure with default values
func NewGame() *Game {
	game := &Game{
		Time:        getTime(),
		UserInput:   -1,
		State:       START,
		Score:       0,
		Level:       1,
		Speed:       SPEED_MIN,
		Pause:       StartScreen,
		Initialized: true,
	}

	route := rand.Intn(TrackCount)*TrackWidth + 1
	game.NextCar = InitCar(RivalCarColor, route, -CAR_HEIGHT)
	game.PlayerCar = InitCar(PlayerColor, LeftBorder+TrackWidth, BottomBorder-CAR_HEIGHT+1)
	game.HighScore = loadHighScore()

	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			if i%4 != 0 && j == 0 {
				game.Field[i][j] = StripeColor
			} else {
				game.Field[i][j] = Empty
			}

		}
	}

	for i := 0; i < NEXT_SIZE; i++ {
		for j := 0; j < NEXT_SIZE; j++ {
			if (i+j)%2 == 0 && i < 2 {
				game.Next[i][j] = LightGreen
			} else {
				game.Next[i][j] = Empty
			}
		}
	}

	return game
}

// InitCar creates a car entity with the specified color and position.
func InitCar(cell Cell, x, y int) Car {
	var car Car
	for i := 0; i < CAR_HEIGHT; i++ {
		for j := 0; j < CAR_WIDTH; j++ {
			car.Data[i][j] = Cell(carTemplate[i][j] & int(cell))
		}
	}

	car.X = x
	car.Y = y

	return car
}

// ApplyCarPosition is appends cells with player's car on game field
func ApplyCarPosition(g *Game) {
	x := g.PlayerCar.X
	y := g.PlayerCar.Y

	for i := 0; i < CAR_HEIGHT; i++ {
		for j := 0; j < CAR_WIDTH; j++ {
			cell := g.PlayerCar.Data[i][j]
			if cell != Empty {
				g.Field[y][x] = cell
			}
			x++
		}
		x = g.PlayerCar.X
		y++
	}
}
