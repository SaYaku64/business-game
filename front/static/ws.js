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

    disconnect(){
        console.log("disconnect");
        this.mysocket.close()
    }
}