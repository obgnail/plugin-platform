package common

type SettingLayoutCardBody struct {
	Action          int32  `json:"action"`
	IssueDetailType int32  `json:"issue_detail_type"`
	SubType         bool   `json:"sub_type"`
	TabType         string `json:"tab_label"`
	Position        int32  `json:"position"`
	CardType        string `json:"card_type"`
	CardLabel       string `json:"card_label"`
	ToPluginLabel   string `json:"to_plugin_label"`
}

type SettingLayoutCardReq struct {
	ReqBody []SettingLayoutCardBody `json:"records"`
}

type LayoutCard interface {
	Setting(*SettingLayoutCardReq) PluginError
}
