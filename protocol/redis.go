package protocol

import (
	"bytes"
	"strconv"
	"strings"
)

//解析redis请求数据包，具体解析的协议可以参考
//http://redis.cn/topics/protocol.html
func (msg *RedisMessage)ToMessage(b []byte) {
	//fmt.Println(string(b))
	buffer := bytes.NewBuffer(b)
	//暂时先跳过前两位
	first,_ := buffer.ReadByte()
	if first == 42 {
		buffer.Next(1)
		buffer.Next(2)
		buffer.Next(1)
		commandLength := buffer.Next(1)
		i, e := strconv.Atoi(string(commandLength))
		if e == nil {
			buffer.Next(2)
			//拿到命令名
			command := string(buffer.Next(i))
			msg.command = strings.ToUpper(command)
		}
		//拿到参数列表
		msg.param = readParam(buffer)
		//for _, v := range msg.param {
		//	fmt.Println("参数列表：", string(v))
		//}
	}
}

//解析出请求包里的各种参数
func readParam(buffer *bytes.Buffer) [][]byte {
	var result [][]byte
	for {
		buffer.Next(2)
		buffer.Next(1)
		i, e := strconv.Atoi(string(buffer.Next(1)))
		if e == nil {
			buffer.Next(2)
			param := buffer.Next(i)
			result = append(result, param)
		} else {
			break
		}
	}
	return result
}

//redis请求包对象
type RedisMessage struct {
	command string
	param [][]byte
}

func (msg *RedisMessage) Command() string {
	return msg.command
}

func (msg *RedisMessage) Param() [][]byte {
	return msg.param
}

