package router

import (
	"encoding/json"

	a "github.com/SaYaku64/business-game/internal/alert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GameLobby struct {
	ID      string   `json:"id"`      // lobbyID
	Players []string `json:"players"` // sessionIDs
	conns   []*websocket.Conn
}

func (r *Router) NextPlayerTurn(lobbyID, sessionID string) {
	r.gMux.RLock()
	gLobby, found := r.games[lobbyID]
	r.gMux.RUnlock()
	if !found {
		return
	}

	for i := range gLobby.conns {
		if gLobby.Players[i] == sessionID {
			gLobby.conns[i].WriteMessage(websocket.TextMessage, []byte("take turn"))
		}
	}
}

func (r *Router) SendMsgChat(lobbyID string, data gin.H) {
	r.gMux.RLock()
	gLobby, found := r.games[lobbyID]
	r.gMux.RUnlock()
	if !found {
		return
	}

	data["struct"] = true

	bMsg, _ := json.Marshal(data)

	for i := range gLobby.conns {
		gLobby.conns[i].WriteMessage(websocket.TextMessage, bMsg)

	}
}

func (r *Router) HandleWSGame(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		a.Error.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	lobbyID := c.Query("lobbyID")
	sessionID := c.Query("sessionID")

	r.gMux.Lock()
	if lobby, ok := r.games[lobbyID]; ok {
		a.Info.Println("games[lobbyID] ok")
		lobby.Players = append(lobby.Players, sessionID)
		lobby.conns = append(lobby.conns, conn)
		r.games[lobbyID] = lobby
		a.Info.Println("games[lobbyID] ok, lobby", *lobby)
	} else {
		a.Info.Println("games[lobbyID] !ok")
		lobby := &GameLobby{
			ID:      lobbyID,
			Players: []string{sessionID},
			conns:   []*websocket.Conn{conn},
		}
		r.games[lobbyID] = lobby
		a.Info.Println("games[lobbyID] ok, lobby", *lobby)
	}
	r.gMux.Unlock()

	r.readGameMsgs(conn, sessionID, lobbyID)
}

func (r *Router) readGameMsgs(
	conn *websocket.Conn,
	sessionID, lobbyID string,
) {
	for {
		_, byteMsg, err := conn.ReadMessage()
		if err != nil {
			r.deletePlayerFromGame(sessionID, lobbyID)

			a.Error.Printf("conn.ReadMessage error. lobbyID: %s; sessionID: %s; err: %s\n", lobbyID, sessionID, err.Error())
			break
		}

		a.Info.Printf("readGameMsgs. lobbyID: %s; sessionID: %s; byteMsg: %s\n", lobbyID, sessionID, string(byteMsg))
	}
}

func (r *Router) deletePlayerFromGame(sessionID, lobbyID string) {
	r.gMux.Lock()
	lobby := r.games[lobbyID]
	for i, player := range lobby.Players {
		if player == sessionID {
			lobby.Players = append(lobby.Players[:i], lobby.Players[i+1:]...)
			lobby.conns = append(lobby.conns[:i], lobby.conns[i+1:]...)
			break
		}
	}
	r.games[lobbyID] = lobby
	r.gMux.Unlock()
}
