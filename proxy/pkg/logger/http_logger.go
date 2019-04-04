package logger

//
//import (
//	"path"
//	"log"
//	"net/http"
//	"github.com/elazarl/goproxy"
//	"fmt"
//	"github.com/elazarl/goproxy/transport"
//	"time"
//	"os"
//)
//
//// HttpLogger is an asynchronous HTTP request/response logger. It traces
//// requests and responses headers in a "log" file in logger directory and dumps
//// their bodies in files prefixed with the session identifiers.
//// Close it to ensure pending items are correctly logged.
//type HttpLogger struct {
//	path  string
//	c     chan *Meta
//	errch chan error
//}
//
//func NewLogger(basepath string) (*HttpLogger, error) {
//	f, err := os.Create(path.Join(basepath, "log"))
//	if err != nil {
//		return nil, err
//	}
//	logger := &HttpLogger{basepath, make(chan *Meta), make(chan error)}
//	go func() {
//		for m := range logger.c {
//			if _, err := m.WriteTo(f); err != nil {
//				log.Println("Can't write meta", err)
//			}
//		}
//		logger.errch <- f.Close()
//	}()
//	return logger, nil
//}
//
//func (logger *HttpLogger) LogResp(resp *http.Response, ctx *goproxy.ProxyCtx) {
//	body := path.Join(logger.path, fmt.Sprintf("%d_resp", ctx.Session))
//	from := ""
//	if ctx.UserData != nil {
//		from = ctx.UserData.(*transport.RoundTripDetails).TCPAddr.String()
//	}
//	if resp == nil {
//		resp = emptyResp
//	} else {
//		resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(body))
//	}
//	logger.LogMeta(&Meta{
//		resp: resp,
//		err:  ctx.Error,
//		t:    time.Now(),
//		sess: ctx.Session,
//		from: from})
//}
//
//var emptyResp = &http.Response{}
//var emptyReq = &http.Request{}
//
//func (logger *HttpLogger) LogReq(req *http.Request, ctx *goproxy.ProxyCtx) {
//	body := path.Join(logger.path, fmt.Sprintf("%d_req", ctx.Session))
//	if req == nil {
//		req = emptyReq
//	} else {
//		req.Body = NewTeeReadCloser(req.Body, NewFileStream(body))
//	}
//	logger.LogMeta(&Meta{
//		req:  req,
//		err:  ctx.Error,
//		t:    time.Now(),
//		sess: ctx.Session,
//		from: req.RemoteAddr})
//}
//
//func (logger *HttpLogger) LogMeta(m *Meta) {
//	logger.c <- m
//}
//
//func (logger *HttpLogger) Close() error {
//	close(logger.c)
//	return <-logger.errch
//}
