package command

import (
	"errors"
	"fmt"
	"gedis/protocol"
	"gedis/reply"
	"sync"
)

//命令处理器管理工具，单例
type handlerManager struct {
	handlerMap map[string]RedisCommandHandler
}

func (manager *handlerManager) FindHandler(commandName string) (RedisCommandHandler, error) {
	handler := manager.handlerMap[commandName]
	if handler != nil {
		return handler, nil
	} else {
		return nil, errors.New("not support command: " + commandName)
	}
}

var manager *handlerManager
var once sync.Once

func GetHandlerManager() *handlerManager {
	once.Do(func() {
		manager = &handlerManager{}
		manager.handlerMap = map[string]RedisCommandHandler{
			//应用启动的时候进行命令处理器的初始化
			"GET": &getCommandHandler{},
			"SET": &setCommandHandler{},
			"DEL": &delCommandHandler{},
		}
	})
	return manager
}

//------------------------------------------分割线---------------------------------------

//模拟redis节点，用于存放数据- -
var dataMap = make(map[string]string)

//redis命令处理器接口，提供处理redis命令的能力
type RedisCommandHandler interface {
	Handler(msg *protocol.RedisMessage) reply.Reply
}

//get命令
type getCommandHandler struct {
}

func (handler *getCommandHandler) Handler(msg *protocol.RedisMessage) reply.Reply {
	key := string(msg.Param()[0])
	value := dataMap[key]
	if value == "" {
		return &reply.BulkReply{}
	} else {
		return &reply.BulkReply{Data: []byte(value)}
	}

}

//set命令
type setCommandHandler struct {
}

func (handler *setCommandHandler) Handler(msg *protocol.RedisMessage) reply.Reply {
	key := string(msg.Param()[0])
	value := string(msg.Param()[1])
	dataMap[key] = value
	return &reply.SimpleReply{Data: "OK"}
}

//del命令
type delCommandHandler struct {
}

func (handler *delCommandHandler) Handler(msg *protocol.RedisMessage) reply.Reply {
	key := string(msg.Param()[0])
	fmt.Println("key:" + key)
	if dataMap[key] == "" {
		return &reply.IntegerReply{Data: 0}
	} else {
		delete(dataMap, key)
		return &reply.IntegerReply{Data: 1}
	}
}
