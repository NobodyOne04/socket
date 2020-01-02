package server

import (
	"fmt"
	"net"
  	"bufio"
)

type Connection struct {
	conn net.Conn
	author string
}

type Message struct{
	text string
	author string
}

var (
	connection []Connection
	message []Message
)

func disconecter(err error, conn net.Conn) bool {
	if err != nil {
		conn.Close()
		return true
	} else { return false }
}

func get(client Connection){
	buf:= make([]byte,256)
	for {
		size, err:= bufio.NewReader(client.conn).Read(buf)
		if disconecter(err, client.conn){
			for iter:= 0; iter < len(connection); iter++ {
					if connection[iter] == client { connection = append(connection[:iter], connection[iter+1:]...)}
					connection[iter].conn.Write([]byte(fmt.Sprintf("[%s HAS LEFT CHAT]\n", client.author)))
					fmt.Println(fmt.Sprintf("[%s HAS LEFT CHAT]\n", client.author))
					return
			}
		} else {
			message = append(message, Message{author:client.author, text:string(buf[:size])})
		}
	}
}

func send() {
	for {
		for len(message) > 0 {
			deliver:= message[0]
			buf:= []byte(deliver.text)
			for iter:= 0; iter < len(connection); iter++ {
				if deliver.author != connection[iter].author { connection[iter].conn.Write(buf) }
			}
			message = message[1:]
		}
	}
}

func StartServer(host string, port int){
	fmt.Println("[SERVER WAS STARTED]")
	listen, _:= net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	go send()
	for {
		buf:= make([]byte, 256)
		conn, err:= listen.Accept()
		conn.Write([]byte("ENTER YOUR LOGIN : "))
		size, err:= bufio.NewReader(conn).Read(buf)
		if !disconecter(err, conn) {
			fmt.Println(fmt.Sprintf("[%s HAS CONNECTED]", string(buf[:size])))
			if connection != nil {
				for iter:= 0; iter < len(connection); iter++ {
					connection[iter].conn.Write([]byte(fmt.Sprintf("[%s HAS JOINED CHAT]\n", string(buf[:size]))))
				}
			}
			connection = append(connection, Connection{conn:conn,author:string(buf[:size])})
			go get(Connection{conn:conn, author:string(buf[:size])})
		}
	}
}