$("document").ready(() => {
    window.connectLobby = function connectLobby(lobbyID) {
        var sessionID = Cookies.get("sessionID")
        // var lobbyID = $(this).attr('sessionID')
        console.log("api/v1/connectLobby lobbyID", lobbyID)

        $.get( "/api/v1/connectLobby", {
            sessionID: sessionID,
            lobbyID: lobbyID
        }, function() {
            // Cookies.set("sessionID", "", { expires: 0});
            // $("#lobbyCreateBtn").prop( "disabled", false );
            // $("#btnName").prop( "disabled", false ); // you cannot change name with active session
            // $("#createdP").hide();
            // document.getElementById("navSessionID").innerHTML = "";
            // window.getLobbies()
            console.log("api/v1/connectLobby success")
            
            var wsMsg = JSON.stringify({
                sessionID: sessionID,
                lobbyID: lobbyID,
                action: 2
            })
            window.WS.sendMsg(wsMsg)
        });
    };
});