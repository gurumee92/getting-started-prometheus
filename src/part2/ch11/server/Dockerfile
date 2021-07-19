FROM  golang:1.14-buster

WORKDIR /app
COPY . .

RUN go build -o http_server http_server.go
RUN go build -o tcp_server tcp_server.go

CMD ./http_server & ./tcp_server