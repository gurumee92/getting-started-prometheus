package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Println(err)
	}

	defer listener.Close()

	for {
		fmt.Println("listen....")
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
			continue
		}

		defer conn.Close()
		go requestHandler(conn)
	}
}

func requestHandler(conn net.Conn) {
	data := make([]byte, 4096) // 클라이언트와 서버간의 데이터 길이를 정의합니다. (바이트 슬라이스)

	for {
		n, err := conn.Read(data)

		if err != nil {

			log.Println(err)
			return
		}

		_, err = conn.Write(data[:n]) // 해당 클라이언트로부터 데이터를 전송합니다.

		if err != nil {
			log.Println(err)
			return
		}
	}
}
