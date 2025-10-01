package main

import (
	"testing"

	"race/game"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	instance *game.Game
}

func (suite *TestSuite) SetupTest() {
	suite.instance = game.Instance()
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestNewGame() {
	suite.instance = game.NewGame()
	assert.Equal(suite.T(), game.START, suite.instance.State)
	assert.Equal(suite.T(), 0, suite.instance.Score)
	assert.Equal(suite.T(), 1, suite.instance.Level)
	assert.Equal(suite.T(), int64(game.SPEED_MIN), suite.instance.Speed)
	assert.Equal(suite.T(), game.StartScreen, suite.instance.Pause)
	assert.False(suite.T(), suite.instance.PlayerCar.IsEmpty())
	assert.False(suite.T(), suite.instance.NextCar.IsEmpty())
}

func (suite *TestSuite) TestMove() {
	suite.instance.State = game.MOVING

	initialX := suite.instance.PlayerCar.X

	suite.instance.UserInput = game.Right
	game.MovingActions()
	assert.Equal(suite.T(), suite.instance.PlayerCar.X, initialX+game.TrackWidth)

	initialX = suite.instance.PlayerCar.X
	suite.instance.UserInput = game.Right
	game.MovingActions()
	assert.Equal(suite.T(), suite.instance.PlayerCar.X, initialX)

	suite.instance.UserInput = game.Left
	game.MovingActions()
	assert.Equal(suite.T(), initialX-game.TrackWidth, suite.instance.PlayerCar.X)
}

func (suite *TestSuite) TestCollision() {
	suite.instance.Field[suite.instance.PlayerCar.Y][suite.instance.PlayerCar.X] = game.RivalCarColor

	assert.True(suite.T(), game.IsColliding())

	suite.instance.Field[suite.instance.PlayerCar.Y][suite.instance.PlayerCar.X] = game.Empty
	assert.False(suite.T(), game.IsColliding())
}

func (suite *TestSuite) TestScoreUpdate() {
	suite.instance.Score = 0

	game.SpawnedCoords = []int{suite.instance.PlayerCar.Y, suite.instance.PlayerCar.Y - 1}
	game.UpdateScore()
	game.UpdateScore()

	assert.Equal(suite.T(), 2, suite.instance.Score)
	assert.Empty(suite.T(), game.SpawnedCoords)

	suite.instance.Score = 50
	suite.instance.HighScore = 10
	game.UpdateScore()
	assert.Equal(suite.T(), 50, suite.instance.HighScore)
}

func (suite *TestSuite) TestLevelUpdate() {
	suite.instance.Score = 0
	game.UpdateScore()
	assert.Equal(suite.T(), 1, suite.instance.Level)
	assert.Equal(suite.T(), int64(game.SPEED_MIN), suite.instance.Speed)

	suite.instance.Score = 10
	game.UpdateScore()
	assert.Equal(suite.T(), 3, suite.instance.Level)
	assert.Equal(suite.T(), game.SPEED_MIN-int64(3*65), suite.instance.Speed)
}

func (suite *TestSuite) TestStateTransitions() {
	suite.instance.State = game.START
	suite.instance.UserInput = game.Start
	game.StartActions()
	assert.Equal(suite.T(), game.MOVING, suite.instance.State)

	suite.instance.UserInput = game.Pause
	game.MovingActions()
	assert.Equal(suite.T(), game.PAUSE, suite.instance.State)

	suite.instance.UserInput = game.Pause
	game.PauseActions()
	assert.Equal(suite.T(), game.MOVING, suite.instance.State)

	suite.instance.UserInput = game.Start
	suite.instance.State = game.MOVING

	for i := 0; i < game.CAR_HEIGHT; i++ {
		for j := 0; j < game.CAR_WIDTH; j++ {
			suite.instance.NextCar.Data[i][j] = game.Empty
		}
	}
	suite.instance.NextCar.Y = 0
	game.Shift()
	assert.Equal(suite.T(), game.SPAWN, suite.instance.State)

	suite.instance.NextCar = game.InitCar(game.RivalCarColor, suite.instance.PlayerCar.X, -4)
	for i := 0; i <= game.HEIGHT-game.CAR_HEIGHT; i++ {
		game.Shift()
	}
	game.MovingActions()
	assert.Equal(suite.T(), game.COLLIDE, suite.instance.State)

	game.FiniteStateMachine()
	assert.Equal(suite.T(), game.GAMEOVER, suite.instance.State)

	suite.instance.UserInput = game.Start
	game.FiniteStateMachine()
	assert.Equal(suite.T(), game.MOVING, game.Instance().State)

	suite.instance = game.NewGame()
	userInput(_Ctype_UserAction_t(game.Start), false)
	game.FiniteStateMachine()
	assert.Equal(suite.T(), game.MOVING, game.Instance().State)
}

func (suite *TestSuite) TestCarGeneration() {
	suite.instance.State = game.SPAWN

	game.GenerateNextCar()

	assert.Equal(suite.T(), game.MOVING, suite.instance.State)
	assert.True(suite.T(), suite.instance.NextCar.Y < 0)
	assert.True(suite.T(), !suite.instance.NextCar.IsEmpty())
}

func (suite *TestSuite) TestAllocate() {
	initGameInfo()
	assert.True(suite.T(), isAllocated)
	for i := 0; i < game.HEIGHT; i++ {
		for j := 0; j < game.WIDTH; j++ {
			assert.True(suite.T(), fieldData[i*game.WIDTH+j] == 0)
		}
	}

	for i := 0; i < game.NEXT_SIZE; i++ {
		for j := 0; j < game.NEXT_SIZE; j++ {
			assert.True(suite.T(), nextData[i*game.NEXT_SIZE+j] == 0)
		}
	}

	suite.instance.Score = 10
	suite.instance.HighScore = 70
	suite.instance.Level = 3
	suite.instance.Speed = 700

	updateCurrentState()
	assert.Equal(suite.T(), suite.instance.Score, int(info.score))
	assert.Equal(suite.T(), suite.instance.HighScore, int(info.high_score))
	assert.Equal(suite.T(), suite.instance.Speed, int64(info.speed))
	assert.Equal(suite.T(), suite.instance.Level, int(info.level))

	freeFields()
	assert.True(suite.T(), !isAllocated)
	assert.True(suite.T(), nextPointers == nil)
	assert.True(suite.T(), fieldPointers == nil)
}
