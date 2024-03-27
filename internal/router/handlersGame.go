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
	if !ok {
		c.JSON(http.StatusBadRequest, answer)
		return
	}

	indexBefore := game.NextPlayerTurn()
	r.NextPlayerTurn(lobbyID, game.GetCurrentPlayer().SessionID, indexBefore, game.CurrentPlayer)
	// r.SendMsgChat(lobbyID, answer)
	r.SendUpdateField(lobbyID, answer)
	r.SendUpdateBalance(lobbyID, answer)

	c.JSON(http.StatusOK, answer)
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

	iB := game.NextPlayerTurn()
	r.NextPlayerTurn(lobbyID, game.GetCurrentPlayer().SessionID, iB, game.CurrentPlayer)
	r.SendMsgChat(lobbyID, answer)
	r.SendUpdateBalance(lobbyID, answer)

	c.JSON(status, answer)
}
