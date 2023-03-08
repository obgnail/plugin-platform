package common_type

type AbilityRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Func string `json:"func"`
	Args []byte `json:"args"`
}

type AbilityResponse struct {
	Data []byte      `json:"data"`
	Err  PluginError `json:"err"`
}
