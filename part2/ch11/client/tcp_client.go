package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	fmt.Println("host: ", host)
	if host == "" {
		host = "localhost"
	}

	addr := host + ":8080"
	client, err := net.Dial("tcp", addr) // 해당 서버로 접속 시도

	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	go func(c net.Conn) {
		data := make([]byte, 4096) // 서버와 클라이언트간의 데이터 길이 동기화

		for {
			n, err := c.Read(data)

			if err != nil {
				log.Fatalln(err)
			}

			log.Println("received: " + string(data[:n]))
			time.Sleep(1 * time.Second)
		}
	}(client)

	go func(c net.Conn) {
		i := 0
		for {
			s := "Hello " + strconv.Itoa(i)
			log.Println("send: " + s)
			_, err := c.Write([]byte(s)) // 서버로 부터 데이터 전송
			if err != nil {
				log.Fatalln(err)
			}

			i++
			time.Sleep(1 * time.Second)
		}
	}(client)

	fmt.Scanln()
}
