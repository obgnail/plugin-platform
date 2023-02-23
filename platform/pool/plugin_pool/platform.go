package plugin_pool

//
//import (
//	"fmt"
//	"github.com/BangWork/ones-platform-api-old/protocol"
//	"github.com/obgnail/plugin-platform/platform/config"
//	"github.com/obgnail/plugin-platform/platform/service/utils"
//	"time"
//)
//
//var platformHandler *FurtherHandler
//
//type FurtherHandler struct {
//	innerEp  connection.Endpoint
//	SendBack bool
//	sync     *connection.Synchronized
//	log      *resourceLog.Log
//}
//
//func (ph *FurtherHandler) doHb() {
//	for {
//		select {
//		case <-time.After(time.Second * 10):
//			{
//				ph.sendHb()
//			}
//		}
//	}
//}
//
//func (ph *FurtherHandler) sendHb() {
//	heartbeatID := utils.CreateCaptcha()
//	controlMessage := protocol.BuildControlMessage()
//	controlMessage.Heartbeat = heartbeatID
//
//	addr := fmt.Sprintf("%s:%d", config.StringOrPanic("host"), config.IntOrPanic("tcp_port"))
//	sourceId, sourceTags := protocol.GetPlatformIDAndTags(config.Config.Version, addr)
//	hostmaps := Pool.Hostmap.All()
//	for _, v := range hostmaps {
//		hostID := v.HostDescriptor.GetHostID()
//		hostVersion := v.HostDescriptor.GetHostVersion().String()
//		hostName := v.HostDescriptor.GetName()
//		distinctId, distinctTags := protocol.GetHostIDAndTags(hostVersion, hostID, hostName)
//		seqNo := utilscommon.CreateCaptcha()
//
//		header := buildHeader(sourceId, distinctId, sourceTags, distinctTags, seqNo)
//		platformMessage := protocol.BuildPlatformMessage()
//		platformMessage.Header = header
//		platformMessage.Control = controlMessage
//		err := ph.SendMessage(platformMessage)
//		if err != nil {
//			ph.log.Error("ph.SendMessage distinctId: " + distinctId + " error: " + err.Error())
//		}
//	}
//}
//
//func RunPlatform() *FurtherHandler {
//	l := resourceLog.Logger
//	p := &connection.PBProto{Log: l}
//	version := config.Config.Version
//
//	addr := fmt.Sprintf("%s:%d", config.StringOrPanic("host"), config.IntOrPanic("tcp_port"))
//	router := connection.NewZmqPlatformEndpoint("R000001", connection.Platform, connection.RouterDumper, version, addr, connection.SocketTypeRouter,
//		l, p, true, true, true)
//	h := &FurtherHandler{
//		sync: connection.NewSynchronized(time.Duration(1000*config.Config.TimeoutSec), l),
//		log:  l,
//	}
//	h.Assign(router)
//	router.StartTimer(time.Second * 3)
//	platformHandler = h
//	go router.Connect()
//	go h.doHb()
//	return h
//}
