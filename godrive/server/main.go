package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"godrive/message"
	"io"
	"io/ioutil"
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
	handler := handlers[header.Type]
	if handler != nil {
		handler(conn, bconn, msg)
	} else {
		log.Println("No handler for message type: ", header.Type)
	}
}

/* Change directory to storage */
func changeDirectory(msg *message.Message) string {
	path, err := os.Getwd()
	msg.Check(err)
	log.Printf("Current directory: %s\n", path)
	if !(strings.HasSuffix(path, "/storage")) {
		if _, err := os.Stat("./storage"); err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir("./storage", 0755)
				msg.Check(err)
			}
		}
		os.Chdir("./storage")
		path, err = os.Getwd()
		msg.Check(err)
		log.Printf("New current working directory: %s\n", path)
	}
	return path
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
	path := changeDirectory(msg)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

/* Handling delete request */
func handleDeleteReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	err := os.Remove(msg.Head.Filename)
	msg.Check(err)
	log.Printf("Deleted file: %s", msg.Head.Filename)
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
