package http_router

import (
	"sync"
)

var paramPool = sync.Pool{
	New: func() interface{} {
		return &Param{
			Keys:   make([]string, 0, 2),
			Values: make([]string, 0, 2),
		}
	},
}

type Param struct {
	Keys   []string
	Values []string
}

func acquireParam() *Param {
	p := paramPool.Get().(*Param)
	return p
}

func releaseParam(p *Param) {
	if p != nil && p.Keys != nil {
		p.Keys = p.Keys[:0]
		p.Values = p.Values[:0]
	}

	paramPool.Put(p)
}

func (param *Param) GetValue(key string) string {
	for i, v := range param.Keys {
		if v == key {
			return param.Values[i]
		}
	}
	return ""
}

func (param *Param) addKV(k, v string) int {
	param.Keys = append(param.Keys, k)
	param.Values = append(param.Values, v)
	return len(param.Keys)
}

func (param *Param) remove(index int) {
	param.Keys = param.Keys[:index]
	param.Values = param.Values[:index]
}
