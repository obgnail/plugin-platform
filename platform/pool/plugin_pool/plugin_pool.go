package plugin_pool

//
//import (
//	"github.com/BangWork/ones-plugin-common/golang/common"
//	"github.com/obgnail/plugin-platform/platform/model/mysql"
//	"github.com/obgnail/plugin-platform/common/errors"
//	"sync"
//	"time"
//)
//
//var (
//	Pool    *PluginPool
//	SendMsg func(message *protocol.PlatformMessage) common.PluginError
//)
//
//type PluginPool struct {
//	m             *PluginMap
//	Hostmap       *HostMap
//	HostPluginMap *HostPluginMap
//	TimeStamp     int64
//	sync.Mutex
//	pluginSecrets map[string]*PluginSecretKey
//}
//
//func NewPluginPool() *PluginPool {
//	return &PluginPool{
//		TimeStamp:     time.Now().Unix(),
//		Hostmap:       NewHostMap(),
//		HostPluginMap: NewHostPluginMap(),
//		pluginSecrets: make(map[string]*PluginSecretKey),
//		m:             NewPluginMap(),
//	}
//}
//
//func (p *PluginPool) StartPlugins() error {
//	m := mysql.ModelPluginInstance()
//	var instances = make([]*mysql.PluginInstance, 0)
//	var arg = new(mysql.PluginInstance)
//	if err := m.All(&instances, arg); err != nil {
//		return errors.Trace(err)
//	}
//
//	for _, instance := range instances {
//		cfg := instance.GetConfig()
//
//		proc := NewPluginProcess(cfg)
//		if instance.Status == PluginStatusRunning {
//			if err := proc.Enable(); err != nil {
//				continue
//			}
//		}
//		p.m.New(instance.InstanceUUID, proc)
//		go p.startPlugin(instance.InstanceUUID, cfg.Service)
//	}
//
//	p.RefreshTimeStamp()
//	return nil
//}
