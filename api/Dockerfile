FROM golang:alpine AS build

RUN apk add --no-cache curl git \
    && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/analogj/lantern/

COPY . .

RUN cd api \
    && dep ensure \
    && go build -o /usr/local/bin/lantern_api cmd/lantern_api.go \
    && ls -alt


###############################################################################
# Runtime Container
###############################################################################
FROM alpine

# root filesystem
COPY api/rootfs /

# copy lantern binary
COPY --from=build /usr/local/bin/lantern_api /usr/local/bin/

# s6 overlay
RUN apk add --no-cache curl \
 && curl -L -s https://github.com/just-containers/s6-overlay/releases/download/v1.18.1.5/s6-overlay-amd64.tar.gz \
  | tar xvzf - -C / \
 && apk del --no-cache curl

ENTRYPOINT ["/init"]