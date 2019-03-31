package main
//
//import (
//	"fmt"
//	"time"
//	"github.com/jinzhu/gorm"
//	"github.com/jinzhu/gorm/dialects/postgres"
//	"github.com/analogj/lantern/common/models"
//	"encoding/json"
//	"log"
//	"os"
//)
//
//func main() {
//
//	db, err := gorm.Open("postgres", "host=database sslmode=disable dbname=lantern user=lantern password=lantern-password")
//	if err != nil {
//		panic(err)
//	}
//
//	defer db.Close()
//
//	db.LogMode(true)
//	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
//
//
//
//	for {
//
//		request := models.DbRequest{
//			Method:        "GET",
//			Url:           "http://www.chromium.org/",
//			Headers:       postgres.Jsonb{json.RawMessage(`{"Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8","Upgrade-Insecure-Requests":"1","User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36"}`)},
//			Body:          "",
//			ContentLength: 0,
//			Host:          "www.chromium.org",
//			RequestedOn: time.Now(),
//		}
//
//		if err = db.Create(&request).Error; err != nil {
//			println("db.Create error!")
//			println(err)
//		}
//
//		fmt.Println("New request ID is:", request.Id)
//
//		time.Sleep(2 * time.Second)
//
//		response := models.DbResponse{
//			RequestId: request.Id,
//			Status: "200 OK",
//			StatusCode:	 200,
//			Headers: 	 postgres.Jsonb{json.RawMessage(`{"Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8","Upgrade-Insecure-Requests":"1","User-Agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36"}`)},
//			Body:		"",
//			ContentLength: 0,
//			MimeType: "text/png",
//			RespondedOn: time.Now(),
//		}
//
//		if err = db.Create(&response).Error; err != nil {
//			println("db.Create error!")
//			println(err)
//		}
//
//		fmt.Println("New response ID is:", response.Id)
//
//
//		// create a new database entry every 5 seconds.
//		time.Sleep(3 * time.Second)
//	}
//}