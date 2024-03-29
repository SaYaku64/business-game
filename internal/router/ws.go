package router

import (
	"encoding/json"

	a "github.com/SaYaku64/business-game/internal/alert"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var allWsReaders []*wsReader

// var (
// 	gameConn = NewGameWS()
// )

// func NewGameWS() *GameWS {
// 	return &GameWS{
// 		conn: make(map[string]*wsGameLobby),
// 	}
// }

// func (gw *GameWS) Find(lobbyID string) (*wsGameLobby, bool) {
// 	gw.mux.RLock()
// 	defer gw.mux.RUnlock()

// 	l, found := gw.conn[lobbyID]

// 	return l, found
// }

// func (gw *GameWS) Add(game *wsGameLobby) {
// 	gw.mux.Lock()
// 	gw.conn[game.lobbyID] = game
// 	gw.mux.Unlock()
// }

// type GameWS struct {
// 	conn map[string]*wsGameLobby // key - lobbyID
// 	mux  sync.RWMutex
// }

// type wsGameLobby struct {
// 	wsConn  *websocket.Conn
// 	lobbyID string

// 	sessions []string
// }

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
	a.Info.Println("socket request")
	if allWsReaders == nil {
		allWsReaders = make([]*wsReader, 0)
	}

	defer func() {
		err := recover()
		if err != nil {
			a.Error.Println(err)
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

// func (r *Router) HandleWSGame(c *gin.Context) {
// 	lobbyID := c.Query("lobbyID")
// 	sessionID := c.Query("sessionID")

// 	defer func() {
// 		err := recover()
// 		if err != nil {
// 			a.Error.Println(err)
// 		}
// 		c.Request.Body.Close()
// 	}()

// 	con, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

// 	lobby, exists := gameConn.Find(lobbyID)
// 	if exists {
// 		lobby.sessions = append(lobby.sessions, sessionID)

// 		return
// 	}

// 	newLobbyWs := &wsGameLobby{
// 		wsConn:  con,
// 		lobbyID: lobbyID,

// 		sessions: []string{sessionID},
// 	}
// 	gameConn.Add(newLobbyWs)

// 	newReader.startListening()
// }

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
		a.Error.Println("Unmarshal error", err, string(byteMsg))
		panic(err)
	}

	a.Info.Println("msg", msg)
	if msg.Action == actionCreate {
		i.lobbyID = msg.LobbyID
		i.sessionID1 = msg.SessionID
		a.Info.Println("waiting lobby for", msg.LobbyID)

		return
	}

	if msg.Action == actionConnect {
		i.sessionID2 = msg.SessionID
		a.Info.Println("connecting", msg.SessionID, msg.LobbyID)

		i.writeBoth(msg.LobbyID, msg.SessionID)

		return
	}

	if msg.Action == actionDelete {
		a.Info.Println("deleting", msg.LobbyID)

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

func (i *wsReader) writeBoth(lobbyID, sessionID string) {
	for index := range allWsReaders {
		if allWsReaders[index].lobbyID == lobbyID {
			i.lobbyID = lobbyID
			allWsReaders[index].sessionID2 = sessionID

			i.writeMsg("redirect")
			allWsReaders[index].writeMsg("redirect")

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
				a.Error.Println(err)
			}
			a.Info.Println("listening finished ", i.sessionID1)
		}()

		for {
			i.read()
		}

	}()
}
