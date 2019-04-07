package main

import (
	"encoding/base64"
	"github.com/elazarl/goproxy"
	_ "github.com/lib/pq"

	"bytes"
	"flag"
	"github.com/analogj/lantern/common/models"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/analogj/lantern/proxy/pkg/cert"
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

	// set the CA
	cert.SetCA("/srv/lantern/certs/ca.crt", "/srv/lantern/certs/ca.key.pem")

	//start proxy server.
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*$"))).HandleConnect(goproxy.AlwaysReject)
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		//log.Printf("printing intercepted request: %v %v\n", req, ctx)

		//jsonBytes, err := json.Marshal(req)
		//if err == nil {
		//	//http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		//	log.Printf("========= DUMPED REQUEST: %q", string(jsonBytes))
		//
		//} else {
		//}

		headers := map[string]interface{}{}
		for k, v := range req.Header {
			headers[k] = strings.Join(v, ";")
		}

		encodedBody, length,  _ := base64EncodedRequestBody(req)
		dbRequest := models.DbRequest{
			Method:        req.Method,
			Url:           req.URL.String(),
			Headers:       headers,
			Body:          encodedBody,
			ContentLength: length,
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

		encodedBody, length, mimeType := base64EncodedResponseBody(resp)
		dbResponse := models.DbResponse{
			RequestId:     ctx.UserData.(map[string]uint)["RequestId"],
			Status:        resp.Status,
			StatusCode:    resp.StatusCode,
			Headers:       headers,
			Body:          encodedBody,
			ContentLength: length,
			MimeType:      mimeType,
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

func base64EncodedRequestBody(req *http.Request) (string, int64, string) {

	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(req.Body)

		// Restore the io.ReadCloser to its original state
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Use the content

	mimeType := http.DetectContentType(bodyBytes)

	return base64.StdEncoding.EncodeToString(bodyBytes), int64(len(bodyBytes)), mimeType
}

func base64EncodedResponseBody(resp *http.Response) (string, int64, string) {

	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Use the content
	mimeType := http.DetectContentType(bodyBytes)

	return base64.StdEncoding.EncodeToString(bodyBytes), int64(len(bodyBytes)), mimeType
}
