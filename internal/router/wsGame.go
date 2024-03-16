package router

import (
	"fmt"

	"github.com/SaYaku64/business-game/internal/alert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Lobby struct {
	ID      string   `json:"id"`
	Players []string `json:"players"`
}

var lobbies = make(map[string]*Lobby)

func (r *Router) HandleWSGame(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		alert.Error("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	lobbyID := c.Query("lobbyID")
	sessionID := c.Query("sessionID")

	if lobby, ok := lobbies[lobbyID]; ok {
		alert.Info("lobbies[lobbyID] ok")
		lobby.Players = append(lobby.Players, sessionID)
		lobbies[lobbyID] = lobby
		alert.Info("lobbies[lobbyID] ok, lobby", *lobby)
	} else {
		alert.Info("lobbies[lobbyID] !ok")
		lobby := &Lobby{
			ID:      lobbyID,
			Players: []string{sessionID},
		}
		lobbies[lobbyID] = lobby
		alert.Info("lobbies[lobbyID] ok, lobby", *lobby)
	}

	readGameMsgs(conn, sessionID, lobbyID)
}

func readGameMsgs(
	conn *websocket.Conn,
	sessionID, lobbyID string,
) {
	for {
		_, byteMsg, err := conn.ReadMessage()
		if err != nil {
			deletePlayerFromLobby(sessionID, lobbyID)

			str := fmt.Sprintf("conn.ReadMessage error. lobbyID: %s; sessionID: %s; err: %s", lobbyID, sessionID, err.Error())
			alert.Error(str)
			break
		}

		str := fmt.Sprintf("readGameMsgs. lobbyID: %s; sessionID: %s; byteMsg: %s", lobbyID, sessionID, string(byteMsg))

		alert.Info(str)
	}
}

func deletePlayerFromLobby(sessionID, lobbyID string) {
	lobby := lobbies[lobbyID]
	for i, player := range lobby.Players {
		if player == sessionID {
			lobby.Players = append(lobby.Players[:i], lobby.Players[i+1:]...)
			break
		}
	}
	lobbies[lobbyID] = lobby
}
