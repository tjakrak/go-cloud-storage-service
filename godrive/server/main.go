package main

import (
	"bufio"
	"godrive/message"
	"io"
	"log"
	"net"
	"os"
)

type RequestHandler func(net.Conn, *bufio.Reader, *message.Message)

var handlers = map[message.MessageType]RequestHandler{
	message.StorageRequest:   handlePutReq,
	message.RetrievalRequest: handleGetReq,
}

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}

func handleConnection(conn net.Conn) {
	log.Println("Inside handle connection")
	defer conn.Close()

	msg := &message.Message{}
	msg.Head.Size = 31
	changeDirectory()

	if msg.Head.Filename == "" {
		msg.Head.Type = 1
		msg.Head.Filename = "test.txt"
	}
	if msg.Head.Type == 0 {
		msg.Get(conn)
	} else if msg.Head.Type == 1 {
		log.Printf("%T\n", conn)
		msg.Put(conn)
	}

	// bconn := bufio.NewReader(conn)
	// decoder := gob.NewDecoder(bconn)
	// msg := &message.Message{}
	// log.Printf("Before decoder: Message header filename in server: %s", msg.Head.Filename)
	// err := decoder.Decode(msg)
	// log.Printf("After decoder: Message header filename in server: %s", msg.Head.Filename)
	// check(err)

	// bconn2 := bufio.NewWriter(conn)
	// encoder := gob.NewEncoder(bconn2)
	// // msg := &message.Message{}
	// log.Printf("Before encoder: Message header filename in server: %s", msg.Head.Filename)
	// err2 := encoder.Encode(msg)
	// log.Printf("After encoder: Message header filename in server: %s", msg.Head.Filename)
	// if err2 != nil {
	// 	log.Fatalln(err2.Error())
	// }

	// file, err := os.OpenFile(msg.Head.Filename, os.O_RDONLY, 0666)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// sz, err := io.Copy(conn, file)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// log.Printf("File size: %d", sz)

	// header := &msg.Head
	// log.Printf("Header type in server: %d", header.Type)
	// changeDirectory(msg.Head)
	// handler := handlers[header.Type]
	// if handler != nil {
	// 	handler(conn, bconn, msg)
	// } else {
	// 	log.Println("No handler for message type: ", header.Type)
	// }
}

func changeDirectory() {
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
}

func handlePutReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	// msg.Get(conn)
	file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", msg.Head.Size)
	bytes, err := io.CopyN(file, bconn, msg.Head.Size)
	check(err)
	log.Printf("New file size: %d\n", bytes)
}

func handleGetReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	//bconn2 := bufio.NewWriter(conn)
	// file, err := os.OpenFile(msg.Head.Filename, os.O_RDONLY, 0666)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// msg.Put(conn)

	// sz, err := io.Copy(conn, file)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// log.Printf("File size: %d", sz)
	// ensure all data is written out to the socket
	// bconn.Flush()
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
