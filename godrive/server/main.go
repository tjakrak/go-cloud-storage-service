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
	msg := &message.MessageHeader{}
	decoder.Decode(msg)
	fmt.Println(msg)

	file, err := os.OpenFile("newfile.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}

	byte, err := io.CopyN(file, conn, int64(msg.Size))
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("New file size: %d\n", byte)

	// io.Copy(file, conn)
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
			log.Println("handling connection...")
			go handleConnection(conn)
		}
	}
}
