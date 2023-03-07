package common_type

import (
	"fmt"
	"net/http"
)

const (
	StatusOK                    = http.StatusOK
	UnknownError                = 11000 // 任何时候都不应该主动使用
	MsgTimeOut                  = 11001
	SocketListenOrDialFailure   = 11002
	TargetEndpointNotFound      = 11003
	ProtoUnmarshalFailure       = 11004
	ProtoMarshalFailure         = 11005
	DbSelectFailure             = 11006
	DbExecFailure               = 11007
	DbSqlSyntaxErr              = 11008
	CreateFileFailure           = 11009
	RenameFileFailure           = 11010
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
	CallPluginHttpFailure       = 11028
	CallPluginConfigFailure     = 11029
	OnEnableFailure             = 11030
	OnDisEnableFailure          = 11031
	OnUpgradeFailure            = 11032
	OnCheckCompatibilityFailure = 11033
	OnCheckStateFailure         = 11034
	OnHeartbeatFailure          = 11035
	OnPluginHttpFailure         = 11036
	AsyncFetchFailure           = 11037
	GetInstanceFailure          = 11038
	CallMainSystemAPIFailure    = 11039
	DbErrorFailure              = 11040
	FileNotFoundFailure         = 11041
	CallSysAPIFailure           = 11042
	OutgoingFailure             = 11043
	SysDbImportSqlFailure       = 11044
	DataBaseNameFailure         = 11045
	NotifyEventFailure          = 11046
	CallAbilityFailure          = 11047
	UnmarshalFailure            = 11048
	AddGroupFieldFailure        = 11051
	UpdateFieldOptionFailure    = 11052
	OnInstallFailure            = 11053
	OnUnInstallFailure          = 11054
	LocalDevelopErr             = 11055
	EndpointReceiveErr          = 11056
	EndpointSendErr             = 11057
	EndpointIdentifyErr         = 11058
)

