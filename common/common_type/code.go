package common_type

import (
	"fmt"
	"net/http"
)

const (
	StatusOK                    = http.StatusOK
	MsgTimeOut                  = 11001
	SocketListenOrDialFailure   = 11002
	TargetEndpointNotFound      = 11003
	ProtoUnmarshalFailure       = 11004
	ProtoMarshalFailure         = 11005
	SysDbSelectFailure          = 11006
	SysDbExecFailure            = 11007
	SysDbCountFailure           = 11008
	CreateFileFailure           = 11009
	ReNameFileFailure           = 11010
	RemoveFileFailure           = 11011
	IsExistFileFailure          = 11012
	CopyFileFailure             = 11013
	ListFileFailure             = 11014
	IsDirFailure                = 11015
	ReadFailure                 = 11016
	ReadLinesFailure            = 11017
	WriteBytesFailure           = 11018
	AppendBytesFailure          = 11019
	WriteStringsFailure         = 11020
	AppendStringsFailure        = 11021
	ZipFailure                  = 11022
	UnZipFailure                = 11023
	GzFailure                   = 11024
	UnGzFailure                 = 11025
	HashFailure                 = 11026
	MakeDirFailure              = 11027
	OnstartFailure              = 11028
	OnstopFailure               = 11029
	OnEnableFailure             = 11030
	OnDisEnableFailure          = 11031
	OnUpgradeFailure            = 11032
	OnCheckCompatibilityFailure = 11033
	OnCheckStateFailure         = 11034
	OnHeartbeatFailure          = 11035
	OnPluginHttpFailure         = 11036
	AsyncFetchFailure           = 11037
	GetInstanceFailure          = 11038
	DbErrorFailure              = 11040
	FileNotFoundFailure         = 11041
	CallSysAPIFailure           = 11042
	OutgoingFailure             = 11043
	SysDbImportSqlFailure       = 11044
	DataBaseNameFailure         = 11045
	NotifyEventFailure          = 11046
	AddLayoutCardPluginFailure  = 11047
	UnmarshalFailure            = 11048
	ItemsAddFailure             = 11049
	FieldsAddFailure            = 11050
	AddGroupFieldFailure        = 11051
	UpdateFieldOptionFailure    = 11052
	OnInstallFailure            = 11053
	OnUnInstallFailure          = 11054
	LocalDevelopErr             = 11055
	EndpointReceiveErr          = 11056
	EndpointSendErr             = 11057
	EndpointIdentifyErr         = 11058
)

