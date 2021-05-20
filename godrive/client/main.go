package main

import (
	"bytes"
	"fmt"
	"godrive/message"
	"log"
	"net"
	"os"
	"strings"
)

type SendRequest func(string) *message.Message

var msgRequester = map[string]SendRequest{
	"put":    sendPutReq,
	"get":    sendGetReq,
	"search": sendSearchReq,
	"delete": sendDeleteReq,
}

/* Creating message for put request */
func sendPutReq(fileName string) *message.Message {

	if strings.Contains(fileName, "/") {
		path := strings.Split(fileName, "/")
		log.Println(path)
		fileName = path[len(path)-1]
		var dirBuf bytes.Buffer

		for i := 0; i < len(path)-1; i++ {
			fmt.Fprintf(&dirBuf, "%s/", path[i])
		}

		dirBuf.Truncate(dirBuf.Len() - 1)
		directory := dirBuf.String()
		log.Println(directory)
		changeDir(directory)
	}

	fileStat, err := os.Stat(fileName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fileSize := fileStat.Size()
	log.Printf("Sending File Size: %d\n", fileSize)
	return message.New(0, fileSize, fileName)
}

/* Creating message for get request*/
func sendGetReq(fileName string) *message.Message {
	return message.New(1, 0, fileName)
}

/* Creating message for search request */
func sendSearchReq(fileName string) *message.Message {
	return message.New(2, 0, fileName)
}

/* Creating message for delete request */
func sendDeleteReq(fileName string) *message.Message {
	return message.New(3, 0, fileName)
}

/* Getting notification message from server */
func receiveNotification(conn net.Conn) {
	message := make([]byte, 128)
	n, _ := conn.Read(message)
	if n == 0 {
		return
	}
	log.Printf("Read %d bytes\n", n)

	if len(message) > 0 {
		msgStr := string(message)
		fmt.Println(msgStr)
	}
}

/* Client's change directory */
func changeDir(path string) {
	log.Printf("Client current directory: %s\n", path)
	os.Chdir(path)
	path, _ = os.Getwd()
	log.Printf("New client current working directory: %s\n", path)
}

func main() {
	userInput := os.Args
	conn, err := net.Dial("tcp", userInput[1])

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	defer conn.Close()
	var msg *message.Message
	request := msgRequester[userInput[2]]
	if request != nil {
		if len(userInput) == 4 {
			msg = request(userInput[3])
		} else {
			msg = request("")
		}
	} else {
		log.Println("No request: ", request)
		return
	}
	msg.Print()
	msg.Send(conn)
	receiveNotification(conn)
}
