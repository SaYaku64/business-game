package router

import (
	a "github.com/SaYaku64/business-game/internal/alert"
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
		a.Error.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	lobbyID := c.Query("lobbyID")
	sessionID := c.Query("sessionID")

	if lobby, ok := lobbies[lobbyID]; ok {
		a.Info.Println("lobbies[lobbyID] ok")
		lobby.Players = append(lobby.Players, sessionID)
		lobbies[lobbyID] = lobby
		a.Info.Println("lobbies[lobbyID] ok, lobby", *lobby)
	} else {
		a.Info.Println("lobbies[lobbyID] !ok")
		lobby := &Lobby{
			ID:      lobbyID,
			Players: []string{sessionID},
		}
		lobbies[lobbyID] = lobby
		a.Info.Println("lobbies[lobbyID] ok, lobby", *lobby)
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

			a.Error.Printf("conn.ReadMessage error. lobbyID: %s; sessionID: %s; err: %s\n", lobbyID, sessionID, err.Error())
			break
		}

		a.Info.Printf("readGameMsgs. lobbyID: %s; sessionID: %s; byteMsg: %s\n", lobbyID, sessionID, string(byteMsg))
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
