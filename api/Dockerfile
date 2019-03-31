FROM golang AS build

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/analogj/lantern/

COPY . .

RUN cd api \
    && dep ensure \
    && go build -o /usr/local/bin/lantern_api cmd/lantern_api.go \
    && ls -alt

CMD ["lantern_api"]