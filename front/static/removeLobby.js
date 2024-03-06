$("document").ready(() => {
    window.removeLobby = function removeLobby() {
        var sessionToDelete = Cookies.get("sessionID")

        $.get( "/api/v1/removeLobby", {
            sessionID: sessionToDelete,
        }, function() {
            Cookies.set("sessionID", "", { expires: 0});
            $("#lobbyCreateBtn").prop( "disabled", false );
            $("#btnName").prop( "disabled", false ); // you cannot change name with active session
            $("#createdP").hide();
            document.getElementById("navSessionID").innerHTML = "";
            window.getLobbies()
            
            var wsMsg = JSON.stringify({
                sessionID: sessionToDelete,
                lobbyID: sessionToDelete,
                action: 3
            })
            window.WS.sendMsg(wsMsg)
        });
    };
});