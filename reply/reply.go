package reply

import (
	"strconv"
)

//换行
var rn = []byte{13, 10}

//返回值接口，提供写回数据的能力
type Reply interface {
	Write() []byte
}

//以下所有代码，为实现redis返回值的协议，具体可以参考
//http://redis.cn/topics/protocol.html

type SimpleReply struct {
	Data string
}

func (reply *SimpleReply) Write() []byte {
	if reply.Data == "" {
		return append([]byte("+-1"), rn...)
	}
	return append([]byte("+"+reply.Data), rn...)
}

type IntegerReply struct {
	Data int
}

func (reply *IntegerReply) Write() []byte {
	return append([]byte(":" + strconv.Itoa(reply.Data)), rn...)
}

type BulkReply struct {
	Data []byte
}

func (reply *BulkReply) Write() []byte {
	if reply.Data == nil {
		return append([]byte("$-1"), rn...)
	}

	pre := []byte("$" + strconv.Itoa(len(reply.Data)))
	pre = append(pre, rn...)
	return append(append(pre, reply.Data...), rn...)
}

type ArrayReply struct {
	Data []string
}

func (reply *ArrayReply) Write() []byte {
	if reply.Data == nil {
		reply.Data = []string{}
	}
	pre := []byte("*" + strconv.Itoa(len(reply.Data)))
	var result []byte
	result = append(result, pre...)
	result = append(result, rn...)
	for _,v := range reply.Data {
		result = append(result, []byte("$" + strconv.Itoa(len([]byte(v))))...)
		result = append(result, rn...)
	}
	return result
}


type ErrorReply struct {
	Data string
}

func (reply *ErrorReply) Write() []byte {
	return append([]byte("-"+reply.Data), rn...)
}