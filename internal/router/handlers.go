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
	fastGame, err1 := strconv.ParseBool(c.PostForm("fastGame"))
	experimental, err2 := strconv.ParseBool(c.PostForm("experimental"))

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong data"})

		return
	}

	sessionID := r.lm.CreateLobby(
		playerName,
		fieldType,
		fastGame,
		experimental,
	)

	c.Set("sessionID", sessionID)

	c.JSON(http.StatusOK, gin.H{"sessionID": sessionID})
}

func (r *Router) GetLobbiesTable(c *gin.Context) {
	sessionID := c.GetString("sessionID")
	if sessionID == "" {
		sessionID = c.Query("sessionID")
	}

	c.JSON(http.StatusOK, gin.H{"lobbiesTable": r.lm.GetLobbiesTableResponse(sessionID)})
}

func (r *Router) RemoveLobby(c *gin.Context) {
	sessionID := c.GetString("sessionID")
	if sessionID == "" {
		sessionID = c.Query("sessionID")
	}

	r.lm.RemoveLobby(sessionID)

	c.Status(http.StatusOK)
}

func (r *Router) ConnectLobby(c *gin.Context) {
	sessionID := c.Query("sessionID")
	lobbyID := c.Query("lobbyID")

	alert.Info("ConnectLobby", sessionID, lobbyID)

	c.Status(http.StatusOK)
}
