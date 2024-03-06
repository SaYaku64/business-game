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

	alert.Info("sessionID", sessionID)
	c.Set("sessionID", sessionID)

	sessionIDT := c.GetString("sessionID")
	alert.Info("sessionIDT", sessionIDT)

	c.JSON(http.StatusOK, gin.H{"sessionID": sessionID})
}

func (r *Router) GetLobbiesTable(c *gin.Context) {
	sessionIDint, ok := c.Get("sessionID")
	alert.Info("GetString sessionIDint, ok", sessionIDint, ok)

	sessionID := c.GetString("sessionID")
	if sessionID == "" {
		sessionID = c.Query("sessionID")
		alert.Info("GetString sessionID from query")
	}

	c.JSON(http.StatusOK, gin.H{"lobbiesTable": r.lm.GetLobbiesTableResponse(sessionID)})
}

func (r *Router) RemoveLobby(c *gin.Context) {
	sessionID := c.GetString("sessionID")
	if sessionID == "" {
		sessionID = c.Query("sessionID")
		alert.Info("RemoveLobby sessionID from query")
	}

	r.lm.RemoveLobby(sessionID)

	c.Status(http.StatusOK)
}
