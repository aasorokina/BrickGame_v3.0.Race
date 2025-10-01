package game

import (
	"math/rand"
	"slices"
)

// FiniteStateMachine executes the main game logic based on the current state.
func FiniteStateMachine() {
	switch Instance().State {
	case START:
		StartActions()
	case SPAWN:
		GenerateNextCar()
	case MOVING:
		MovingActions()
	case SHIFTING:
		Shift()
	case COLLIDE:
		crash()
	case PAUSE:
		PauseActions()
	case GAMEOVER:
		gameoverActions()
	case EXIT:
		Instance().Level = -1
		Instance().Initialized = false
	default:
		break
	}
	Instance().UserInput = -1
}

// StartActions handles game initialization and menu interactions.
func StartActions() {
	switch Instance().UserInput {
	case Start:
		Instance().State = MOVING
		Instance().Pause = PlayScreen
	case Terminate:
		Instance().State = EXIT
	default:
		break
	}
}

// MovingActions processes player input, checks shifting timer and collision with other cars.
func MovingActions() {
	if getTime()-Instance().Time >= Instance().Speed {
		Instance().State = SHIFTING
		Instance().Time = getTime()
	}

	switch Instance().UserInput {
	case Left:
		moveTowards(Left)
	case Right:
		moveTowards(Right)
	case Up:
		Shift()
	case Pause:
		Instance().Pause = PauseScreen
		Instance().State = PAUSE
	case Terminate:
		Instance().State = EXIT
	default:
		break
	}

	if IsColliding() {
		Instance().State = COLLIDE
	}
}

// moveTowards adjusts the player car position based on directional input.
func moveTowards(act UserAction) {
	x := Instance().PlayerCar.X
	switch act {
	case Left:
		if x-TrackWidth >= LeftBorder {
			Instance().PlayerCar.X -= TrackWidth
		}
	case Right:
		if (x+TrackWidth)+(CAR_WIDTH-1) <= RightBorder {
			Instance().PlayerCar.X += TrackWidth
		}
	default:
		break
	}
}

// PauseActions manages game pause state and resume operations.
func PauseActions() {
	switch Instance().UserInput {
	case Pause:
		Instance().Pause = PlayScreen
		Instance().State = MOVING
		Instance().Time = getTime()
	case Terminate:
		Instance().State = EXIT
	default:
		break
	}
}

// GenerateNextCar creates next car at random route and distance
func GenerateNextCar() {
	prevRoute := Instance().NextCar.X
	route := rand.Intn(TrackCount)*TrackWidth + 1
	var newCarY int
	if route == prevRoute {
		newCarY = rand.Intn(5) + 1 + CAR_HEIGHT
	} else {
		newCarY = rand.Intn(2) + 6 + CAR_HEIGHT
	}
	Instance().NextCar = InitCar(RivalCarColor, route, -newCarY)
	Instance().State = MOVING
}

// Shift performs the game field movement by scrolling everything downward.
func Shift() {
	Instance().NextCar.Y++
	topLeft := topLeftBorder()

	for i := HEIGHT - 1; i >= 0; i-- {
		for j := 0; j < WIDTH; j++ {
			if i == 0 {
				Instance().Field[i][j] = Empty
			} else {
				Instance().Field[i][j] = Instance().Field[i-1][j]
			}
		}
	}

	if Instance().NextCar.Y > -CAR_HEIGHT {
		shiftNextCar()
	}

	Instance().Field[0][0] = topLeft

	if Instance().NextCar.IsEmpty() {
		Instance().State = SPAWN
	} else {
		Instance().State = MOVING
	}

	UpdateScore()
}

// topLeftBorder calculates top left cell for left road border basing on border pattern
func topLeftBorder() Cell {
	field := Instance().Field

	if field[0][0] != Empty && field[1][0] != Empty && field[2][0] != Empty {
		return Empty
	}

	return StripeColor
}

// shiftNextCar moves the next car on main game field when it becomes visible.
func shiftNextCar() {
	k := CAR_HEIGHT - 1
	x := Instance().NextCar.X
	for j := 0; j < CAR_WIDTH; j++ {
		Instance().Field[0][j+x] = Instance().NextCar.Data[k][j]
	}

	for i := CAR_HEIGHT - 1; i >= 0; i-- {
		for j := 0; j < CAR_WIDTH; j++ {
			if i == 0 {
				Instance().NextCar.Data[i][j] = Empty
			} else {
				Instance().NextCar.Data[i][j] = Instance().NextCar.Data[i-1][j]
			}
		}
	}

	if Instance().NextCar.Y >= 0 {
		SpawnedCoords = append(SpawnedCoords, Instance().NextCar.Y)
	}
}

// UpdateScore manages player scoring based on overtaken cars.
func UpdateScore() {
	for i := range SpawnedCoords {
		SpawnedCoords[i]++
	}

	for i := len(SpawnedCoords) - 1; i >= 0; i-- {
		if SpawnedCoords[i] > Instance().PlayerCar.Y {
			Instance().Score++
			SpawnedCoords = slices.Delete(SpawnedCoords, i, i+1)
		}
	}

	if Instance().Score > Instance().HighScore {
		Instance().HighScore = Instance().Score
		saveHighScore()
	}

	Instance().Level = min(Instance().Score/5+1, LEVEL_MAX)
	if Instance().Level != LEVEL_MIN {
		Instance().Speed = SPEED_MIN - int64(Instance().Level*65)
	}
}

// IsColliding detects whether the player car intersects with other cars
func IsColliding() bool {
	x := Instance().PlayerCar.X
	y := Instance().PlayerCar.Y

	for i := 0; i < CAR_HEIGHT; i++ {
		for j := 0; j < CAR_WIDTH; j++ {
			if Instance().Field[y][x] != Empty {
				return true
			}
			x++
		}
		x = Instance().PlayerCar.X
		y++
	}
	return false
}

// crash handles collision consequences by removing the player car and creating explosion effect.
func crash() {
	for i := 0; i < CAR_HEIGHT; i++ {
		for j := 0; j < CAR_WIDTH; j++ {
			Instance().PlayerCar.Data[i][j] = Empty
		}
	}

	x := Instance().PlayerCar.X
	y := Instance().PlayerCar.Y
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if x >= WIDTH || y >= HEIGHT {
				continue
			}
			Instance().Field[y][x] = Cell(burstTemplate[i][j] & int(PlayerColor))
			x++
		}
		x = Instance().PlayerCar.X
		y++
	}

	Instance().State = GAMEOVER
	Instance().Pause = GameOverScreen
}

// gameoverActions manages post-collision game state and restart options.
func gameoverActions() {
	switch Instance().UserInput {
	case Start:
		instance = NewGame()
		Instance().State = MOVING
		Instance().Pause = PlayScreen
	case Terminate:
		Instance().State = EXIT
	default:
		break
	}
}
