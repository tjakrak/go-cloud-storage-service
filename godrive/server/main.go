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

	file, err := os.OpenFile("newfile.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}

	var size int64 // get from header
	_, err = io.CopyN(file, conn, size)
	if err != nil {
		fmt.Println(err.Error())
	}
	// use io.CopyN for a certain byte for exp header of file
}

func main() {
	listener, err := net.Listen("tcp", ":9995")
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
