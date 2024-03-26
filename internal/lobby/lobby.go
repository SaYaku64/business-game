package lobby

import (
	"fmt"
	"sync"
	"time"

	a "github.com/SaYaku64/business-game/internal/alert"
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
		lobbies map[string]*Lobby // key - LobbyID
		lMux    sync.RWMutex
	}

	Lobby struct {
		LobbyID string

		SessionIDs  []string
		PlayerNames []string

		// sessionID1  string // deprecated
		// sessionID2  string // deprecated
		// playerName1 string // deprecated
		// playerName2 string // deprecated

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
		lobbies: map[string]*Lobby{},
	}
}

func (lm *LobbyModule) CreateLobby(
	playerName string,
	sessionID string,
	fieldType string,
	fastGame bool,
	experimental bool,
) (
	LobbyID string,
) {
	LobbyID = generateLobbyID()

	gameLobby := &Lobby{
		LobbyID: LobbyID,

		SessionIDs:  []string{sessionID},
		PlayerNames: []string{playerName},

		// sessionID1:  sessionID,
		// playerName1: playerName,

		gameSettings: Settings{
			FieldType:    fieldType,
			FastGame:     fastGame,
			Experimental: experimental,
		},
	}

	lm.addToMap(gameLobby)

	return
}

func (lm *LobbyModule) AddPlayerToLobby(LobbyID, playerName, sessionID string) (lb *Lobby, err error) {
	lobby, found := lm.getLobbyByLobbyID(LobbyID)
	if !found {
		err = fmt.Errorf("lobby not found")

		return
	}

	lm.lMux.Lock()

	lobby.PlayerNames = append(lobby.PlayerNames, playerName)
	lobby.SessionIDs = append(lobby.SessionIDs, sessionID)

	// lobby.playerName2 = playerName
	// lobby.sessionID2 = sessionID
	lobby.isStarted = true
	lm.lMux.Unlock()

	a.Info.Println("AddPlayerToLobby lobby", lobby)

	lb = lobby

	return
}

func (lm *LobbyModule) IsLobbyExists(LobbyID string) (exists bool) {
	_, exists = lm.getLobbyByLobbyID(LobbyID)

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

func (lm *LobbyModule) RemoveLobby(LobbyID string) {
	lm.removeFromMap(LobbyID)
}

func generateLobbyID() string {
	long := fmt.Sprint(time.Now().UnixNano() / 100)
	return long[len(long)-5:]
}

func (lm *LobbyModule) addToMap(gameLobby *Lobby) {
	lm.lMux.Lock()
	lm.lobbies[gameLobby.LobbyID] = gameLobby
	lm.lMux.Unlock()
}

func (lm *LobbyModule) removeFromMap(LobbyID string) {
	lm.lMux.Lock()
	delete(lm.lobbies, LobbyID)
	lm.lMux.Unlock()
}

func (lm *LobbyModule) getLobbyByLobbyID(LobbyID string) (gameLobby *Lobby, ok bool) {
	lm.lMux.RLock()
	gameLobby, ok = lm.lobbies[LobbyID]
	lm.lMux.RUnlock()

	if !ok {
		gameLobby = &Lobby{}
	}

	return
}

func (l *Lobby) formatTableResponse(sessionID string) string {
	if l.isStarted {
		return ""
	}

	a.Info.Printf("formatTableResponse. l.LobbyID: %s, l.SessionIDs[0]: %s, sessionID: %s", l.LobbyID, l.SessionIDs[0], sessionID)

	fastGame := "<i>" + dashSvg + "</i>"
	if l.gameSettings.FastGame {
		fastGame = "<i>" + checkSvg + "</i>"
	}

	experimental := "<i>" + dashSvg + "</i>"
	if l.gameSettings.Experimental {
		experimental = "<i>" + checkSvg + "</i>"
	}

	action := "<i onclick='window.connectLobby(\"" + l.LobbyID + "\");' style=\"cursor: pointer;\">" + connectSvg + "</i>"
	selection := ""
	if l.SessionIDs[0] == sessionID {
		action = "<i onclick='window.removeLobby();' style=\"cursor: pointer;\">" + deleteSvg + "</i>"
		selection = "class=\"table-info\""
	}

	result := `<tr ` + selection + `>
              <td>` + l.PlayerNames[0] + `</td>
              <td>` + l.gameSettings.FieldType + `</td>
              <td>` + fastGame + `</td>
              <td>` + experimental + `</td>
              <td>` + action + `</td>
        </tr>
`

	return result

	//$(this).attr('sessionID')

}
