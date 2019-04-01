package main

import (
	"github.com/elazarl/goproxy"
	_ "github.com/lib/pq"
    "encoding/base64"

	"log"
	"flag"
	"net/http"
	"github.com/jinzhu/gorm"
	"os"
	"time"
	"github.com/analogj/lantern/common/models"
	"strings"
	"io/ioutil"
	"bytes"
)

func main() {
	//parse flags
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":9000", "proxy listen address")
	flag.Parse()

	//open database connection
	db, err := gorm.Open("postgres", "host=database sslmode=disable dbname=lantern user=lantern password=lantern-password")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))



	//start proxy server.

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		log.Printf("printing intercepted request: %v %v\n", req, ctx)

		headers := map[string]interface{}{}
		for k, v := range req.Header {
			headers[k] = strings.Join(v, ";")
		}

		dbRequest := models.DbRequest{
			Method:        req.Method,
			Url:           req.RequestURI,
			Headers:       headers,
			Body:          base64EncodedRequestBody(req),
			ContentLength: req.ContentLength,
			Host:          req.Host,
			RequestedOn:   time.Now(),
		}

		if err = db.Create(&dbRequest).Error; err != nil {
			println("db.Create error!")
			println(err)
		}

		ctx.UserData = map[string]uint{
			"RequestId": dbRequest.Id,
		}

		//logger.LogReq(req, ctx)
		return req, nil
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

		log.Printf("printing intercepted response: %v %v\n", resp, ctx)
		//logger.LogResp(resp, ctx)

		headers := map[string]interface{}{}
		for k, v := range resp.Header {
			headers[k] = strings.Join(v, ";")
		}


		dbResponse := models.DbResponse{
			RequestId:     ctx.UserData.(map[string]uint)["RequestId"],
			Status:        resp.Status,
			StatusCode:    resp.StatusCode,
			Headers:       headers,
			Body:          base64EncodedResponseBody(resp),
			ContentLength: resp.ContentLength,
			MimeType:      resp.Header.Get("Content-Type"),
			RespondedOn:   time.Now(),
		}

		if err = db.Create(&dbResponse).Error; err != nil {
			println("db.Create error!")
			println(err)
		}

		return resp
	})

	log.Printf("Starting proxy server on %v\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, proxy))
}

func base64EncodedRequestBody(req *http.Request) string {

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _= ioutil.ReadAll(req.Body)
	}
	// Restore the io.ReadCloser to its original state
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	return base64.StdEncoding.EncodeToString(bodyBytes)
}

func base64EncodedResponseBody(resp *http.Response) string {

	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _= ioutil.ReadAll(resp.Body)
	}
	// Restore the io.ReadCloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	return base64.StdEncoding.EncodeToString(bodyBytes)
}
