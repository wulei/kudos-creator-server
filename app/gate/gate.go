package gate

import (
	"errors"
	"fmt"
	"github.com/kudoochui/kudosServer/config"
	"github.com/kudoochui/kudos/app"
	"github.com/kudoochui/kudos/component/connector"
	rpcClient "github.com/kudoochui/kudos/component/proxy"
	rpcServer "github.com/kudoochui/kudos/component/remote"
	"github.com/kudoochui/kudos/component/timers"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
)

type Gate struct {
	*app.ServerDefault

	msgHandler *MsgHandler
}

func init()  {
	app.RegisterCreateServerFunc("gate", func(serverId string) app.Server {
		return &Gate{
			ServerDefault: app.NewServerDefault(serverId),
		}
	})
}

//没有登入过服务器是不让访问服务器的，可以设置拦截访问的路由，如果没有登入，就返回错误
func gateRoute(session *rpc.Session, servicePath, serviceName string) (string,error) {
	if (servicePath + "." + serviceName) != "User.Login" {
		if session.GetUserId() == 0 {
			return "", errors.New("not login")
		}
	}
	return servicePath, nil
}

func (g *Gate) OnStart(){
	// 开始加载组件，Gate必须加载的组件有connector,remote和proxy。connector负责与客户端连接通信，proxy负责与后端节点通信。在onStart中添加组件：
	settings, err := config.ServersConfig.GetMap("gate")
	if err != nil {
		log.Error("%s", err)
	}
	serverSetting := settings[g.ServerId].(map[string]interface{})

	// connector 与客户端连接通信
	wsAddr := fmt.Sprintf("%s:%.f",serverSetting["host"], serverSetting["clientPort"])
	conn := connector.NewConnector(
		connector.WSAddr(wsAddr),
		)
	g.Components["connector"] = conn
	//连接组件设置路由函数
	conn.Route(gateRoute)

	// remote
	remoteAddr := fmt.Sprintf("%s:%.f",serverSetting["host"], serverSetting["port"])
	remote := rpcServer.NewRemote(
		rpcServer.Addr(remoteAddr),
		rpcServer.RegistryType(config.RegistryConfig.String("registry")),
		rpcServer.RegistryAddr(config.RegistryConfig.String("addr")),
		rpcServer.BasePath(config.RegistryConfig.String("basePath")))
	g.Components["remote"] = remote
	g.msgHandler = &MsgHandler{r:remote}

	// remote 负责与后端节点通信
	proxy := rpcClient.NewProxy(
		rpcClient.RegistryType(config.RegistryConfig.String("registry")),
		rpcClient.RegistryAddr(config.RegistryConfig.String("addr")),
		rpcClient.BasePath(config.RegistryConfig.String("basePath")))
	g.Components["proxy"] = proxy

	timer := timers.NewTimer()
	g.Components["timer"] = timer

	for _,com := range g.Components {
		com.OnInit()
	}

	// register service.  Note: must behind remote OnInit
	g.msgHandler.RegisterHandler()

	//连接器收到的消息全部转发到代理组件
	conn.SetRouter(proxy)
	//注册连接器暴露的session服务和channel服务
	conn.SetRegisterServiceHandler(remote)
	//代理组件收到的回复发给连接器，由连接器返回给客户端
	proxy.SetRpcResponder(conn)
	//断线回调
	conn.SetConnectionListener(g)
}

func (g *Gate) Run(closeSig chan bool){
	for _,com := range g.Components {
		go com.Run(closeSig)
	}

	<-closeSig
	//closing
	log.Info("gate closing")
}

func (g *Gate) OnStop(){
	for _,com := range g.Components {
		com.OnDestroy()
	}
}

func (g *Gate) OnDisconnect(session *rpc.Session) {
	proxy := g.GetComponent("proxy").(*rpcClient.Proxy)
	args := &rpc.Args{
		Session: *session,
	}
	reply := &rpc.Reply{}
	//断线事件，转发给后端User节点
	proxy.RpcCall("User", "OnOffline", args, reply)
}