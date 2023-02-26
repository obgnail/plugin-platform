package errors

import "net/http"

const (
	AccessDenied                        = "AccessDenied"
	AlreadyExists                       = "AlreadyExists"
	AuthFailure                         = "AuthFailure"
	BadConfig                           = "BadConfig"
	Blocked                             = "Blocked"
	ConstraintViolation                 = "ConstraintViolation"
	CorruptedData                       = "CorruptedData"
	Gone                                = "Gone"
	InUse                               = "InUse"
	InvalidEnum                         = "InvalidEnum"
	InvalidParameter                    = "InvalidParameter"
	KeyConflict                         = "KeyConflict"
	LimitExceeded                       = "LimitExceeded"
	MissingParameter                    = "MissingParameter"
	NotFound                            = "NotFound"
	SourceDeleted                       = "SourceDeleted"
	OK                                  = "OK"
	Deleted                             = "Deleted"
	PermissionDenied                    = "PermissionDenied"
	RedisError                          = "RedisError"
	ServerError                         = "ServerError"
	SQLError                            = "SQLError"
	Timeout                             = "Timeout"
	TypeMismatch                        = "TypeMismatch"
	StatusError                         = "StatusError"
	UnexpectedArguments                 = "UnexpectedArguments"
	UnknownError                        = "UnknownError"
	VerificationFailure                 = "VerificationFailure"
	InvalidFileExt                      = "InvalidFileExt"
	PluginAlreadyRunning                = "PluginAlreadyRunning"
	PluginInstanceInstallationFailure   = "PluginInstanceInstallationFailure"
	PluginInstanceUninstallationFailure = "PluginInstanceUninstallationFailure"
	PluginInstanceUploadFailure         = "PluginInstanceUploadFailure"
	PluginInstanceEnableFailure         = "PluginInstanceEnableFailure"
	PluginInstanceDisableFailure        = "PluginInstanceDisableFailure"
	PluginInstanceInternalError         = "PluginInstanceInternalError"
)

var (
	DefaultStatusCodeBinding = map[string]int{
		AccessDenied:                        http.StatusForbidden,
		AlreadyExists:                       http.StatusConflict,
		StatusError:                         http.StatusConflict,
		AuthFailure:                         http.StatusUnauthorized,
		BadConfig:                           http.StatusInternalServerError,
		Blocked:                             http.StatusForbidden,
		ConstraintViolation:                 http.StatusForbidden,
		CorruptedData:                       http.StatusInternalServerError,
		Gone:                                http.StatusGone,
		InUse:                               http.StatusBadRequest,
		InvalidEnum:                         http.StatusInternalServerError,
		InvalidParameter:                    http.StatusBadRequest,
		KeyConflict:                         http.StatusInternalServerError,
		LimitExceeded:                       http.StatusForbidden,
		MissingParameter:                    http.StatusBadRequest,
		NotFound:                            http.StatusNotFound,
		Deleted:                             http.StatusInternalServerError,
		SourceDeleted:                       http.StatusInternalServerError,
		OK:                                  http.StatusOK,
		PermissionDenied:                    http.StatusForbidden,
		RedisError:                          http.StatusInternalServerError,
		ServerError:                         http.StatusInternalServerError,
		SQLError:                            http.StatusInternalServerError,
		Timeout:                             http.StatusBadRequest,
		TypeMismatch:                        http.StatusInternalServerError,
		UnexpectedArguments:                 http.StatusInternalServerError,
		UnknownError:                        http.StatusInternalServerError,
		VerificationFailure:                 http.StatusBadRequest,
		InvalidFileExt:                      http.StatusBadRequest,
		PluginAlreadyRunning:                http.StatusBadRequest,
		PluginInstanceInstallationFailure:   http.StatusBadRequest,
		PluginInstanceUninstallationFailure: http.StatusBadRequest,
		PluginInstanceUploadFailure:         http.StatusBadRequest,
		PluginInstanceEnableFailure:         http.StatusBadRequest,
		PluginInstanceDisableFailure:        http.StatusBadRequest,
		PluginInstanceInternalError:         http.StatusInternalServerError,
	}
)
