package router

import (
	"net/http"
	"strconv"

	a "github.com/SaYaku64/business-game/internal/alert"
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

	a.Info.Println("GetLobbiesTable sessionID:", sessionID)

	c.JSON(http.StatusOK, gin.H{"lobbiesTable": r.lm.GetLobbiesTableResponse(sessionID)})
}

func (r *Router) RemoveLobby(c *gin.Context) {
	lobbyID := c.Query("lobbyID")

	r.lm.RemoveLobby(lobbyID)

	c.Status(http.StatusOK)
}

func (r *Router) ConnectLobby(c *gin.Context) {
	lobbyID := c.PostForm("lobbyID")
	playerName := c.PostForm("playerName")
	sessionID := c.PostForm("sessionID")

	lobby, err := r.lm.AddPlayerToLobby(lobbyID, playerName, sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if lobby != nil {
		r.gm.SetGame(*lobby)
	}

	c.Status(http.StatusOK)
}

func (r *Router) CheckActiveGame(c *gin.Context) {
	lobbyID := c.PostForm("lobbyID")
	playerName := c.PostForm("playerName")
	sessionID := c.PostForm("sessionID")

	active, plrTurn, current := r.gm.CheckActiveGame(lobbyID, playerName, sessionID)
	if active {
		c.JSON(http.StatusOK, gin.H{"turn": plrTurn, "current": current})
	} else {

		c.Status(http.StatusTeapot)
	}
}

func (r *Router) UpdatePlates(c *gin.Context) {
	lobbyID := c.PostForm("lobbyID")

	players, ok := r.gm.UpdatePlates(lobbyID)
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	obj := []gin.H{}

	for _, plr := range players {
		obj = append(obj, gin.H{
			"index":   plr.Index,
			"name":    plr.Name,
			"balance": plr.Balance,
		})
	}

	c.JSON(http.StatusOK, obj)
}

func (r *Router) IsLobbyExists(c *gin.Context) {
	lobbyID := c.PostForm("lobbyID")

	active := r.lm.IsLobbyExists(lobbyID)
	if active {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusTeapot)
	}
}
