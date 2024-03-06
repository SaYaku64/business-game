package lobby

import (
	"fmt"
	"sync"
	"time"
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

	queryJS = `$.get( "/api/v1/removeLobby", {
	sessionID: Cookies.get("sessionID"),//$(this).attr("sessionID"),
}, function() {
	console.log("removeLobby function")
	Cookies.set("sessionID", "", { expires: 0});
	$("#lobbyCreateBtn").prop( "disabled", false );
	$("#btnName").prop( "disabled", false ); // you cannot change name with active session
	$("#createdP").hide();
	document.getElementById("navSessionID").innerHTML = "";
	$.get( "/api/v1/getLobbiesTable", {
		sessionID: Cookies.get("sessionID"),
	}, function(data) {
		document.getElementById("tableLobbies").innerHTML = data.lobbiesTable
	});
});`
)

type (
	LobbyModule struct {
		lobbies map[string]*lobby // key - SessionID1
		lMux    sync.RWMutex
	}

	lobby struct {
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
	fieldType string,
	fastGame bool,
	experimental bool,
) (
	plrSessionID string,
) {
	plrSessionID = generateSessionID()

	gameLobby := &lobby{
		sessionID1:  plrSessionID,
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

func (lm *LobbyModule) GetLobbiesTableResponse(sessionID string) string {
	resp := ""

	// test
	gameLobby := &lobby{
		sessionID1:  "123123",
		playerName1: "playerName",

		gameSettings: Settings{
			FieldType:    "test",
			FastGame:     true,
			Experimental: true,
		},
	}
	lm.lobbies["123123"] = gameLobby
	// test

	lm.lMux.RLock()
	for i := range lm.lobbies {
		resp += lm.lobbies[i].formatTableResponse(sessionID)
	}
	lm.lMux.RUnlock()

	return resp
}

func (lm *LobbyModule) RemoveLobby(sessionID string) {
	lm.removeFromMap(sessionID)
}

func generateSessionID() string {
	return fmt.Sprint(time.Now().UnixNano())
}

func (lm *LobbyModule) addToMap(gameLobby *lobby) {
	lm.lMux.Lock()
	lm.lobbies[gameLobby.sessionID1] = gameLobby
	lm.lMux.Unlock()
}

func (lm *LobbyModule) removeFromMap(sessionID string) {
	lm.lMux.Lock()
	delete(lm.lobbies, sessionID)
	lm.lMux.Unlock()
}

// func (lm *LobbyModule) getLobbyBySessionID(sessionID string) *lobby {
// 	lm.lMux.RLock()
// 	gameLobby, ok := lm.lobbies[sessionID]
// 	lm.lMux.RUnlock()

// 	if !ok {
// 		return &lobby{}
// 	}

// 	return gameLobby
// }

func (l *lobby) formatTableResponse(sessionID string) string {
	fastGame := "<i>" + dashSvg + "</i>"
	if l.gameSettings.FastGame {
		fastGame = "<i>" + checkSvg + "</i>"
	}

	experimental := "<i>" + dashSvg + "</i>"
	if l.gameSettings.Experimental {
		experimental = "<i>" + checkSvg + "</i>"
	}

	action := "<i id=\"connectAndPlay\" sessionID=" + sessionID + " style=\"cursor: pointer;\">" + connectSvg + "</i>"
	selection := ""
	if l.sessionID1 == sessionID {
		action = "<i class=\"removeLobby\" onclick='" + queryJS + "' style=\"cursor: pointer;\">" + deleteSvg + "</i>"
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

	//$(this).attr('terr')

}
