FROM golang:1.16

WORKDIR /app

COPY . .

RUN go mod init app
RUN go get
RUN go build -o app .

EXPOSE 8080

CMD ["/app/app"]