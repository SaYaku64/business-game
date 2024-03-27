$("document").ready(() => {
    var audioElement = document.createElement('audio');
    audioElement.setAttribute('src', '/dice.mp3');
    
    $(".dice").hide()
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
        window.toggleActivePlate(data.current)
    }).fail(function () {
        console.log("checkActiveGame failed")
        window.location.href = "/"
    });

    window.getColorById = function getColorById(id) {
        switch (id) {
            case 0:
                return "blue"
            case 1:
                return "green"
            case 2:
                return "yellow"
            case 3:
                return "red"
        }
    }

    window.showFirstDice = function showFirstDice(id) {
        audioElement.play();
        switch (id) {
            case 1:
                $(".dice.first.first-face").show()
                return
            case 2:
                $(".dice.first.second-face").show()
                return
            case 3:
                $(".dice.first.third-face").show()
                return
            case 4:
                $(".dice.first.fourth-face").show()
                return
            case 5:
                $(".dice.first.fifth-face").show()
                return
            case 6:
                $(".dice.first.sixth-face").show()
                return
        }
    }

    window.showSecondDice = function showSecondDice(id) {
        switch (id) {
            case 1:
                $(".dice.second.first-face").show()
                return
            case 2:
                $(".dice.second.second-face").show()
                return
            case 3:
                $(".dice.second.third-face").show()
                return
            case 4:
                $(".dice.second.fourth-face").show()
                return
            case 5:
                $(".dice.second.fifth-face").show()
                return
            case 6:
                $(".dice.second.sixth-face").show()
                return
        }
    }

    window.toggleActivePlate = function toggleActivePlate(id) {
        $("#player" + id + "-plate").toggleClass("grad-" + window.getColorById(id))
    }

    window.turn = function turn() {
        $('#turnPlate button').prop("disabled", true);
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/turn", function (data) {
            $('#turnPlate').hide();
            $('#turnPlate button').prop("disabled", false);

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
        $('#buyPlate button').prop("disabled", true);
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/buy", function (data) {
            $('#buyPlate').hide();
            $('#buyPlate button').prop("disabled", false);
            //$(".overflow-auto").append("<p class='small'>"+data.msg+"</p>");

            console.log(data)
        }).fail(function (result) {
            alert("error"+result.responseJSON.error)
        });;
    };
    
    window.payRent = function payRent() {
        $('#payRentPlate button').prop("disabled", true);
        $.ajaxSetup({
            headers: {
                'lobbyID': Cookies.get('lobbyID'),
                'sessionID': Cookies.get('sessionID')
            }
        });
        $.get("/api/v1/game/payRent", function (data) {
            $('#payRentPlate').hide();
            $('#payRentPlate button').prop("disabled", false);
            //$(".overflow-auto").append("<p class='small'>"+data.msg+"</p>");

            console.log(data)
        }).fail(function (result) {
            alert("error"+result.responseJSON.error)
        });;
    };
});