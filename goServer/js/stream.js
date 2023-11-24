console.log("js_work")

var socket = new WebSocket("ws://localhost:8080/ws?id=js");

socket.onopen = function() {
    alert("Соединение установлено.");
  };
  
  socket.onclose = function(event) {
    if (event.wasClean) {
      alert('Соединение закрыто чисто');
    } else {
      alert('Обрыв соединения');
    }
    alert('Код: ' + event.code + ' причина: ' + event.reason);
  };
  
  socket.onmessage = function(event) {
    //document.getElementById("content").src=event.data

    jsonn = JSON.parse(event.data);

    let f = document.getElementById("f")
    f.textContent = jsonn["f"]
    f.style.color = "red"
    document.getElementById("content").src = jsonn["b64"]
    //console.log(event.data);
  };
  
  socket.onerror = function(error) {
    alert("Ошибка " + error.message);
  };

  async function start(){
    const response = await fetch('/start');
    //alert(response.body);
  }