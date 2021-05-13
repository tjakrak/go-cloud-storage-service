package main

import (
	"bufio"
	"encoding/gob"
	"godrive/message"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type RequestHandler func(net.Conn, *bufio.Reader, *message.Message)

var handlers = map[message.MessageType]RequestHandler{
	message.StorageRequest:   handlePutReq,
	message.RetrievalRequest: handleGetReq,
}

/* Check error */
func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}

/* Handling connections from client */
func handleConnection(conn net.Conn) {
	log.Println("Inside handle connection")
	defer conn.Close()
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	msg := &message.Message{}
	decoder.Decode(msg)
	changeDirectory()

	log.Printf("Filename: %s", msg.Head.Filename)
	log.Printf("Type: %d", msg.Head.Type)

	header := &msg.Head
	log.Printf("Header type in server: %d", header.Type)
	handler := handlers[header.Type]
	if handler != nil {
		handler(conn, bconn, msg)
	} else {
		log.Println("No handler for message type: ", header.Type)
	}
}

/* Change directory to storage */
func changeDirectory() {
	path, err := os.Getwd()
	check(err)
	if strings.HasSuffix(path, "./storage") {
		if _, err := os.Stat("./storage"); err != nil {
			if os.IsNotExist(err) {
				err2 := os.Mkdir("./storage", 0755)
				check(err2)
			}
		}
		os.Chdir("./storage")
		newDir, err := os.Getwd()
		check(err)
		log.Printf("New current Working Directory: %s\n", newDir)
	}
}

/* Handling put request */
func handlePutReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("SERVER -> Header size: %d\n", msg.Head.Size)
	bytes, err := io.CopyN(file, bconn, msg.Head.Size)
	check(err)
	log.Printf("SERVER -> New file size: %d\n", bytes)
}

/* Handling get request */
func handleGetReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	fileStat, err := os.Stat(msg.Head.Filename)
	check(err)
	msg.Head.Size = fileStat.Size()
	msg.Put(conn)
}

func main() {
	listener, err := net.Listen("tcp", "192.168.122.212:7777")
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
