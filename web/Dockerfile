FROM golang:alpine AS build

RUN apk add --no-cache curl git \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/analogj/lantern/

COPY . .

RUN cd web \
    && dep ensure \
    && go build -o /usr/local/bin/lantern_api cmd/lantern_api.go \
    && ls -alt


###############################################################################
# Runtime Container
###############################################################################
FROM alpine
WORKDIR /srv/lantern

#download devtools frontend code.

RUN apk add curl \
    && curl -O -L https://github.com/ChromeDevTools/devtools-frontend/archive/6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz \
    && tar -xvf 6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz --strip 1 \
    && rm -rf 6bd6d4996c0e1dd424c85540e298338c9aa913ad.tar.gz \
    && curl -o front_end/InspectorBackendCommands.js -L https://chrome-devtools-frontend.appspot.com/serve_file/@38db055e5fc20b2eddca2c829c324fb49de07cbf/InspectorBackendCommands.js \
    && curl -o front_end/SupportedCSSProperties.js -L https://chrome-devtools-frontend.appspot.com/serve_file/@38db055e5fc20b2eddca2c829c324fb49de07cbf/SupportedCSSProperties.js



# root filesystem
COPY web/rootfs /

# copy lantern binary
COPY --from=build /usr/local/bin/lantern_api /usr/local/bin/

# s6 overlay
RUN curl -L -s https://github.com/just-containers/s6-overlay/releases/download/v1.18.1.5/s6-overlay-amd64.tar.gz \
  | tar xvzf - -C / \
 && apk del --no-cache curl

ENTRYPOINT ["/init"]