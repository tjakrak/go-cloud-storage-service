package main

import (
	"godrive/message"
	"log"
	"net"
	"os"
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
	fileStat, err := os.Stat(fileName)
	check(err)
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
	return message.New(1, 0, fileName)
}

/* Creating message for delete request */
func sendDeleteReq(fileName string) *message.Message {
	return message.New(1, 0, fileName)
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
		msg = request(userInput[3])
	} else {
		log.Println("No request: ", request)
		return
	}
	msg.Print()
	msg.Send(conn)
}

/* Check error */
func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
