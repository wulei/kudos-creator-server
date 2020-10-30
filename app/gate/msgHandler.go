package gate

import (
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/msgService"
)

// register server service to remote
var msgArray = []interface{}{}

func RegisterHandler(msg interface{}){
	msgArray = append(msgArray, msg)
}

type MsgHandler struct {
	r rpc.HandlerRegister
}

func (m *MsgHandler)RegisterHandler()  {
	for _,v := range msgArray {
		m.r.RegisterHandler(v,"")
	}
}

func init() {
	// 注册服务，它会导出Arith上所有方法
	RegisterHandler(new(Arith))

	// 消息消息路由，将导入导出参数和路由对应起来
	msgService.GetMsgService().Register("Arith.Mul", &Args{}, &Reply{})
}