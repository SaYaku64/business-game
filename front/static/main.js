$("document").ready(() => {
    var sessionCookie = Cookies.get('sessionID');
    if (sessionCookie == undefined || sessionCookie == "") {
        $.get( "/api/v1/getSessionID", function( data ) {
            Cookies.set('sessionID', data.sessionID, { expires: 365});
            document.getElementById("navSessionID").innerHTML = data.sessionID;
        });
    } else {
        document.getElementById("navSessionID").innerHTML = sessionCookie;
    }

    window.checkActiveLobby = function () {
        var lobbyCookie = Cookies.get('lobbyID');
        if (lobbyCookie == undefined || lobbyCookie == "") {
            $('#lobbyCreateBtn').prop( "disabled", false );
            $('#btnName').prop( "disabled", false );
            $('#createdP').hide();
        } else {
            $('#lobbyCreateBtn').prop( "disabled", true );
            $('#btnName').prop( "disabled", true ); // you cannot change name with active session
            $('#createdP').show();
        }
    }
    window.checkActiveLobby();

    var onLoadName = Cookies.get('name');
    if (onLoadName == undefined || onLoadName == "") {
        $("#modalName").show()
        $("#modalClose").hide()
    } else {
        $("#nameField").val(onLoadName)
        $("#inputName").val(onLoadName)
        $("#modalName").hide()
    }

    $("#btnRand").click(() => {
        $.get( "/api/v1/randomName", function( data ) {
            $("#inputName").val(data);
        });
    });
    
    $("#btnSubmit").click(() => {
        var name = $("#inputName").val();
        Cookies.set('name', name, { expires: 365});
        $("#nameField").val(name)
        $("#modalName").hide()
    });
    
    $("#btnName").click(() => {
        $("#modalName").show()
        // $("#inputName").val(onLoadName)
        $("#modalClose").show()
    });
    
    $("#modalClose").click(() => {
        $("#modalName").hide()
    });

});
