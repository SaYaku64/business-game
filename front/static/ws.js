class FrontWS{
    constructor(){
        this.mysocket =  null;
        // this.vMsgContainer = document.getElementById("msgcontainer");
        // this.vMsgIpt = document.getElementById("ipt");
    }

    // showMessage(text, myself){
    //     var div = document.createElement("div"); 
    //     div.innerHTML = text;
    //     var cself = (myself)? "self" : "";
    //     div.className="msg " + cself;
    //     this.vMsgContainer.appendChild(div);
    // }

    // send(){
    //     var txt = this.vMsgIpt.value; 
    //     this.showMessage("<b>Me</b> " + txt,true);
    //     this.mysocket.send(txt);
    //     this.vMsgIpt.value = ""
    // }
 
    // keypress(e){
    //     if (e.keyCode == 13) {
    //         this.send();
    //     }
    // }

    sendMsg(msg){
        console.log("sendMsg", msg);
        this.mysocket.send(msg);
    }

    acceptMsg(msg){
        console.log("acceptMsg", msg)
        if (msg == "redirect") {
            window.redirectToLobby();
        }
    }

    acceptMsgGame(msg){
        console.log("acceptMsgGame", msg)

        if (msg.match("^{")) {
            let jObj = JSON.parse(msg)
            if (jObj.struct == true) {
                // $(".overflow-auto").append("<p class='small'>"+jObj.msg+"</p>");
                $("#chatBottom").prepend("<p class='small'>"+jObj.msg+"</p>");
                return
            }
        }
        
        if (msg == "take turn") {
            $('#turnPlate').show();

            return
        }

    }

    connect(){
        console.log("connect");
        var socket = new WebSocket("ws://localhost:8080/ws");
        this.mysocket = socket;

        socket.onmessage = (e)=>{
        //    this.showMessage(e.data, false);
           this.acceptMsg(e.data)
        }
        
        socket.onopen =  ()=> {
           console.log("socket opened")
        };  
        socket.onclose = ()=> {
           console.log("socket closed")
        }
    }

    connectGame(){
        console.log("connectGame");
        var params = "?lobbyID="+Cookies.get('lobbyID')+"&sessionID="+Cookies.get('sessionID')

        var socket = new WebSocket("ws://localhost:8080/ws/game"+params);
        this.mysocket = socket;

        socket.onmessage = (e)=>{
        //    this.showMessage(e.data, false);
           this.acceptMsgGame(e.data)
        }
        
        socket.onopen =  ()=> {
           console.log("socket opened")
        };  
        socket.onclose = ()=> {
           console.log("socket closed")
        }
    }

    disconnect(){
        console.log("disconnect");
        this.mysocket.close()
    }
}