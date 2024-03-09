$("document").ready(() => {
    window.connectLobby = function connectLobby(lobbyID) {
        var sessionID = Cookies.get("sessionID")

        $.post( "/api/v1/connectLobby", {
            lobbyID: lobbyID,
            playerName: Cookies.get("name"),
            sessionID: sessionID
        }, function(result) {
            if (result.error != null) {
                window.errorAlert(result.message);

                return
            }
            Cookies.set('lobbyID', lobbyID, { expires: 365});

            window.getLobbies();
            
            var wsMsg = JSON.stringify({
                sessionID: sessionID,
                lobbyID: lobbyID,
                action: 2
            })
            window.WS.sendMsg(wsMsg)

        }).fail(function(result) {
            window.errorAlert(result.responseJSON.error);
          });
    };

    
    window.redirectToLobby = function redirectToLobby() {
        window.location.href = "/game"
    };
});