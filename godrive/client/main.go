package main

import (
	"godrive/message"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":9997") // connect to localhost port 9999
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// something := message.SearchRequest
	// m := message.Message{Name: "GoDrive"}
	// fmt.Println(m, something)
	defer conn.Close()
	msg := message.New(message.SearchRequest, 300)
	msg.Print()
	msg.Send(conn) // pass in our connection
	file, err := os.OpenFile("test.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	io.Copy(conn, file)
	// scanner := bufio.NewScanner(os.Stdin)
	// for {
	// 	fmt.Print("message> ")
	// 	result := scanner.Scan() // Reads up to a \n newline character
	// 	if result == false {
	// 		break
	// 	}

	// 	message := scanner.Text()
	// 	if len(message) != 0 {
	// 		msgBytes := make([]byte, 128) // create a slice w/ 128 byte capacity
	// 		copy(msgBytes, message)       // won't copy more than 128 bytes
	// 		conn.Write(msgBytes)
	// 	}
	// }
}
