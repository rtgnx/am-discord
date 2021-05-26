FROM golang:1.16.4-alpine3.13

RUN mkdir -p /go/src/github.com/rtgnx/am-discord
WORKDIR /go/src/github.com/rtgnx/am-discord
COPY . .
RUN go build -o /usr/bin/am-discord
CMD ["/usr/bin/am-discord"]