package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
        "com/client"
        //"com/logger"
        "com/models"
	"strconv"
)

func showError(e error) {
    if e != nil {
        log.Println(e)
    }
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}
var obs *client.Client = nil
const (
  host     = "192.168.0.149"
  port     = 8635
  user     = "mryan"
  password = "termita000"
  dbname   = "noti"
)
func putObject(contenido []byte,id int64) {
        input := new(models.PutObjectInput)
        input.Bucket = "notificaciones"
        input.Object = "noti1"+strconv.FormatInt(id,10)+".pdf"
        input.ACL = models.PUBLIC_READ_WRITE
        input.Body = contenido
        //input.SourceFile = f.Name //"C:\\install.log"

        //requst, output := obs.PutObject(input)
        requst,output:=obs.PutObject(input)
        if requst.err != nil {
          fmt.Printf("err:%s,statusCode:%d,code:%s,message:%s\n", requst.Err, requst.StatusCode, requst.Code, requst.Message)
          fmt.Printf("ETag:%s,VersionId:%s\n", output.ETag, output.VersionId)
        }        
}
func main() {
    obs = client.FactoryEx("YGKTCULBWQMWXTB6SHUJ", "sgA4tMY8o41uAXeS4Zzpl6N5SIxQkPNrWWOeFfI9", "", "", "https://obs.sa-argentina-1.telefonicaopencloud.com", true)

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=require",
    host, port, user, password, dbname)
    /*--verify-ca sslrootcert=/home/linux/ca-bundle.pem",*/

    db, err := sql.Open("postgres", psqlInfo)
    check(err)
    defer db.Close()
      err = db.Ping()
      if err != nil {
        panic(err)
      }

    rows, err := db.Query("select id,texto_pdf from notificaciones where texto_pdf is not null order by id;")
    showError(err)
    check(err)
    defer rows.Close()
    for rows.Next() {
       var id int64
       texto:= make([]byte,16000)
       err = rows.Scan(&id, &texto)

       showError(err)
       if err != nil {
          // handle this error
	  fmt.Printf("error haciend selectpara id:%d <--\n", id)
          //panic(err)
       }
       //sube la nofificacion firmada
       putObject(texto,id)
        //fmt.Print(".")
    }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }
}
