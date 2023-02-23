package common_type

type PluginLogger interface {
	Info(string)
	Debug(string)
	Warn(string)
	Error(string)
	ErrorTrace(string)
	Fatal(string)
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
