package user

import (
	"fmt"
	"github.com/kudoochui/kudosServer/config"
	"github.com/kudoochui/kudos/app"
	rpcClient "github.com/kudoochui/kudos/component/proxy"
	rpcServer "github.com/kudoochui/kudos/component/remote"
	"github.com/kudoochui/kudos/log"
)

type UserServer struct {
	*app.ServerDefault

	msgHandler *MsgHandler
}

func init()  {
	app.RegisterCreateServerFunc("user", func(serverId string) app.Server {
		return &UserServer{
			ServerDefault: app.NewServerDefault(serverId),
		}
	})
}

func (g *UserServer) OnStart(){
	//开始加载组件，如果这个节点只提供服务，不需要调用其它节点的服务，只需要添加remote组件，在onStart中添加组件
	settings, err := config.ServersConfig.GetMap("user")
	if err != nil {
		log.Error("%s", err)
	}
	serverSetting := settings[g.ServerId].(map[string]interface{})
	remoteAddr := fmt.Sprintf("%s:%.f",serverSetting["host"], serverSetting["port"])

	remote := rpcServer.NewRemote(
		rpcServer.Addr(remoteAddr),
		rpcServer.RegistryType(config.RegistryConfig.String("registry")),
		rpcServer.RegistryAddr(config.RegistryConfig.String("addr")),
		rpcServer.BasePath(config.RegistryConfig.String("basePath")))
	g.Components["remote"] = remote
	g.msgHandler = &MsgHandler{r:remote}

	// 需要调用别的服务才需要proxy
	proxy := rpcClient.NewProxy(
		rpcClient.RegistryType(config.RegistryConfig.String("registry")),
		rpcClient.RegistryAddr(config.RegistryConfig.String("addr")),
		rpcClient.BasePath(config.RegistryConfig.String("basePath")))
	g.Components["proxy"] = proxy

	for _,com := range g.Components {
		com.OnInit()
	}

	// register service
	g.msgHandler.RegisterHandler()
}

func (g *UserServer) Run(closeSig chan bool){
	for _,com := range g.Components {
		go com.Run(closeSig)
	}
	<-closeSig
	//closing
	log.Info("user closing")
}

func (g *UserServer) OnStop(){
	for _,com := range g.Components {
		com.OnDestroy()
	}
}