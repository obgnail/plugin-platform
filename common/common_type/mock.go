package common_type

var _ IInstanceDescription = (*MockInstanceDesc)(nil)
var _ IPluginDescriptor = (*MockPluginDescriptor)(nil)

type MockInstanceDesc struct {
	PluginInstanceID string
	PluginDescriptor *MockPluginDescriptor
}

func (i *MockInstanceDesc) InstanceID() string { return i.PluginInstanceID }
func (i *MockInstanceDesc) PluginDescription() IPluginDescriptor {
	return i.PluginDescriptor
}

type MockPluginDescriptor struct {
	AppID      string
	PluginName string
	Lang       string
	LangVer    string
	AppVer     string
	HostVer    string
	MinSysVer  string
}

func (i *MockPluginDescriptor) ApplicationID() string { return i.AppID }
func (i *MockPluginDescriptor) Name() string          { return i.PluginName }
func (i *MockPluginDescriptor) Language() string      { return i.Lang }
func (i *MockPluginDescriptor) LanguageVersion() IVersion {
	result, _ := ParseVersionString(i.LangVer)
	return result
}
func (i *MockPluginDescriptor) ApplicationVersion() IVersion {
	result, _ := ParseVersionString(i.AppVer)
	return result
}
func (i *MockPluginDescriptor) HostVersion() IVersion {
	result, _ := ParseVersionString(i.HostVer)
	return result
}
func (i *MockPluginDescriptor) MinSystemVersion() IVersion {
	result, _ := ParseVersionString(i.MinSysVer)
	return result
}

type MockHost struct {
	Info   HostInfo
	Status HostStatus
}

func (h *MockHost) GetInfo() HostInfo {
	return h.Info
}

func (h *MockHost) GetStatus() HostStatus {
	return h.Status
}

//
//// Store 存在内存表
//func (h *MockHost) Store() {
//	return
//}

//
//func (h *MockHost) KillSelf() {
//	return
//}
//
//func (h *MockHost) KillPlugin(instanceUUID string) {
//	return
//}

type MockHostBoot struct {
	Info   HostBootInfo
	Status HostBootStatus
}

func (b *MockHostBoot) GetInfo() HostBootInfo {
	return b.Info
}

func (b *MockHostBoot) GetStatus() HostBootStatus {
	return b.Status
}
