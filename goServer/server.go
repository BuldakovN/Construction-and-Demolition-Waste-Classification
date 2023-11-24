package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"maslov/hack/db"
	"maslov/hack/entity"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xuri/excelize/v2"
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
  db.InitDataBase()

  clients = make(map[*websocket.Conn]bool)
  r := gin.Default()
  r.Static("/js", "./js")
  r.LoadHTMLGlob("templates/*.html")
  r.Static("/css", "./css")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/Waste.xlsx", "./xls/Waste.xlsx")



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
    db.WriteDataToDB(res["video_name"].(string),res["result"].(string))
		c.JSON(200,res)
    	//log.Println(res["b64"])
		resp.Body.Close()
	})
  
  r.GET("/xls", func(c *gin.Context) {
    Openfile, err := os.Open("xls/Waste.xlsx") //Open the file to be downloaded later
    
  
    if err != nil {		
      http.Error(c.Writer, "File not found.", 404) //return 404 if file is not found
      return
    }
  
    tempBuffer := make([]byte, 512) //Create a byte array to read the file later
    Openfile.Read(tempBuffer) //Read the file into  byte
    FileContentType := http.DetectContentType(tempBuffer) //Get file header
  
    FileStat, _ := Openfile.Stat() //Get info from file
    FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
  
    Filename := "result"
  
    //Set the headers
    c.Writer.Header().Set("Content-Type", FileContentType+";"+Filename)
    c.Writer.Header().Set("Content-Length", FileSize)
  
    Openfile.Seek(0, 0) //We read 512 bytes from the file already so we reset the offset back to 0
    io.Copy(c.Writer, Openfile) //'Copy' the file to the client
    Openfile.Close() //Close after function return
	})

  r.GET("/data",func(c *gin.Context) {
		results := db.GetFromDB()
		
		
    f := excelize.NewFile()
    defer func() {
      if err := f.Close(); err != nil {
          log.Println(err)
      }
    }()
    f.SetCellValue("Sheet1", "A1", "номер")
    f.SetCellValue("Sheet1", "B1", "имя видео")
    f.SetCellValue("Sheet1", "C1", "результат")
    f.SetCellValue("Sheet1", "D1", "дата обработки")
    n:=0
		for i, result := range results {
      n = i+2
      f.SetCellValue("Sheet1", "A"+strconv.Itoa(n), strconv.Itoa(result.Id))
      f.SetCellValue("Sheet1", "B"+strconv.Itoa(n), result.Video_name)
      f.SetCellValue("Sheet1", "C"+strconv.Itoa(n), result.Result)
      f.SetCellValue("Sheet1", "D"+strconv.Itoa(n), result.Date)
      //log.Println(result)
		}
		
    if err := f.SaveAs("xls/Waste.xlsx"); err != nil {
      log.Println(err)
  }
		c.String(200,"выгрузка завершена")

	})
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}