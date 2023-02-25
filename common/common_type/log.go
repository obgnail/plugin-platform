package common_type

type PluginLogger interface {
	Trace(string)
	Info(string)
	Warn(string)
	Error(string)
}

type VersionTime int

const (
	Earlier VersionTime = -1
	Same    VersionTime = 0
	Later   VersionTime = 1
)

func (vt *VersionTime) ToString() string {
	switch *vt {
	case Earlier:
		return "Earlier"
	case Same:
		return "Same"
	case Later:
		return "Later"
	default:
		return "Unknown"
	}
}
