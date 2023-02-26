package errors

import (
	jujuerrors "github.com/juju/errors"
	"net/http"
	"strings"
)

const (
	httpStatusKey  = "http_status"
	modelKey       = "model"
	fieldKey       = "field"
	reasonKey      = "reason"
	permissionKey  = "permission"
	formatKey      = "format"
	readableKey    = "readable"
	actionKey      = "action"
	descriptionKey = "description"
)

type Err struct {
	*jujuerrors.Err
	Code   string
	values map[string]interface{}
}

// eg. err := errors.New("InvalidParameter", "Task.Summary", "TooLong")
func New(codeparts ...string) error {
	code := mustCode(codeparts)
	jerr := jujuerrors.NewErr("<" + code + ">")
	jerr.SetLocation(1)
	return &Err{
		Code: code,
		Err:  &jerr,
	}
}

// 获取错误对应的 http 状态码
func (e *Err) HttpStatus() int {
	if e == nil {
		return http.StatusOK
	}
	status, ok := e.IntValue(httpStatusKey)
	if ok {
		return status
	}
	errtype := strings.Split(e.Code, ".")[0]
	status = DefaultStatusCodeBinding[errtype]
	if status == 0 {
		status = http.StatusInternalServerError
	}
	return status
}

// 设置自定义数据，设置的数据会返回给客户端
func (e *Err) SetValue(key string, value interface{}) {
	if e == nil {
		return
	}
	if e.values == nil {
		e.values = make(map[string]interface{})
	}
	e.values[key] = value
}

// 获取自定义数据
func (e *Err) Value(key string) (interface{}, bool) {
	if e == nil || e.values == nil {
		return nil, false
	}
	v, ok := e.values[key]
	return v, ok
}

// 获取所有自定义数据
func (e *Err) Values() map[string]interface{} {
	if e == nil {
		return nil
	}
	if e.values == nil {
		e.values = make(map[string]interface{})
	}
	return e.values
}

// 获取 int 类型的自定义数据
func (e *Err) IntValue(key string) (int, bool) {
	v, ok := e.Value(key)
	if !ok {
		return 0, false
	}
	iv, ok := v.(int)
	return iv, ok
}

// 获取 string 类型的自定义数据
func (e *Err) StringValue(key string) (string, bool) {
	v, ok := e.Value(key)
	if !ok {
		return "", false
	}
	sv, ok := v.(string)
	return sv, ok
}

// 额外指定 error 的自定义字段
func WithValue(err error, key string, value interface{}) error {
	_err, ok := err.(*Err)
	if !ok {
		var newErr error
		if err != nil {
			newErr = Errorf(UnknownError, err.Error())
		} else {
			newErr = New(UnknownError)
		}
		_err = Wrap(err, newErr).(*Err)
	}
	_err.SetValue(key, value)
	return _err
}

// 在现有错误的基础上包装一层错误
// eg:
// err := tx.Exec(sql, args...)
// if err != nil {
//     return errors.Wrap(err, errors.New(errors.SQLError))
// }
func Wrap(other error, err error) error {
	var code string
	var values map[string]interface{}
	var jerr *jujuerrors.Err
	if _err, ok := err.(*Err); ok {
		code = _err.Code
		values = _err.values
		jerr = jujuerrors.Wrap(other, _err.Err).(*jujuerrors.Err)
	} else {
		code = UnknownError
		prefixed := jujuerrors.NewErr(prefixWithCode(err.Error(), code))
		jerr = jujuerrors.Wrap(other, &prefixed).(*jujuerrors.Err)
	}
	jerr.SetLocation(1)
	return &Err{
		Code:   code,
		Err:    jerr,
		values: values,
	}
}

// 根据错误码生成一个错误，并加入自定义信息
func Errorf(code string, format string, args ...interface{}) error {
	jerr := jujuerrors.NewErr(prefixWithCode(format, code), args...)
	jerr.SetLocation(1)
	return &Err{
		Code: code,
		Err:  &jerr,
	}
}

func ErrorStack(err error) string {
	return jujuerrors.ErrorStack(err)
}

func prefixWithCode(s string, code string) string {
	return "<" + code + "> " + s
}

func mustCode(parts []string) string {
	if len(parts) == 0 {
		panic("error code parts is empty")
	}
	return strings.Join(parts, ".")
}

// 从里向外原样返回错误时，必须调用这个方法，以记录里层错误的栈信息
// ex:
// if err := SomeFunc(); err != nil {
//     return errors.Trace(err)
// }
func Trace(other error) error {
	if other == nil {
		return nil
	}
	newErr := new(Err)
	if err, ok := other.(*Err); ok {
		newErr.Code = err.Code
		newErr.values = err.values
	} else {
		newErr.Code = UnknownError
	}
	newErr.Err = jujuerrors.Trace(other).(*jujuerrors.Err)
	newErr.SetLocation(1)
	return newErr
}
