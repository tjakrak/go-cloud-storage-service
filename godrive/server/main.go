package main

import (
	"encoding/gob"
	"fmt"
	"godrive/message"
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	msg := &message.MessageHeader{}
	decoder.Decode(msg)
	fmt.Println(msg)

	file, err := os.OpenFile("newfile.txt", os.O_CREATE|os.O_TRUNC|os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	io.Copy(file, conn)
	// use io.CopyN for a certain byte for exp header of file

	// for {
	// 	message := make([]byte, 128)
	// 	n, err := conn.Read(message)
	// 	if err != nil {
	// 		log.Println(err.Error())
	// 		break
	// 	}
	// 	log.Printf("Read %d bytes\n", n)

	// 	if len(message) > 0 {
	// 		msgStr := string(message)
	// 		fmt.Println(msgStr)
	// 	} else {
	// 		break
	// 	}
	// }
}

func main() {
	listener, err := net.Listen("tcp", ":9997")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConnection(conn)
		}
	}
}
