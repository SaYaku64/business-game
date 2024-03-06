$("document").ready(() => {
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
