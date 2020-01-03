package client

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"strings"
	"../protector"
)

var (
	currentKey string
	protectorObj *protector.SessionProtector
)

func recovery() { if recv:= recover(); recv != nil {fmt.Println(recv); os.Exit(1)} }

func disconecter(err error, conn net.Conn) {
	if err != nil {
		conn.Close()
		panic(err)
	}
}

func get(conn net.Conn, login string) {
	defer recovery()
	buf:= make([]byte, 256)
	for {
		size, _:= bufio.NewReader(conn).Read(buf)
		currentKey = protectorObj.Next_session_key(currentKey)
		msg:= strings.Split(string(buf[:size]), "\n") 
		fmt.Println("client key : " + string(currentKey) + " server key : " + msg[0])
		if msg[0] != currentKey { disconecter(fmt.Errorf("Server protector key error"), conn) }
		fmt.Println(msg[1])
	}
}

func send(conn net.Conn, login string) {
	defer recovery()
	for {
		mess, _:= bufio.NewReader(os.Stdin).ReadString('\n')
		currentKey = protectorObj.Next_session_key(currentKey)
		conn.Write([]byte(currentKey + "\n" + fmt.Sprintf("%s:%s", login, mess)))
	}
}

func StartClient(server string, port int) {
	defer recovery()
	var login string
	initHash:= protector.Get_hash_str()
	currentKey = protector.Get_session_key()
	protectorObj = protector.NewSessionProtector(initHash)
	buf:= make([]byte, 256)
  	addr, err:= net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", server, port))
	conn, err:= net.DialTCP("tcp", nil, addr)
	disconecter(err, conn)
	conn.Write([]byte(initHash + "\n" + currentKey))
	size, err:= bufio.NewReader(conn).Read(buf)
	fmt.Println(string(buf[:size]))
	fmt.Scanln(&login)
	conn.Write([]byte(login))
  	go send(conn, login)
  	go get(conn, login)
}
