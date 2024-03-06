function removeLobby() {
    console.log("removeLobby start")
    $.get( "/api/v1/removeLobby", {
        sessionID: Cookies.get('sessionID'),//$(this).attr('sessionID'),
    }, function() {
        console.log("removeLobby function")
        Cookies.set('sessionID', "", { expires: 0});
        $('#lobbyCreateBtn').prop( "disabled", false );
        $('#btnName').prop( "disabled", false ); // you cannot change name with active session
        $('#createdP').hide();
        document.getElementById("navSessionID").innerHTML = "";
        getLobbies()
    });
    
    console.log("removeLobby end")
};