var (
	MsgTimeOutError                  = fmt.Errorf("timeout %d", MsgTimeOut)
	SocketListenOrDialFailureError   = fmt.Errorf("socketListenOrDialFailureError %d", SocketListenOrDialFailure)
	TargetEndpointNotFoundError      = fmt.Errorf("targetEndpointNotFoundError %d", TargetEndpointNotFound)
	ProtoUnmarshalFailureError       = fmt.Errorf("ProtoUnmarshalFailureError %d", ProtoUnmarshalFailure)
	ProtoMarshalFailureError         = fmt.Errorf("ProtoMarshalFailure %d", ProtoMarshalFailure)
	SysDbSelectFailureError          = fmt.Errorf("SysDbSelectFailure %d", SysDbSelectFailure)
	SysDbExecFailureError            = fmt.Errorf("SysDbExecFailure %d", SysDbExecFailure)
	SysDbCountFailureError           = fmt.Errorf("SysDbCountFailure %d", SysDbCountFailure)
	CreateFileFailureError           = fmt.Errorf("CreateFileFailure %d", CreateFileFailure)
	ReNameFileFailureError           = fmt.Errorf("ReNameFileFailure %d", ReNameFileFailure)
	RemoveFileFailureError           = fmt.Errorf("RemoveFileFailure %d", RemoveFileFailure)
	IsExistFileFailureError          = fmt.Errorf("IsExistFileFailure %d", IsExistFileFailure)
	CopyFileFailureError             = fmt.Errorf("CopyFileFailure %d", CopyFileFailure)
	ListFileFailureError             = fmt.Errorf("ListFileFailure %d", ListFileFailure)
	IsDirFailureError                = fmt.Errorf("IsDirFailure %d", IsDirFailure)
	ReadFailureError                 = fmt.Errorf("ReadFailure %d", ReadFailure)
	ReadLinesFailureError            = fmt.Errorf("ReadLinesFailure %d", ReadLinesFailure)
	WriteBytesFailureError           = fmt.Errorf("WriteBytesFailure %d", WriteBytesFailure)
	AppendBytesFailureError          = fmt.Errorf("AppendBytesFailure %d", AppendBytesFailure)
	WriteStringsFailureError         = fmt.Errorf("WriteStringsFailure %d", WriteStringsFailure)
	AppendStringsFailureError        = fmt.Errorf("AppendStringsFailure %d", AppendStringsFailure)
	ZipFailureError                  = fmt.Errorf("ZipFailure %d", ZipFailure)
	UnZipFailureError                = fmt.Errorf("UnZipFailure %d", UnZipFailure)
	GzFailureError                   = fmt.Errorf("GzFailure %d", GzFailure)
	UnGzFailureError                 = fmt.Errorf("UnGzFailure %d", UnGzFailure)
	HashFailureError                 = fmt.Errorf("HashFailure %d", HashFailure)
	MakeDirFailureError              = fmt.Errorf("MakeDirFailure %d", MakeDirFailure)
	AsyncFetchFailureError           = fmt.Errorf("AsyncFetchFailure %d", AsyncFetchFailure)
	GetInstanceFailureError          = fmt.Errorf("GetInstanceFailure %d", GetInstanceFailure)
	OnstartFailureError              = fmt.Errorf("OnstartFailure %d", OnstartFailure)
	OnstopFailureError               = fmt.Errorf("OnstopFailure %d", OnstopFailure)
	OnEnableFailureError             = fmt.Errorf("OnEnableFailure %d", OnEnableFailure)
	OnDisEnableFailureError          = fmt.Errorf("OnDisEnableFailure %d", OnDisEnableFailure)
	OnUpgradeFailureError            = fmt.Errorf("OnUpgradeFailure %d", OnUpgradeFailure)
	OnCheckCompatibilityFailureError = fmt.Errorf("OnCheckCompatibilityFailure %d", OnCheckCompatibilityFailure)
	OnCheckStateFailureError         = fmt.Errorf("OnCheckStateFailure %d", OnCheckStateFailure)
	OnHeartbeatFailureError          = fmt.Errorf("OnHeartbeatFailure %d", OnHeartbeatFailure)
	OnPluginHttpFailureError         = fmt.Errorf("OnPluginHttpFailure %d", OnPluginHttpFailure)
	DbError                          = fmt.Errorf("DbError %d", DbErrorFailure)
	FileNotFoundError                = fmt.Errorf("FileNotFoundError %d", FileNotFoundFailure)
	CallSysAPIError                  = fmt.Errorf("CallSysAPIFailure %d", CallSysAPIFailure)
	OutgoingError                    = fmt.Errorf("OutgoingFailure %d", OutgoingFailure)
	SysDbImportSqlError              = fmt.Errorf("SysDbImportSqlFailure %d", SysDbImportSqlFailure)
	DataBaseNameError                = fmt.Errorf("DataBaseNameFailure %d", DataBaseNameFailure)
	NotifyEventError                 = fmt.Errorf("NotifyEventFailure %d", NotifyEventFailure)
	AddLayoutCardPluginEventError    = fmt.Errorf("AddLayoutCardPluginEventFailure %d", AddLayoutCardPluginFailure)
	UnmarshalError                   = fmt.Errorf("UnmarshalFailure %d", UnmarshalFailure)
	ItemsAddError                    = fmt.Errorf("ItemsAddFailure %d", ItemsAddFailure)
	FieldsAddError                   = fmt.Errorf("FieldsAddFailure %d", FieldsAddFailure)
	AddGroupFieldError               = fmt.Errorf("AddGroupFieldFailure %d", AddGroupFieldFailure)
	UpdateFieldOptionError           = fmt.Errorf("UpdateFieldOptionFailure %d", UpdateFieldOptionFailure)
	OnInstallFailureError            = fmt.Errorf("OnInstallFailure %d", OnInstallFailure)
	OnUnInstallFailureError          = fmt.Errorf("OnUninstallFailure %d", OnUnInstallFailure)
	LocalDevelopError                = fmt.Errorf("LocalDevelopError %d", LocalDevelopErr)
	EndpointReceiveError             = fmt.Errorf("EndpointReceiveError %d", EndpointReceiveErr)
	EndpointSendError                = fmt.Errorf("EndpointSendError %d", EndpointSendErr)
	EndpointIdentifyError            = fmt.Errorf("EndpointIdentifyError %d", EndpointIdentifyErr)
)
