package common

type PluginError interface {
	Code() int
	Error() string
	Msg() string
}

type PError struct {
	msg   string
	code  int
	error string
}

func NewPluginError(code int, error string, msg string) PluginError {
	pe := &PError{
		msg:   msg,
		code:  code,
		error: error,
	}
	var e PluginError
	e = pe
	return e
}

func (pe *PError) Msg() string {
	return pe.msg
}

func (pe *PError) Code() int {
	return pe.code
}

func (pe *PError) Error() string {
	return pe.error
}
