package common_type

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

func (pe *PError) Msg() string {
	return pe.msg
}

func (pe *PError) Code() int {
	return pe.code
}

func (pe *PError) Error() string {
	return pe.error
}

func NewPluginError(code int, msg string) PluginError {
	err := getErr(code)
	if err == nil {
		code = UnknownError
		err = m[UnknownError]
	}

	pe := &PError{
		code:  code,
		error: err.Error(),
		msg:   msg,
	}
	return pe
}
