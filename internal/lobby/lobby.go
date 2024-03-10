package lobby

import (
	"fmt"
	"sync"
	"time"

	"github.com/SaYaku64/monopoly/internal/alert"
)

const (
	checkSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-check2-square" viewBox="0 0 16 16">
	<path d="M3 14.5A1.5 1.5 0 0 1 1.5 13V3A1.5 1.5 0 0 1 3 1.5h8a.5.5 0 0 1 0 1H3a.5.5 0 0 0-.5.5v10a.5.5 0 0 0 .5.5h10a.5.5 0 0 0 .5-.5V8a.5.5 0 0 1 1 0v5a1.5 1.5 0 0 1-1.5 1.5z"/>
	<path d="m8.354 10.354 7-7a.5.5 0 0 0-.708-.708L8 9.293 5.354 6.646a.5.5 0 1 0-.708.708l3 3a.5.5 0 0 0 .708 0"/>
  </svg>`

	dashSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-dash-square" viewBox="0 0 16 16">
  <path d="M14 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
  <path d="M4 8a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7A.5.5 0 0 1 4 8"/>
</svg>`

	connectSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-chevron-double-left" viewBox="0 0 16 16">
	<path fill-rule="evenodd" d="M8.354 1.646a.5.5 0 0 1 0 .708L2.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0"/>
	<path fill-rule="evenodd" d="M12.354 1.646a.5.5 0 0 1 0 .708L6.707 8l5.647 5.646a.5.5 0 0 1-.708.708l-6-6a.5.5 0 0 1 0-.708l6-6a.5.5 0 0 1 .708 0"/>
  </svg>`
	deleteSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-slash-square-fill" viewBox="0 0 16 16">
	<path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2zm9.354 5.354-6 6a.5.5 0 0 1-.708-.708l6-6a.5.5 0 0 1 .708.708"/>
  </svg>`
)

type (
	LobbyModule struct {
		lobbies map[string]*lobby // key - lobbyID
		lMux    sync.RWMutex
	}

	lobby struct {
		lobbyID string

		sessionID1  string
		sessionID2  string
		playerName1 string
		playerName2 string

		gameSettings Settings

		isStarted bool
	}

	Settings struct {
		FieldType    string `json:"fieldType"`
		FastGame     bool   `json:"fastGame"`
		Experimental bool   `json:"experimental"`
	}

	// LobbyElemResp struct {
	// 	User      string   `json:"user"`
	// 	SessionID string   `json:"sessionID"`
	// 	Settings  Settings `json:"settings"`
	// }
)

func CreateLobbyModule() *LobbyModule {
	return &LobbyModule{
		lobbies: map[string]*lobby{},
	}
}

func (lm *LobbyModule) CreateLobby(
	playerName string,
	sessionID string,
	fieldType string,
	fastGame bool,
	experimental bool,
) (
	lobbyID string,
) {
	lobbyID = generateLobbyID()

	gameLobby := &lobby{
		lobbyID:     lobbyID,
		sessionID1:  sessionID,
		playerName1: playerName,

		gameSettings: Settings{
			FieldType:    fieldType,
			FastGame:     fastGame,
			Experimental: experimental,
		},
	}

	lm.addToMap(gameLobby)

	return
}

func (lm *LobbyModule) AddPlayerToLobby(lobbyID, playerName, sessionID string) error {
	lobby, found := lm.getLobbyByLobbyID(lobbyID)
	if !found {
		return fmt.Errorf("lobby not found")
	}

	lm.lMux.Lock()
	lobby.playerName2 = playerName
	lobby.sessionID2 = sessionID
	lobby.isStarted = true
	lm.lMux.Unlock()

	// start game

	alert.Info("AddPlayerToLobby lobby", lobby)
	return nil
}

func (lm *LobbyModule) CheckActiveGame(lobbyID, playerName, sessionID string) bool {
	lobby, found := lm.getLobbyByLobbyID(lobbyID)
	if !found {
		return false
	}

	if playerName != lobby.playerName1 && playerName != lobby.playerName2 {
		return false
	}

	if sessionID != lobby.sessionID1 && sessionID != lobby.sessionID2 {
		return false
	}

	return lobby.isStarted
}

func (lm *LobbyModule) IsLobbyExists(lobbyID string) (exists bool) {
	_, exists = lm.getLobbyByLobbyID(lobbyID)

	return
}

func (lm *LobbyModule) GetLobbiesTableResponse(sessionID string) string {
	resp := ""

	lm.lMux.RLock()
	for i := range lm.lobbies {
		resp += lm.lobbies[i].formatTableResponse(sessionID)
	}
	lm.lMux.RUnlock()

	return resp
}

func (lm *LobbyModule) RemoveLobby(lobbyID string) {
	lm.removeFromMap(lobbyID)
}

func generateLobbyID() string {
	long := fmt.Sprint(time.Now().UnixNano() / 100)
	return long[len(long)-5:]
}

func (lm *LobbyModule) addToMap(gameLobby *lobby) {
	lm.lMux.Lock()
	lm.lobbies[gameLobby.lobbyID] = gameLobby
	lm.lMux.Unlock()
}

func (lm *LobbyModule) removeFromMap(lobbyID string) {
	lm.lMux.Lock()
	delete(lm.lobbies, lobbyID)
	lm.lMux.Unlock()
}

func (lm *LobbyModule) getLobbyByLobbyID(lobbyID string) (gameLobby *lobby, ok bool) {
	lm.lMux.RLock()
	gameLobby, ok = lm.lobbies[lobbyID]
	lm.lMux.RUnlock()

	if !ok {
		gameLobby = &lobby{}
	}

	return
}

func (l *lobby) formatTableResponse(sessionID string) string {
	if l.isStarted {
		return ""
	}

	alert.Info("formatTableResponse", l.lobbyID, l.sessionID1, sessionID)

	fastGame := "<i>" + dashSvg + "</i>"
	if l.gameSettings.FastGame {
		fastGame = "<i>" + checkSvg + "</i>"
	}

	experimental := "<i>" + dashSvg + "</i>"
	if l.gameSettings.Experimental {
		experimental = "<i>" + checkSvg + "</i>"
	}

	action := "<i onclick='window.connectLobby(\"" + l.lobbyID + "\");' style=\"cursor: pointer;\">" + connectSvg + "</i>"
	selection := ""
	if l.sessionID1 == sessionID {
		action = "<i onclick='window.removeLobby();' style=\"cursor: pointer;\">" + deleteSvg + "</i>"
		selection = "class=\"table-info\""
	}

	result := `<tr ` + selection + `>
              <td>` + l.playerName1 + `</td>
              <td>` + l.gameSettings.FieldType + `</td>
              <td>` + fastGame + `</td>
              <td>` + experimental + `</td>
              <td>` + action + `</td>
        </tr>
`

	return result

	//$(this).attr('sessionID')

}
