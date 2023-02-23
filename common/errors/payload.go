package errors

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrPayload struct {
	Code       string
	Desc       string
	HttpStatus int
	Values     map[string]interface{}
}

func (p *ErrPayload) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if p.Values != nil {
		for k, v := range p.Values {
			if k != httpStatusKey {
				m[k] = v
			}
		}
	}
	m["errcode"] = p.Code
	m["code"] = p.HttpStatus
	if len(p.Desc) > 0 {
		m["desc"] = p.Desc
	}
	m["type"] = strings.Split(p.Code, ".")[0]
	return json.Marshal(m)
}

func (p *ErrPayload) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	p.Code, _ = m["errcode"].(string)
	delete(m, "errcode")
	p.Desc, _ = m["desc"].(string)
	delete(m, "desc")
	p.HttpStatus, _ = m["code"].(int)
	delete(m, "code")
	p.Values = m
	return nil
}

func NewErrPayload(err error) *ErrPayload {
	var code, desc string
	var httpStatus int
	var values map[string]interface{}
	switch err.(type) {
	case *Err:
		e := err.(*Err)
		if e == nil {
			code = OK
			desc = ""
			httpStatus = http.StatusOK
		} else {
			code = e.Code
			desc = e.Error()
			if desc == "<"+code+">" {
				desc = ""
			}
			httpStatus = e.HttpStatus()
			values = e.values
		}

	case nil:
		code = OK
		desc = ""
		httpStatus = http.StatusOK

	default:
		code = TypeMismatch
		desc = TypeMismatchError(err, "errors.*Err", "nil").Error()
		httpStatus = http.StatusInternalServerError
	}

	return &ErrPayload{
		Code:       code,
		Desc:       desc,
		HttpStatus: httpStatus,
		Values:     values,
	}
}
