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

type MessageHandler func(*message.MessageHeader, *message.Message, net.Conn)

var handlers = map[message.MessageType]MessageHandler{
	message.StorageRequest:   handlePutReq,
	message.RetrievalRequest: handleGetReq,
}

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	msg := &message.Message{}
	//msg := message.New(conn)
	// header := msg.ReadHeader()
	header := &msg.Head
	log.Println("Incoming message header: ", header)
	log.Println("Message type: ", header.Type)

	if _, err := os.Stat("./storage"); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir("./storage", 0755)
			check(err2)
		}
	}
	os.Chdir("./storage")
	newDir, err := os.Getwd()
	check(err)
	log.Printf("Current Working Directory: %s\n", newDir)

	handler := handlers[header.Type]
	if handler != nil {
		handler(header, msg, conn)
	} else {
		log.Println("No handler for message type: ", header.Type)
	}
}

func handlePutReq(header *message.MessageHeader, msg *message.Message, conn net.Conn) {
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	err := decoder.Decode(msg)
	check(err)
	fmt.Println(msg)
	file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", msg.Head.Size)
	bytes, err := io.CopyN(file, bconn, msg.Head.Size)
	check(err)
	log.Printf("New file size: %d\n", bytes)
}

func handleGetReq(header *message.MessageHeader, msg *message.Message, conn net.Conn) {
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
