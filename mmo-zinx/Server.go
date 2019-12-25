package main

import (
	"fmt"
	"mmo-go/zinx/ziface"
	"mmo-go/zinx/znet-tcp"
)

//ping test 自定义路由
type PingRouter struct {
	znet_tcp.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//	先读取客户端的数据，然后再回写ping...ping...ping
	fmt.Println("recv from client msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet_tcp.BaseRouter
}

//Test Handle
func (this *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle")
	//	先读取客户端的数据，然后再回写ping...ping...ping
	fmt.Println("recv from client msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("hello"))
	if err != nil {
		fmt.Println(err)
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	//设置两个链接属性，在连接创建之后
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.jianshu.com/u/35261429b7f1")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}

	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	//	创建一个server句柄,使用zinx的api
	s := znet_tcp.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	// 启动server
	s.Serve()
}
