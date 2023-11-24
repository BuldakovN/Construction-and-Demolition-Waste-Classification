package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"maslov/hack/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//var host ="model"
var host ="localhost"

var clients map[*websocket.Conn]bool

var upgrade = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func wsHandler(w http.ResponseWriter,r *http.Request,id string){

  conn, err:= upgrade.Upgrade(w,r,nil)
  if err!=nil{
    log.Println(err)
    return
  }
  if(id=="js"){
  clients[conn] = true
  }
  defer conn.Close()
  defer delete(clients,conn)

  for{//ждем/получаем    ответа от model 
    _,mess,err:=conn.ReadMessage()
    if err!=nil{
      log.Println(err)
      return
    }
    //log.Println(string(mess))
    writeMessagesToJs(mess)

  }
}

func writeMessagesToJs(mess []byte){
  //log.Println(len(clients))
  for conn := range clients {
		conn.WriteMessage(websocket.TextMessage, mess)
    
	}
}

func main() {
  clients = make(map[*websocket.Conn]bool)
  r := gin.Default()
  r.Static("/js", "./js")
  r.LoadHTMLGlob("templates/*.html")
  r.Static("/css", "./css")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/load.gif", "./css/images/load.gif")
  r.GET("/",func(c *gin.Context) {
    c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			},
		)
  })
  r.GET("/stream",func(c *gin.Context) {
    c.HTML(
			http.StatusOK,
			"index2.html",
			gin.H{
				"title": "Home Page",
			},
		)
  })
  
  r.GET("/ws", func(c *gin.Context) {
    id := c.DefaultQuery("id","-")
    log.Println(id)
    go wsHandler(c.Writer,c.Request,id)
  })
  r.GET("/start", func(c *gin.Context) {
    
    resp, err := http.Get("http://"+host+":5000/start") 
    	if err != nil { 
        	log.Println(err) 
    	}
      answer,_:=io.ReadAll(resp.Body)
      sb:= string(answer)
      //log.Println("Стрим запущен")
      c.String(200,sb)
  })


  r.POST("/sendVideo", func(c *gin.Context) {
		//запись в структуру go для промежуточной обработки данных
		var data entity.FormBase64
		c.ShouldBindJSON(&data)
		bytesRepresentation, err := json.Marshal(data)
		if err != nil {
			log.Fatalln(err)
		}
		//пост запрос на сервер с нейронкой
		log.Println("Запрос на py серрвер")
		resp, err := http.Post("http://"+host+":5000/sendFromPyVideo", "application/json", bytes.NewBuffer(bytesRepresentation)) 
    	if err != nil { 
        	log.Println(err) 
    	}
		log.Println("Ответ с py сервера")
		//defer resp.Body.Close()
		var res map[string]interface{}
		
		//json
    json.NewDecoder(resp.Body).Decode(&res)
		log.Println(res["result"])
		log.Println("Ответ на js")
		//entity.Base64ToFile(res["b64"].(string))
		c.JSON(200,res)
    	//log.Println(res["b64"])
		resp.Body.Close()
	})
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}