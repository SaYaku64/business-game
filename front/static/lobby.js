$("document").ready(() => {
    window.WS = new FrontWS()
    window.WS.connect();

    var sessionCookie = Cookies.get('sessionID');
    if (sessionCookie == undefined || sessionCookie == "") {
        $('#lobbyCreateBtn').prop( "disabled", false );
        $('#btnName').prop( "disabled", false );
    } else {
        $('#lobbyCreateBtn').prop( "disabled", true );
        $('#btnName').prop( "disabled", true ); // you cannot change name with active session
        $('#createdP').show();
        document.getElementById("navSessionID").innerHTML = sessionCookie;
    }
    

    $("#lobbyCreateBtn").click(() => {
        var fieldType = $("#fieldType").val()
        
        var fastGame = false
        var experimental = false

        if ($("#fast-game").is(':checked')) { fastGame = true; }
        if ($("#experimental-settings").is(':checked')) { experimental = true; }

        $.post("/api/v1/createLobby", {
            fieldType: fieldType,
            fastGame: fastGame,
            experimental: experimental,
            playerName: Cookies.get('name'),
        }, function (result) {
            if (result.error != null) {
                alert(result.error)
                alert(result.message)
                errorAlert(result.message);
            } else {
                Cookies.set('sessionID', result.sessionID, { expires: 365});
                $('#lobbyCreateBtn').prop( "disabled", true );
                $('#btnName').prop( "disabled", true ); // you cannot change name with active session
                $('#createdP').show();
                document.getElementById("navSessionID").innerHTML = sessionCookie;
                $('#createLobbyContainer').hide()
                $('#connectToLobbyContainer').show()
                window.getLobbies()

                var wsMsg = JSON.stringify({
                    sessionID: result.sessionID,
                    lobbyID: result.sessionID,
                    action: 1
                })
                window.WS.sendMsg(wsMsg)
            }
        }).fail(function(result) {
            errorAlert(result.responseJSON.error);
          });
    });

    
    $("#connectToLobbyBtn").click(() => {
        $('#connectToLobbyContainer').show(); 
        $('#createLobbyContainer').hide()

        window.getLobbies()
    });
    
    window.getLobbies = function getLobbies() {
        $.get( "/api/v1/getLobbiesTable", {
            sessionID: Cookies.get('sessionID'),
        }, function(data) {
            document.getElementById("tableLobbies").innerHTML = data.lobbiesTable
        });
    };

    function errorAlert(message) {
        document.getElementById("errorMenu").innerHTML = "<p class=\"bg-danger dropdown-item text-white font-weight-bold\">"+message+"</p>";
        setTimeout(() => $( "#errorMenu" ).load(window.location.href + " #errorMenu" ), 3500);
    };
});