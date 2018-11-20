package main

import (
	"fmt"
	"gedis/command"
	"gedis/protocol"
	"gedis/reply"
	"net"
	"os"
)

func main() {
	service := ":6379"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//启动协程处理redis请求
		go handleClient(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		allData := readAll(conn)
		//貌似这么写可以处理客户端连接断开的情况，不知道有没有坑- -
		if allData == nil {
			return
		}
		message := protocol.RedisMessage{}
		//将数据包封装成RedisMessage对象，方便使用
		message.ToMessage(allData)
		var re reply.Reply
		//拿到对应命令的处理器
		handler, err := command.GetHandlerManager().FindHandler(message.Command())
		if err != nil {
			re = &reply.ErrorReply{Data: err.Error()}
		} else {
			//处理请求，得到返回结果
			re = handler.Handler(&message)
		}
		//将返回结果写回
		conn.Write(re.Write())
	}
}

//读取所有的数据，不知道这么写有没有坑- -
func readAll(conn net.Conn) []byte {
	var result []byte
	var buf [2048]byte
	var n = 2048
	var err error
	for n == 2048 {
		n, err = conn.Read(buf[:])
		if err != nil {
			return nil
		}
		result = append(result, buf[0:n]...)
	}
	return result
}
