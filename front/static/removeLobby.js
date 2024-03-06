$("document").ready(() => {
    window.removeLobby = function removeLobby() {
        var lobbyToDelete = Cookies.get("lobbyID")

        $.get( "/api/v1/removeLobby", {
            lobbyID: lobbyToDelete,
        }, function() {
            Cookies.set("lobbyID", "", { expires: 0});
            $("#lobbyCreateBtn").prop( "disabled", false );
            $("#btnName").prop( "disabled", false ); // you cannot change name with active session
            $("#createdP").hide();
            window.getLobbies()
            
            var wsMsg = JSON.stringify({
                lobbyID: lobbyToDelete,
                action: 3
            })
            window.WS.sendMsg(wsMsg)
        });
    };
});