var m = map[int]error{
	UnknownError:                fmt.Errorf("UnknownError"),
	MsgTimeOut:                  fmt.Errorf("timeout %d", MsgTimeOut),
	SocketListenOrDialFailure:   fmt.Errorf("socketListenOrDialFailureError %d", SocketListenOrDialFailure),
	TargetEndpointNotFound:      fmt.Errorf("targetEndpointNotFoundError %d", TargetEndpointNotFound),
	ProtoUnmarshalFailure:       fmt.Errorf("ProtoUnmarshalFailureError %d", ProtoUnmarshalFailure),
	ProtoMarshalFailure:         fmt.Errorf("ProtoMarshalFailure %d", ProtoMarshalFailure),
	DbSelectFailure:             fmt.Errorf("DbSelectFailure %d", DbSelectFailure),
	DbExecFailure:               fmt.Errorf("DbExecFailure %d", DbExecFailure),
	DbSqlSyntaxErr:              fmt.Errorf("DbSqlSyntaxError %d", DbSqlSyntaxErr),
	CreateFileFailure:           fmt.Errorf("CreateFileFailure %d", CreateFileFailure),
	RenameFileFailure:           fmt.Errorf("RenameFileFailure %d", RenameFileFailure),
	RemoveFileFailure:           fmt.Errorf("RemoveFileFailure %d", RemoveFileFailure),
	IsExistFileFailure:          fmt.Errorf("IsExistFileFailure %d", IsExistFileFailure),
	CopyFileFailure:             fmt.Errorf("CopyFileFailure %d", CopyFileFailure),
	ListFileFailure:             fmt.Errorf("ListFileFailure %d", ListFileFailure),
	IsDirFailure:                fmt.Errorf("IsDirFailure %d", IsDirFailure),
	ReadFailure:                 fmt.Errorf("ReadFailure %d", ReadFailure),
	ReadLinesFailure:            fmt.Errorf("ReadLinesFailure %d", ReadLinesFailure),
	WriteBytesFailure:           fmt.Errorf("WriteBytesFailure %d", WriteBytesFailure),
	AppendBytesFailure:          fmt.Errorf("AppendBytesFailure %d", AppendBytesFailure),
	WriteStringsFailure:         fmt.Errorf("WriteStringsFailure %d", WriteStringsFailure),
	AppendStringsFailure:        fmt.Errorf("AppendStringsFailure %d", AppendStringsFailure),
	ZipFailure:                  fmt.Errorf("ZipFailure %d", ZipFailure),
	UnZipFailure:                fmt.Errorf("UnZipFailure %d", UnZipFailure),
	GzFailure:                   fmt.Errorf("GzFailure %d", GzFailure),
	UnGzFailure:                 fmt.Errorf("UnGzFailure %d", UnGzFailure),
	HashFailure:                 fmt.Errorf("HashFailure %d", HashFailure),
	MakeDirFailure:              fmt.Errorf("MakeDirFailure %d", MakeDirFailure),
	CallPluginHttpFailure:       fmt.Errorf("CallPluginHttpFailure %d", CallPluginHttpFailure),
	CallPluginConfigFailure:     fmt.Errorf("CallPluginConfigFailure %d", CallPluginConfigFailure),
	OnEnableFailure:             fmt.Errorf("OnEnableFailure %d", OnEnableFailure),
	OnDisEnableFailure:          fmt.Errorf("OnDisEnableFailure %d", OnDisEnableFailure),
	OnUpgradeFailure:            fmt.Errorf("OnUpgradeFailure %d", OnUpgradeFailure),
	OnCheckCompatibilityFailure: fmt.Errorf("OnCheckCompatibilityFailure %d", OnCheckCompatibilityFailure),
	OnCheckStateFailure:         fmt.Errorf("OnCheckStateFailure %d", OnCheckStateFailure),
	OnHeartbeatFailure:          fmt.Errorf("OnHeartbeatFailure %d", OnHeartbeatFailure),
	OnPluginHttpFailure:         fmt.Errorf("OnPluginHttpFailure %d", OnPluginHttpFailure),
	AsyncFetchFailure:           fmt.Errorf("AsyncFetchFailure %d", AsyncFetchFailure),
	GetInstanceFailure:          fmt.Errorf("GetInstanceFailure %d", GetInstanceFailure),
	DbErrorFailure:              fmt.Errorf("DbError %d", DbErrorFailure),
	FileNotFoundFailure:         fmt.Errorf("FileNotFoundError %d", FileNotFoundFailure),
	CallSysAPIFailure:           fmt.Errorf("CallSysAPIFailure %d", CallSysAPIFailure),
	OutgoingFailure:             fmt.Errorf("OutgoingFailure %d", OutgoingFailure),
	SysDbImportSqlFailure:       fmt.Errorf("SysDbImportSqlFailure %d", SysDbImportSqlFailure),
	DataBaseNameFailure:         fmt.Errorf("DataBaseNameFailure %d", DataBaseNameFailure),
	NotifyEventFailure:          fmt.Errorf("NotifyEventFailure %d", NotifyEventFailure),
	CallAbilityFailure:          fmt.Errorf("CallAbilityFailure %d", CallAbilityFailure),
	UnmarshalFailure:            fmt.Errorf("UnmarshalFailure %d", UnmarshalFailure),
	AddGroupFieldFailure:        fmt.Errorf("AddGroupFieldFailure %d", AddGroupFieldFailure),
	UpdateFieldOptionFailure:    fmt.Errorf("UpdateFieldOptionFailure %d", UpdateFieldOptionFailure),
	OnInstallFailure:            fmt.Errorf("OnInstallFailure %d", OnInstallFailure),
	OnUnInstallFailure:          fmt.Errorf("OnUninstallFailure %d", OnUnInstallFailure),
	LocalDevelopErr:             fmt.Errorf("LocalDevelopError %d", LocalDevelopErr),
	EndpointReceiveErr:          fmt.Errorf("EndpointReceiveError %d", EndpointReceiveErr),
	EndpointSendErr:             fmt.Errorf("EndpointSendError %d", EndpointSendErr),
	EndpointIdentifyErr:         fmt.Errorf("EndpointIdentifyError %d", EndpointIdentifyErr),
}

func getErr(code int) error {
	err, ok := m[code]
	if ok {
		return err
	}
	return nil
}
