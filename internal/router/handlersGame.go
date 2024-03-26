package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) Turn(c *gin.Context) {
	lobbyID := c.GetHeader("lobbyID")
	sessionID := c.GetHeader("sessionID")

	game, found := r.gm.GetGame(lobbyID)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
	}

	plr := game.GetCurrentPlayer()
	if plr.SessionID != sessionID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not players turn"})
		return
	}

	result := game.RollAction(plr)

	answer := game.CalculateRollActionResult(plr, result)
	r.SendMsgChat(lobbyID, answer)

	c.JSON(http.StatusOK, answer)
}

func (r *Router) Buy(c *gin.Context) {
	lobbyID := c.GetHeader("lobbyID")
	sessionID := c.GetHeader("sessionID")

	game, found := r.gm.GetGame(lobbyID)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
	}

	plr := game.GetCurrentPlayer()
	if plr.SessionID != sessionID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not players turn"})
		return
	}

	answer, ok := game.Buy(plr)
	status := http.StatusBadRequest
	if ok {
		status = http.StatusOK
		game.NextPlayerTurn()
		r.NextPlayerTurn(lobbyID, game.GetCurrentPlayer().SessionID)
		r.SendMsgChat(lobbyID, answer)
	}

	c.JSON(status, answer)
}

func (r *Router) PayRent(c *gin.Context) {
	lobbyID := c.GetHeader("lobbyID")
	sessionID := c.GetHeader("sessionID")

	game, found := r.gm.GetGame(lobbyID)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "game not found"})
		return
	}

	plr := game.GetCurrentPlayer()
	if plr.SessionID != sessionID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not players turn"})
		return
	}

	answer, ok := game.PayRent(plr)
	status := http.StatusOK
	if !ok {
		status = http.StatusBadRequest
	}

	game.NextPlayerTurn()
	r.NextPlayerTurn(lobbyID, game.GetCurrentPlayer().SessionID)
	r.SendMsgChat(lobbyID, answer)

	c.JSON(status, answer)
}
