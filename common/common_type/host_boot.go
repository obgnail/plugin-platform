package common_type

type HostBootStatus int

const (
	HostBootStatusNormal HostBootStatus = 1
	HostBootStatusDrift  HostBootStatus = 2
)

type IHostBoot interface {
	GetInfo() HostBootInfo
	GetStatus() HostBootStatus
}

type HostBootInfo struct {
	ID      string
	Name    string
	Version string
}
