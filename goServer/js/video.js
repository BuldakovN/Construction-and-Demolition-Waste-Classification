const toBase64 = file => new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = reject;
  });


async function workWithVideo(){
    let loadingGif = document.createElement("img")
    let div = document.getElementById("videoDiv")
    while (div.firstChild) {
      div.removeChild(div.firstChild);
    }
    loadingGif.src = "https://i.gifer.com/ZKZg.gif"
    loadingGif.id = "loadingGif"

    div.appendChild(loadingGif)

    if(document.getElementById("videoFile").files[0]=== undefined){
      alert("Вы не выбрали видео")
      return
    }
    let videoFile = document.getElementById("videoFile").files[0]
    let res_video = await toBase64(videoFile)
    let slv = res_video.split(",")
  
    var dataObj = {
      type : slv[0],
      b64 : slv[1]
    }
  
    var json = JSON.stringify(dataObj);
    while (div.firstChild) {
      div.removeChild(div.firstChild);
    }
    requestVideo(json)
  }

  function requestVideo(json){
    var url = '/sendVideo';
      // Формируем запрос
      response = fetch(url, {
          // Метод, если не указывать, будет использоваться GET
          method: 'POST',
          // Заголовок запроса
          headers: {
          'Content-Type': 'application/json'
          },
          // Данные
          body: json
      })
      .then((resp) => resp.json())
      .then((data)=> finalVideo(data)
      );
  }
  
  function finalVideo(data)
  {
    //console.log(Object.keys(data));
    //console.log(Object.keys(data["video"]));
    //console.log(Object.keys(data["frames"]));
    //console.log(data["result"])
    document.getElementById("result").textContent = data["result"]
    return
  }