package router

import (
	"encoding/json"

	"github.com/SaYaku64/monopoly/internal/alert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var allWsReaders []*wsReader

// wsReader struct
type wsReader struct {
	wsConn     *websocket.Conn
	lobbyID    string
	sessionID1 string
	sessionID2 string
}

/*
wsMsg.sessionID = result.sessionID
wsMsg.create = true
*/
func removeFromSlice(s []*wsReader, i int) []*wsReader {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func deleteReader(lobbyID string) {
	for i := range allWsReaders {
		if allWsReaders[i].lobbyID == lobbyID {
			allWsReaders = removeFromSlice(allWsReaders, i)
		}
	}
}

func (r *Router) HandleWebSocket(c *gin.Context) {
	alert.Info("socket request")
	if allWsReaders == nil {
		allWsReaders = make([]*wsReader, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			alert.Error(err)
		}
		c.Request.Body.Close()
	}()

	con, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

	newReader := &wsReader{
		wsConn: con,
	}

	allWsReaders = append(allWsReaders, newReader)

	newReader.startListening()
}

// func (i *wsReader) broadcast(str string) {
// 	for _, g := range allWsReaders {

// 		if g == i {
// 			// no send message to himself
// 			continue
// 		}

// 		if g.mode == 1 {
// 			// no send message to connected user before user write his name
// 			continue
// 		}
// 		g.writeMsg(i.name, str)
// 	}
// }

type ReaderMsg struct {
	// SessionID1 string `json:"sessionID1"`
	// SessionID2 string `json:"sessionID2"`
	// Create     bool   `json:"create"`

	SessionID string `json:"sessionID"`
	Action    int    `json:"action"` // 1 - create; 2 - connect; 3 - delete
	LobbyID   string `json:"lobbyID"`
}

const (
	actionCreate = iota + 1
	actionConnect
	actionDelete
)

func (i *wsReader) read() {
	_, byteMsg, err := i.wsConn.ReadMessage()
	if err != nil {
		panic(err)
	}

	msg := ReaderMsg{}

	if err := json.Unmarshal(byteMsg, &msg); err != nil {
		alert.Error("Unmarshal error", err, string(byteMsg))
		panic(err)
	}

	alert.Info("msg", msg)
	if msg.Action == actionCreate {
		i.lobbyID = msg.LobbyID
		i.sessionID1 = msg.SessionID
		alert.Info("waiting lobby for", msg.LobbyID)

		return
	}

	if msg.Action == actionConnect {
		i.sessionID2 = msg.SessionID
		alert.Info("connecting", msg.SessionID, msg.LobbyID)

		i.writeBoth(msg.LobbyID)

		return
	}

	if msg.Action == actionDelete {
		alert.Info("deleting", msg.LobbyID)

		deleteReader(msg.LobbyID)
	}
	// log.Println(i.mode)

	// if i.mode == 1 {
	// 	i.name = string(byteMsg)
	// 	i.writeMsg("System", "Welcome "+i.name+", please write a message and we will broadcast it to other users.")
	// 	i.mode = 2 // real msg mode

	// 	return
	// }

	// i.broadcast(string(byteMsg))

	// log.Println(i.name + " " + string(byteMsg))
}

func (i *wsReader) writeBoth(lobbyID string) {
	for index := range allWsReaders {
		if allWsReaders[index].sessionID1 == lobbyID {
			i.writeMsg("redirect, please connected: " + lobbyID)
			allWsReaders[index].writeMsg("redirect, please creator: " + lobbyID)

			return
		}
	}
}

func (i *wsReader) writeMsg(str string) {
	i.wsConn.WriteMessage(websocket.TextMessage, []byte(str))
}

func (i *wsReader) startListening() {
	// i.writeMsg("startListening started")
	// i.mode = 1 //mode 1 get user name

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				alert.Error(err)
			}
			alert.Info("listening finished ", i.sessionID1)
		}()

		for {
			i.read()
		}

	}()
}

// func (r *Router) HandleWebSocket(c *gin.Context) {
// 	// Upgrade HTTP connection to websocket connection
// 	ws, err := websocket.Upgrader(c.Writer, c.Request, nil, 1024, 1024)
// 	if err != nil {
// 		http.Error(c.Writer, "Could not open websocket connection", http.StatusBadRequest)
// 		return
// 	}

// 	// Register client
// 	clients[ws] = true

// 	// Listen for messages from client
// 	go listenForMessages(ws)
// }
