package server

import (
	"fmt"
	"net"
  	"bufio"
  	"strings"
  	"../protector"
)

type Connection struct {
	conn net.Conn
	author string
	protectorObj *protector.SessionProtector 
	currentKey string
}

type Message struct {
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

func get(index int){
	buf:= make([]byte, 256)
	for {
		size, err:= bufio.NewReader(connection[index].conn).Read(buf)
		disconecter(err, connection[index].conn)
		connection[index].currentKey = connection[index].protectorObj.Next_session_key(connection[index].currentKey);
		fmt.Println("client key : " + strings.Split(string(buf[:size]), "\n")[0] + " server key : " + connection[index].currentKey)
		if connection[index].currentKey != strings.Split(string(buf[:size]), "\n")[0] { disconecter(fmt.Errorf("key error"), connection[index].conn); return } 
		if disconecter(err, connection[index].conn) {
			for iter:= 0; iter < len(connection); iter++ {
					if connection[iter] == connection[index] { connection = append(connection[:iter], connection[iter+1:]...); continue }
					connection[iter].conn.Write([]byte(fmt.Sprintf("[%s HAS LEFT CHAT]\n", connection[index].author)))
					fmt.Println(fmt.Sprintf("[%s HAS LEFT CHAT]\n", connection[index].author))
					return
			}
		} else {
			message = append(message, Message{author:connection[index].author, text:strings.Split(string(buf[:size]), "\n")[1]})
		}
	}
}

func send() {
	for {
		for len(message) > 0 {
			deliver:= message[0]
			for iter:= 0; iter < len(connection); iter++ {
				if deliver.author != connection[iter].author { 
					connection[iter].currentKey = connection[iter].protectorObj.Next_session_key(connection[iter].currentKey)
					connection[iter].conn.Write([]byte(connection[iter].currentKey + "\n" + deliver.text))
				}
			}
			message = message[1:]
		}
	}
}

func StartServer(host string, port int, client_number int){
	fmt.Println("[SERVER WAS STARTED]")
	listen, _:= net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	go send()
	for len(connection) < client_number {
		loginBuf, buf:= make([]byte, 256), make([]byte, 256)
		conn, _:= listen.Accept()
		size, _:= bufio.NewReader(conn).Read(buf)
		access:= strings.Split(string(buf[:size]), "\n")
		fmt.Println(access)
		conn.Write([]byte("ENTER YOUR LOGIN : "))
		loginSize, err:= bufio.NewReader(conn).Read(loginBuf)
		if !disconecter(err, conn) {
			fmt.Println(fmt.Sprintf("[%s HAS CONNECTED]", string(loginBuf[:loginSize])))
			newConnection := Connection {conn:conn,
										author:string(loginBuf[:loginSize]),
										protectorObj:protector.NewSessionProtector(access[0]),
										currentKey:access[1]}
			connection = append(connection, newConnection)
			go get(len(connection)-1)
		}
	}
}