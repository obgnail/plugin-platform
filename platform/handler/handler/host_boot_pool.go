package handler

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"sync"
)

type HostBootPool struct {
	alive sync.Map //  在业务上运行的host_boot map[hostBootID]common_type.IHostBoot
}

func NewHostBootPool() *HostBootPool {
	return &HostBootPool{}
}

func (pool *HostBootPool) One() common_type.IHostBoot {
	var ret common_type.IHostBoot
	pool.Range(func(hostBootID string, hostBoot common_type.IHostBoot) bool {
		if hostBoot != nil {
			ret = hostBoot
			return false
		}
		return true
	})
	return ret
}

func (pool *HostBootPool) Add(host common_type.IHostBoot) {
	id := host.GetInfo().ID
	pool.alive.Store(id, host)
}

func (pool *HostBootPool) Delete(host common_type.IHostBoot) {
	id := host.GetInfo().ID
	pool.alive.Delete(id)
}

func (pool *HostBootPool) Exist(hostBootID string) bool {
	if _, ok := pool.alive.Load(hostBootID); ok {
		return true
	}
	return false
}

func (pool *HostBootPool) GetAll() []common_type.IHostBoot {
	var result []common_type.IHostBoot
	pool.alive.Range(func(key, value any) bool {
		result = append(result, value.(common_type.IHostBoot))
		return true
	})
	return result
}

func (pool *HostBootPool) Range(f func(hostBootID string, hostBoot common_type.IHostBoot) bool) {
	pool.alive.Range(func(key, value any) bool {
		return f(key.(string), value.(common_type.IHostBoot))
	})
}

func NewHostBoot(msg *protocol.ControlMessage_HostBootReportMessage, status common_type.HostBootStatus) *common_type.MockHostBoot {
	hostBoot := msg.GetBoot()

	b := &common_type.MockHostBoot{
		Info: common_type.HostBootInfo{
			ID:      hostBoot.BootID,
			Name:    hostBoot.Name,
			Version: message_utils.VersionPb2String(hostBoot.BootVersion),
		},
		Status: status,
	}
	return b
}
