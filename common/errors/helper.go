package errors

import (
	"reflect"
)

// 生成 TypeMismatch 错误
func TypeMismatchError(value interface{}, expectedTypes ...string) error {
	var actualType = "nil"
	if value != nil {
		actualType = reflect.ValueOf(value).Type().String()
	}
	err := Errorf(TypeMismatch, "expected %v, got [%s](%v)", expectedTypes, actualType, value)
	err.(*Err).SetLocation(1)
	return err
}

func PluginUploadError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceUploadFailure, "", "", reason)
}

func PluginInstallError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceInstallationFailure, "", "", reason)
}

func PluginEnableError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceEnableFailure, "", "", reason)
}

func PluginDisableError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceDisableFailure, "", "", reason)
}

func PluginUninstallError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceUninstallationFailure, "", "", reason)
}

func PluginUpgradeError(reason string) error {
	return errorWithModelFieldReason(PluginInstanceUpgradeFailure, "", "", reason)
}

func PluginMessageError(reason string) error {
	return errorWithModelFieldReason(GetPluginMessageFailure, "", "", reason)
}

func MissingParameterError(model string, field string) error {
	return errorWithModelFieldReason(MissingParameter, model, field, "")
}

func errorWithModelFieldReason(t string, model string, field string, reason string) error {
	parts := []string{t}
	values := make(map[string]interface{})
	if len(model) > 0 {
		parts = append(parts, model)
		values[modelKey] = model
	}
	if len(field) > 0 {
		parts = append(parts, field)
		values[fieldKey] = field
	}
	if len(reason) > 0 {
		parts = append(parts, reason)
		values[reasonKey] = reason
	}
	err := New(parts...).(*Err)
	err.values = values
	err.SetLocation(2)
	return err
}
