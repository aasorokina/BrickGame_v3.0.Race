package web

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"server/internal/model"
	"server/internal/wrappers"
	"server/internal/wrappers/race"
	"server/internal/wrappers/snake"
	"server/internal/wrappers/tetris"
)

var (
	games = []model.GameInfo{
		{ID: 1, Name: "Tetris"},
		{ID: 2, Name: "Snake"},
		{ID: 3, Name: "Race"},
	}

	mu            sync.Mutex
	currentGameID int
)

// GameHandler is responsible for handling game-related HTTP requests.
type GameHandler struct {
	currentGame wrappers.BrickGame
}

// NewGameHandler creates and returns a new GameHandler instance.
func NewGameHandler() *GameHandler {
	return &GameHandler{
		currentGame: nil,
	}
}

// GetGames returns the list of available games as JSON.
func (h *GameHandler) GetGames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, model.GamesList{Games: games})
}

// StartGame starts a new game based on the provided game ID.
func (h *GameHandler) StartGame(c *gin.Context) {
	idStr := c.Param("id")
	gameID, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.ErrorMessage{Message: "invalid game id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if h.currentGame != nil {
		c.IndentedJSON(http.StatusConflict, model.ErrorMessage{Message: "another game is already running"})
		return
	}

	switch gameID {
	case 1:
		h.currentGame = tetris.NewTetris()
	case 2:
		h.currentGame = snake.NewSnake()
	case 3:
		h.currentGame = race.NewRace()
	default:
		c.IndentedJSON(http.StatusNotFound, model.ErrorMessage{Message: "game not found"})
		return
	}

	c.Status(http.StatusOK)
}

// PostAction applies a user action to the current game.
func (h *GameHandler) PostAction(c *gin.Context) {
	var act model.UserAction
	if err := c.ShouldBindJSON(&act); err != nil {
		c.IndentedJSON(http.StatusBadRequest, model.ErrorMessage{Message: "invalid action"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if h.currentGame == nil {
		c.IndentedJSON(http.StatusBadRequest, model.ErrorMessage{Message: "no game started"})
		return
	}

	h.currentGame.UserInput(act.ActionID, act.Hold)

	if act.ActionID == 2 {
		h.currentGame.UpdateCurrentState()
		h.currentGame = nil
	}

	c.Status(http.StatusOK)
}

// GetState returns the current state of the running game as JSON.
func (h *GameHandler) GetState(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	if h.currentGame == nil {
		c.IndentedJSON(http.StatusBadRequest, model.ErrorMessage{Message: "no game started"})
		return
	}

	c.IndentedJSON(http.StatusOK, h.currentGame.UpdateCurrentState())
}
