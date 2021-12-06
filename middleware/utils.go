package middleware

import (
	"Book-App/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type _Request struct {
	*http.Request
	GetBody interface{}
	Cancel  interface{}
}

func GetLogFile(failed bool) string {
	/*
		Logging requests into a single log file is bad practice because eventually the single file will bloat up, and take more disk space.
		Thus we will create a new log file for each day, so that the log files can be rotated easily ( suppose upload to drive etc using a CRON, etc )
		- The file path where the logs are stored can be taken by the env file.
	*/
	// #####################################
	REQUEST_LOG_FILE_PATH := os.Getenv("REQUEST_LOG_FILE_PATH")
	year, month, date := time.Now().Date()
	if !failed {
		file_name := fmt.Sprintf("request-log-%v-%v-%v.json", year, month, date)
		filePath := REQUEST_LOG_FILE_PATH + file_name
		filePath = filepath.FromSlash(filePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Log file is not created for today, so create one
			ioutil.WriteFile(filePath, []byte(""), os.ModePerm)
		}
		return filePath
	}
	file_name := fmt.Sprintf("failed-request-log-%v-%v-%v.json", year, month, date)
	filePath := REQUEST_LOG_FILE_PATH + file_name
	filePath = filepath.FromSlash(filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Log file is not created for today, so create one
		ioutil.WriteFile(filePath, []byte(""), os.ModePerm)
	}
	return filePath

}

func logToFile(b []byte, failed bool) {
	if failed {
		filePath := GetLogFile(true)
		b = append(b, 44) // Inject an "," at the end of the json
		f, err := os.OpenFile(filePath, os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(b); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		filePath := GetLogFile(false)
		b = append(b, 44) // Inject an "," at the end of the json
		f, err := os.OpenFile(filePath, os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(b); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

}

func logToDB(b []byte, failed bool) {

	if failed {
		var u1 = uuid.Must(uuid.NewV4())
		rawRequest := models.RequestLog{
			ID:         u1,
			RawRequest: string(b),
		}

		models.DB.Create(&rawRequest)
		return
	}
	var u1 = uuid.Must(uuid.NewV4())
	rawRequest := models.FailedRequestLog{
		ID:         u1,
		RawRequest: string(b),
	}
	models.DB.Create(&rawRequest)

}

// RequestLogger is middleware to log the requests in the database or local log file
func RequestLogger(c *gin.Context) {
	j, _ := json.MarshalIndent(_Request{Request: c.Request}, "", "    ")
	switch os.Getenv("LOGGING") {
	case "DB":
		logToDB(j, false)
	case "FILE":
		go logToFile(j, false)
	default:
		panic("LOGGING FLAG INCORRECTLY SET. Logging flag should be DB/FILE")
	}
	c.Next()
}

// RequestLogger is middleware to log the requests in the database or local log file
func FailedRequestLogger(c *gin.Context) {
	j, _ := json.MarshalIndent(_Request{Request: c.Request}, "", "    ")
	switch os.Getenv("LOGGING") {
	case "DB":
		logToDB(j, true)
	case "FILE":
		go logToFile(j, true)
	default:
		panic("LOGGING FLAG INCORRECTLY SET. Logging flag should be DB/FILE")
	}
	c.Next()
}
