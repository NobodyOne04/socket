package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
  "bufio"
	"os"
)

type clientConn struct {
	conn net.Conn
	author string
}

type messageData struct{
	text string
	author string
}

var (
	host   = flag.String("h", "localhost", "host")
	port   = flag.Int("p", 9090, "port")
	server =flag.Bool("s",false,"Start as server")
	serverIP=flag.String("ip","78.140.8.143","Server ip")
	message []messageData
	connection []clientConn
	wg sync.WaitGroup
	start bool
)

func processErr(err error) {
	if err != nil {
		panic(err)
	}
}

func recv() {
	r := recover()
	if r != nil {
		fmt.Println(r)
	}
}

func startServer(){
	fmt.Println("[SERVER WAS STARTED]")
	addr:=fmt.Sprintf("%s:%d",*host,*port)
	listen,err:=net.Listen("tcp",addr)
	processErr(err)
	for {
		buf:=make([]byte,256)
		conn,err:=listen.Accept()
		if err!=nil{
			fmt.Println("[CAN`T ACCEPT CLIENT CONNECTION]")
		}
		conn.Write([]byte("ENTER YOUR LOGIN : "))
		size,err:=bufio.NewReader(conn).Read(buf)
		if err!=nil{
			conn.Close()
			fmt.Println("[CLIENT WAS DISCONECTED]")
		}
		log:=string(buf[:size])
		notice:=fmt.Sprintf("[%s HAS CONNECTED]",string(buf[:size]))
		fmt.Println(notice)
		if connection !=nil{
			iter:=0
			for iter<len(connection){
				notice:=fmt.Sprintf("[%s HAS JOINED CHAT]",log)
				connection[iter].conn.Write([]byte(notice))
				iter++
			}
			fmt.Println(connection)
		}
		newClient:=clientConn{conn:conn,author:log}
		connection=append(connection,newClient)
		go processGet(newClient)
		fmt.Println(connection)
		if start{
			go processSend()
			start=false
		}
	}
}

func processSend(){
	for{
		iter:=0
		for len(message)>0{
				deliver:=message[0]
				buf:=[]byte(deliver.text)
				for iter<len(connection){
					if deliver.author != connection[iter].author{
					connection[iter].conn.Write(buf)
				}
					iter++
				}
				message=message[1:]
		}
	}
}

func processGet(client clientConn){
	buf:=make([]byte,256)
	for {
		size,err:=bufio.NewReader(client.conn).Read(buf)
		if err!=nil{
      client.conn.Close()
      fmt.Println("[CLIENT HAS DISCONECTED]")
			return
		}else{
			message=append(message,messageData{author:client.author,text:string(buf[:size])})
		}
	}
}

func GetMessage(conn net.Conn,login string) {
	var message string
	buf:=make([]byte,256)
	for {
		size,err:=bufio.NewReader(conn).Read(buf)
		if err!=nil{
      conn.Close()
      panic("[YOU HAVE DISCONECTED]")
		}
		message = string(buf[:size])
		fmt.Println(message)
	}
}

func SendMessage(conn net.Conn,login string) {
	for {
		mess,err:=bufio.NewReader(os.Stdin).ReadString('\n')
		processErr(err)
		logNmess:=fmt.Sprintf("%s:%s",login,mess)
    buf:=[]byte(logNmess)
		conn.Write(buf)
	}
}

func startClient() {
	var login string
	buf:=make([]byte,256)
	serv_addr := fmt.Sprintf("%s:%d", *serverIP, *port)
  addr, err := net.ResolveTCPAddr("tcp", serv_addr)
  processErr(err)
	conn, err := net.DialTCP("tcp",nil, addr)
	size,err:=bufio.NewReader(conn).Read(buf)
	if err!=nil{
		conn.Close()
		panic("[YOU HAVE DISCONECTED]")
	}
	fmt.Println(string(buf[:size]))
	fmt.Scanln(&login)
	conn.Write([]byte(login))
	processErr(err)
  go SendMessage(conn,login)
  go GetMessage(conn,login)
}

func main() {
	start=true
	flag.Parse()
	defer recv()
	if *server{
		wg.Add(1)
		go startServer()
	}else{
	wg.Add(1)
	go startClient()
}
	wg.Wait()
}
