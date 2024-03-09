$("document").ready(() => {
    var onLoadName = Cookies.get('name');
    if (onLoadName == undefined || onLoadName == "") {
        Cookies.set('name', "Хитрун", { expires: 365});
    }

    var onLoadSessionID = Cookies.get('sessionID');
    if (onLoadSessionID == undefined || onLoadSessionID == "") {
        console.log("empty sessionID")
        window.location.href = "/"
    }

    var onLoadLobbyID = Cookies.get('lobbyID');
    if (onLoadLobbyID == undefined || onLoadLobbyID == "") {
        console.log("empty lobbyID")
        window.location.href = "/"
    }

    $.post( "/api/v1/checkActiveGame", {
        lobbyID: Cookies.get('lobbyID'),
        playerName: Cookies.get("name"),
        sessionID: Cookies.get('sessionID')
    }).fail(function() {
        console.log("checkActiveGame failed")
        window.location.href = "/"
    });
});