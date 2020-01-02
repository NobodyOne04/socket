package client

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

func disconecter(err error, conn net.Conn) {
	if err != nil {
		conn.Close()
		panic("[YOU HAVE DISCONECTED]")
	}
}

func get(conn net.Conn, login string) {
	buf:= make([]byte, 256)
	for {
		size, err:= bufio.NewReader(conn).Read(buf)
		disconecter(err, conn)
		fmt.Printf(string(buf[:size]))
	}
}

func send(conn net.Conn, login string) {
	for {
		mess, _:= bufio.NewReader(os.Stdin).ReadString('\n')
		conn.Write([]byte(fmt.Sprintf("%s:%s", login, mess)))
	}
}

func StartClient(server string, port int) {
	var login string
	buf:= make([]byte, 256)
  	addr, err:= net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", server, port))
	conn, err:= net.DialTCP("tcp", nil, addr)
	size, err:= bufio.NewReader(conn).Read(buf)
	disconecter(err, conn)
	fmt.Println(string(buf[:size]))
	fmt.Scanln(&login)
	conn.Write([]byte(login))
  	go send(conn, login)
  	go get(conn, login)
}