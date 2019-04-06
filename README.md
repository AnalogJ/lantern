# lantern

<p align="center">
  <a href="https://github.com/AnalogJ/lantern">
  <img width="300" alt="lexicon_view" src="https://github.com/AnalogJ/lantern/blob/master/logo.svg">
  </a>
</p>


Peer into your requests.

# Introduction

Lantern is an open-source debugging proxy (similar to Fiddler/Charles Proxy) that is written in Go and can be
hosted on a server. It uses the Chrome DevTools Inspector as a frontend UI, providing developers with a familiar
interface for inspecting their network traffic.

# Setup

# Usage

```
docker-compose build
docker-compose up
```


# TroubleShooting & Useful Tools


# TODO:

- [-] SSL interception
    - [ ] Add links to the mobileprofile & CA certificate in the Web UI
- [x] backfill requests when new Websocket connection opened
- [x] command responses should be separated from event responses.
- [x] reorganize code, cleanup of event generation
- [ ] README.md documentation.
- [ ] variables for connection strings.
- [ ] better error handling.
- [ ] **STRETCH** live request interception/hijacking & modification.
- [-] Generate custom pdl file with only tabs/domains that we support.
- [x] Use embedded version of Devtools UI.
- [ ] Devtools Theme
    - https://chrome.google.com/webstore/detail/devtools-author/egfhcfdfnajldliefpdoaojgahefjhhi
    - https://chrome.google.com/webstore/detail/devtools-theme-zero-dark/bomhdjeadceaggdgfoefmpeafkjhegbo
- [x] Move the web container into static content served by the API container.

# License

# Contributing


# References

- Logo: [Lantern by zidney](https://thenounproject.com/zidney0721/uploads/?i=1500728)

## Web Frontend

- https://chromedevtools.github.io/devtools-protocol/
- https://github.com/ChromeDevTools
- https://github.com/christian-bromann/devtools-backend
- https://github.com/ChromeDevTools/awesome-chrome-devtools
- https://github.com/ChromeDevTools/devtools-frontend/issues/95
- https://github.com/chromedp/cdproto-gen
- https://blog.hqcodeshop.fi/archives/402-iPhone-Mobile-Profile-for-a-new-CA-root-certificate-Case-CAcert.org.html
- https://www.howtogeek.com/253325/how-to-create-an-ios-configuration-profile-and-alter-hidden-settings/
- https://mdzlog.alcor.net/2012/11/15/decoding-a-mobileconfig-file-containing-a-cisco-ipsec-vpn-configuration/
- https://stackoverflow.com/questions/16727038/how-to-make-mobileconfig-file-on-ios-device
- https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/iPhoneOTAConfiguration/Introduction/Introduction.html
- https://github.com/mritd/strongswan/blob/master/generate-mobileconfig.sh
- https://github.com/ChromeDevTools/devtools-protocol/tree/master/pdl
- https://github.com/ChromeDevTools/devtools-frontend
- https://bit.ly/devtools-contribution-guide
- https://groups.google.com/forum/#!forum/google-chrome-developer-tools
- https://twitter.com/DevToolsCommits
- https://gist.github.com/vbsessa/e337d0add70a71861b8c580d5e16996e


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
- https://github.com/elazarl/goproxy/issues/9
- `curl -v -L -x http://localhost:8082 -p  --proxy-insecure -k https://www.example.com`


