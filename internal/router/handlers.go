package router

import (
	"net/http"
	"strconv"

	"github.com/SaYaku64/monopoly/internal/alert"
	"github.com/gin-gonic/gin"
)

func (r *Router) CreateLobbyHandler(c *gin.Context) {
	fieldType := c.PostForm("fieldType")
	playerName := c.PostForm("playerName")
	sessionID := c.PostForm("sessionID")
	fastGame, err1 := strconv.ParseBool(c.PostForm("fastGame"))
	experimental, err2 := strconv.ParseBool(c.PostForm("experimental"))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data"})

		return
	}

	lobbyID := r.lm.CreateLobby(
		playerName,
		sessionID,
		fieldType,
		fastGame,
		experimental,
	)

	c.JSON(http.StatusOK, gin.H{"lobbyID": lobbyID})
}

func (r *Router) GetSessionID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sessionID": r.generateSessionToken()})
}

func (r *Router) GetLobbiesTable(c *gin.Context) {
	sessionID := c.Query("sessionID")

	alert.Info("GetLobbiesTable", sessionID)

	c.JSON(http.StatusOK, gin.H{"lobbiesTable": r.lm.GetLobbiesTableResponse(sessionID)})
}

func (r *Router) RemoveLobby(c *gin.Context) {
	lobbyID := c.Query("lobbyID")

	r.lm.RemoveLobby(lobbyID)

	c.Status(http.StatusOK)
}

func (r *Router) ConnectLobby(c *gin.Context) {
	sessionID := c.Query("sessionID")
	lobbyID := c.Query("lobbyID")

	alert.Info("ConnectLobby", sessionID, lobbyID)

	c.Status(http.StatusOK)
}
