FROM golang AS build

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/analogj/lantern/

COPY . .

RUN ls -alt proxy/vendor && cd proxy \
    && dep ensure \
    && go build -o /usr/local/bin/lantern_proxy cmd/lantern_proxy.go \
    && go build -o /usr/local/bin/lantern_generator cmd/lantern_generator.go \
    && ls -alt

CMD ["lantern_proxy"]