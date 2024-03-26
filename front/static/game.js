$("document").ready(() => {
    window.WS = new FrontWS()
    window.WS.connectGame();

    const RollAction_Buy = 1
	const RollAction_PayRent = 2

    var onLoadName = Cookies.get('name');
    if (onLoadName == undefined || onLoadName == "") {
        Cookies.set('name', "Хитрун", { expires: 365 });
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

    $.post("/api/v1/checkActiveGame", {
        lobbyID: Cookies.get('lobbyID'),
        playerName: Cookies.get("name"),
        sessionID: Cookies.get('sessionID')
    }, function (data) {
        if (data.turn) {
            $('#turnPlate').show();
        }
    }).fail(function () {
        console.log("checkActiveGame failed")
        window.location.href = "/"
    });

    window.toggleActivePlate = function toggleActivePlate(id) {
        let plateColor

        switch (id) {
            case 0:
                plateColor = "blue"
                break;
            case 1:
                plateColor = "green"
                break;
            case 2:
                plateColor = "yellow"
                break;
            case 3:
                plateColor = "red"
                break;
        }

        $("#player" + id + "-plate").toggleClass("grad-" + plateColor)
    }

    window.turn = function turn() {
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/turn", function (data) {
            $('#turnPlate').hide();
            //$(".overflow-auto").append("<p class='small'>"+data.msg+"</p>");

            alert("Dices: " + data.result.firstDice + " & " + data.result.secondDice)

            switch (data.result.status) {
                case RollAction_Buy:
                    $('#buyPlate').show();
                    break;
                case RollAction_PayRent:
                    $('#payRentPlate').show();
                    break;
            }

            console.log(data)
        });
    };
    
    window.buy = function buy() {
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/buy", function (data) {
            $('#buyPlate').hide();
            //$(".overflow-auto").append("<p class='small'>"+data.msg+"</p>");

            console.log(data)
        }).fail(function (result) {
            alert("error"+result.responseJSON.error)
        });;
    };
    
    window.payRent = function payRent() {
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/payRent", function (data) {
            $('#payRentPlate').hide();
            //$(".overflow-auto").append("<p class='small'>"+data.msg+"</p>");

            console.log(data)
        }).fail(function (result) {
            alert("error"+result.responseJSON.error)
        });;
    };
});