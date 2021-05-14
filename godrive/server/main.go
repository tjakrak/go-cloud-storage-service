package main

import (
	"bufio"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"godrive/message"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
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
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	var note string
	log.Printf("SERVER PUT -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		msg.Check(err)
		log.Printf("SERVER PUT -> Header size: %d\n", msg.Head.Size)
		bytes, err := io.CopyN(file, bconn, msg.Head.Size)
		msg.Check(err)
		log.Printf("SERVER PUT -> New file size: %d\n", bytes)
		file.Seek(0, 0)
		hashFile(file, msg.Head.Filename)
		note = "File " + msg.Head.Filename + " is stored"
		defer file.Close()
	} else {
		note = "File already exists"
	}
	sendMessage(note, conn)
}

/* Handling get request */
func handleGetReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	log.Printf("SERVER GET -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("File doesn't exist")
		msg.Body = "File doesn't exist"
		return
	} else {
		fileStat, err := os.Stat(msg.Head.Filename)
		msg.Check(err)
		msg.Head.Size = fileStat.Size()
		note := msg.Head.Filename + " is received"
		msg.Body = note
		msg.PutRequest(conn)
	}
}

/* Handling search request */
func handleSearchReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	files, err := ioutil.ReadDir(path)
	msg.Check(err)

	var returnFiles string
	r, _ := regexp.Compile(".*?" + msg.Head.Filename + ".*")
	for _, f := range files {
		if r.MatchString(f.Name()) {
			fmt.Println(f.Name())
			returnFiles += f.Name() + "\n"
		}
	}
	sendMessage(returnFiles, conn)
}

/* Handling delete request */
func handleDeleteReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	var note string
	log.Printf("SERVER DEL -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		note = "File doesn't exist"
	} else {
		err := os.Remove(msg.Head.Filename)
		msg.Check(err)
		note = msg.Head.Filename + " is deleted"
	}
	fmt.Println(note)
	sendMessage(note, conn)
}

/* Sending notification message to client */
func sendMessage(note string, conn net.Conn) {
	conn.Write([]byte(note))
}

/* Get the md5sum of the file */
func hashFile(file *os.File, fileName string) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Hash: %x\n", hash.Sum(nil))
	fileSplit := strings.Split(fileName, ".")
	fileHash := fileSplit[0] + "hash.txt"
	writeToFile(hex.EncodeToString(hash.Sum(nil)), fileHash)
}

/* Write to filehash.txt */
func writeToFile(hash string, file string) {
	openFile(file)
	hashData := []byte(hash)
	err := ioutil.WriteFile(file, hashData, 0777)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

/* Open file */
func openFile(file string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("failed opening file: %s", err.Error())
		return
	}
	defer f.Close()
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
