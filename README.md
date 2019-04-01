<p align="center">
  <a href="https://github.com/AnalogJ/lantern">
  <img width="300" alt="lexicon_view" src="https://github.com/AnalogJ/lantern/blob/master/logo.svg">
  </a>
</p>



# lantern
Peer into your requests.

# TODO:

- [ ] SSL interception
- [ ] backfill requests when new Websocket connection opened
- [ ] command responses should be separated from event responses.
- [ ] fix mimetype retrieval
- [ ] reorganize code, cleanup of event generation


# References

- Logo: [Lantern by zidney](https://thenounproject.com/zidney0721/uploads/?i=1500728)

## Web Frontend

- https://chromedevtools.github.io/devtools-protocol/
- https://github.com/ChromeDevTools
- https://github.com/christian-bromann/devtools-backend
- https://github.com/ChromeDevTools/awesome-chrome-devtools
- https://github.com/ChromeDevTools/devtools-frontend/issues/95
- https://github.com/chromedp/cdproto-gen

## API/Websockets

- https://github.com/ChromeDevTools/awesome-chrome-devtools#chrome-devtools-protocol
- https://medium.freecodecamp.org/million-websockets-and-go-cc58418460bb
- https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
- https://github.com/chromedp/cdproto
- https://godoc.org/github.com/chromedp/cdproto
- https://gobyexample.com/non-blocking-channel-operations
- https://github.com/kdzwinel/betwixt/blob/master/src/main.js
- https://chromedevtools.github.io/devtools-protocol/tot/Network
- https://medium.com/rungo/anatomy-of-channels-in-go-concurrency-in-go-1ec336086adb


## Database

- http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/

## Proxy

- https://github.com/elazarl/goproxy
- https://gist.github.com/Soulou/6048212
- https://github.com/docker/go-docker/blob/master/hijack.go
- https://stackoverflow.com/questions/23812330/go-hijack-client-connection
- https://golang.org/pkg/net/http/httptest/#NewRequest
- https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
- https://medium.com/@mlowicki/https-proxies-support-in-go-1-10-b956fb501d6bv
- https://github.com/sethgrid/fakettp
- https://github.com/go-httpproxy/httpproxy
- http://speakmy.name/2014/07/29/http-request-debugging-in-go/
- https://github.com/roglew/pappy-proxy/tree/master/pappyproxy/interface/repeater
