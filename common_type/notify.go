package common

type Notify interface {
	SendMessage(*SendSubMessageReq) PluginError
}

type SendSubMessageReq struct {
	Title       string
	ToUsers     []string
	NotifyWay   string
	MessageBody []Body
	Ext         string
	Source      string
}

type Body struct {
	Body interface{}
	Url  string
}
