package main

import (
	"bufio"
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
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	// msg := &message.MessageHeader{}
	msg := &message.Message{}
	err := decoder.Decode(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(msg)

	file, err := os.OpenFile("newfile.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("Header size: %d\n", msg.Head.Size)
	bytes, err := io.CopyN(file, conn, msg.Head.Size)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("New file size: %d\n", bytes)

	// io.Copy(file, conn)
	// use io.CopyN for a certain byte for exp header of file
}

func main() {
	listener, err := net.Listen("tcp", "192.168.122.212:9995")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			log.Println("handling connection...")
			go handleConnection(conn)
		}
	}
}
