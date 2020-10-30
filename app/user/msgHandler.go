package user

import (
	"github.com/kudoochui/kudos/rpc"
	msgService "github.com/kudoochui/kudos/service/msgService"
	"github.com/kudoochui/kudosServer/app/user/msg"
	"github.com/kudoochui/kudosServer/app/user/remote"
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
	RegisterHandler(new(remote.Hello))
	RegisterHandler(new(remote.Hi))
	room := &remote.Room{
		RoomRemote: &remote.RoomRemote{},
	}
	room.RoomRemote.Room = room
	RegisterHandler(room.RoomRemote)
	RegisterHandler(&remote.User{Room: room})

	// register msg type
	msgService.GetMsgService().Register("Hello.Say", &remote.HelloReq{}, &remote.HelloResp{})
	msgService.GetMsgService().Register("Hi.Say", &msg.HiReq{}, &msg.HiResp{})
	msgService.GetMsgService().Register("RoomRemote.Join", &remote.RoomJoin{}, &remote.RoomResp{})
	msgService.GetMsgService().Register("RoomRemote.Leave", &remote.RoomLeave{}, &remote.RoomResp{})
	msgService.GetMsgService().RegisterPush("onNotify")
	msgService.GetMsgService().RegisterPush("onLeave")
	msgService.GetMsgService().RegisterPush("onJoin")
	msgService.GetMsgService().RegisterPush("onSay")
}
