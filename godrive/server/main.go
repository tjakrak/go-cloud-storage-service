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
	message.SearchRequest:    handleSearchReq,
	message.DeleteRequest:    handleDeleteReq,
}

/* Handling connections from client */
func handleConnection(conn net.Conn) {
	log.Println("handling connection...")
	defer conn.Close()
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	msg := &message.Message{}
	decoder.Decode(msg)
	changeDirectory(msg)

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
func changeDirectory(msg *message.Message) {
	path, err := os.Getwd()
	msg.Check(err)
	log.Printf("Current directory: %s\n", path)
	log.Println(strings.HasSuffix(path, "/storage"))
	if !(strings.HasSuffix(path, "/storage")) {
		if _, err := os.Stat("./storage"); err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir("./storage", 0755)
				msg.Check(err)
			}
		}
		os.Chdir("./storage")
		newDir, err := os.Getwd()
		msg.Check(err)
		log.Printf("New current Working Directory: %s\n", newDir)
	}
}

/* Handling put request */
func handlePutReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	msg.Check(err)
	log.Printf("SERVER -> Header size: %d\n", msg.Head.Size)
	bytes, err := io.CopyN(file, bconn, msg.Head.Size)
	msg.Check(err)
	log.Printf("SERVER -> New file size: %d\n", bytes)
}

/* Handling get request */
func handleGetReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	fileStat, err := os.Stat(msg.Head.Filename)
	msg.Check(err)
	msg.Head.Size = fileStat.Size()
	msg.PutRequest(conn)
}

/* Handling search request */
func handleSearchReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {

}

/* Handling delete request */
func handleDeleteReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	log.Println("Inside server delete request")
	err := os.Remove(msg.Head.Filename)
	msg.Check(err)
}

func main() {
	listener, err := net.Listen("tcp", "192.168.122.212:7777")
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
