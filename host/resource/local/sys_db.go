package local

import (
	"github.com/obgnail/plugin-platform/common/common_type"
)

var _ common_type.SysDB = (*SysDBOp)(nil)

type SysDBOp struct {
	//msg    *protocol.DatabaseMessage_DatabaseRequestMessage
	plugin common_type.IPlugin
}

func NewSysDB(plugin common_type.IPlugin) common_type.SysDB {
	return &SysDBOp{plugin: plugin}
}

//func (db *SysDBOp) buildMessage(databaseRequestMessage *protocol.DatabaseMessage_DatabaseRequestMessage) *protocol.PlatformMessage {
//	//msg := utils.GetInitMessage()
//	//msg.Resource = &protocol.ResourceMessage{
//	//	Database: &protocol.DatabaseMessage{
//	//		DBRequest: databaseRequestMessage,
//	//	},
//	//}
//	//return msg
//	return nil
//}
//
//func (db *SysDBOp) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
//	//return SyncSendToHost(db.plugin, platformMessage)
//	return nil, nil
//}
//
//func (db *SysDBOp) Common(dbName int, op string) (protocol.DatabaseMessage_DBInstanceType, protocol.DatabaseMessage_DBOperationType) {
//	//var instance protocol.DatabaseMessage_DBInstanceType
//	//switch dbName {
//	//case 0:
//	//	instance = protocol.DatabaseMessage_Project
//	//case 1:
//	//	instance = protocol.DatabaseMessage_Wiki
//	//case 2:
//	//	instance = protocol.DatabaseMessage_Local
//	//}
//	//
//	//var operation protocol.DatabaseMessage_DBOperationType
//	//switch strings.ToLower(op) {
//	//case "query":
//	//	operation = protocol.DatabaseMessage_Query
//	//case "insert":
//	//	operation = protocol.DatabaseMessage_Insert
//	//case "update":
//	//	operation = protocol.DatabaseMessage_Update
//	//case "delete":
//	//	operation = protocol.DatabaseMessage_Delete
//	//case "create":
//	//	operation = protocol.DatabaseMessage_Create
//	//case "alter":
//	//	operation = protocol.DatabaseMessage_Alter
//	//case "drop":
//	//	operation = protocol.DatabaseMessage_Drop
//	//}
//	//
//	//return instance, operation
//	return 0, 0
//}

func (d *SysDBOp) Select(db, sql string) ([]*common_type.RawData, []*common_type.ColumnDesc, common_type.PluginError) {
	//b := validateSysDbName(dbName)
	//if !b {
	//	return nil, nil, common_type.NewPluginError(common_type.DataBaseNameFailure, common_type.DataBaseNameError.Error(), "Database Name illegal")
	//}
	//instance, operation := db.Common(dbName, op)
	//databaseRequestMessage := &protocol.DatabaseMessage_DatabaseRequestMessage{
	//	Instance:  instance,
	//	Operation: operation,
	//	Statement: sql,
	//}
	//msg, err := db.sendMsgToHost(db.buildMessage(databaseRequestMessage))
	//if err != nil {
	//	return nil, nil, err
	//}
	//if msg.GetResource().GetDatabase().GetDBResponse().GetError() != nil {
	//	reterr := msg.GetResource().GetDatabase().GetDBResponse().GetError()
	//	return nil, nil, common_type.NewPluginError(int(reterr.GetCode()), reterr.GetError(), reterr.GetMsg())
	//}
	//if msg.GetResource().GetDatabase().GetDBResponse().GetDBError() != nil {
	//	dberr := msg.GetResource().GetDatabase().GetDBResponse().GetDBError()
	//	return nil, nil, common_type.NewPluginError(int(dberr.GetCode()), dberr.GetError(), dberr.GetMsg())
	//}
	//rawDatas, columns := resourceutils.ParseTableData(msg.GetResource().GetDatabase().GetDBResponse().GetData())
	//return rawDatas, columns, nil
	return nil, nil, nil
}

func (d *SysDBOp) AsyncSelect(db, sql string, callback common_type.DBCallBack) {
	//instance, operation := db.Common(dbName, op)
	//databaseRequestMessage := &protocol.DatabaseMessage_DatabaseRequestMessage{
	//	Instance:  instance,
	//	Operation: operation,
	//	Statement: sql,
	//}
	//asyncSysDb := new(AsyncSysDb)
	//asyncSysDb.callBackHandler = callback
	//AsyncSendToHost(db.plugin, db.buildMessage(databaseRequestMessage), asyncObj, asyncSysDb.callBack, nil)
	return
}

func (d *SysDBOp) Unmarshal(rawData []*common_type.RawData, columnDesc []*common_type.ColumnDesc, v interface{}) common_type.PluginError {
	//err := resourceutils.Unmarshal(rawData, columnDesc, v)
	//if err != nil {
	//	return common_type.NewPluginError(common_type.UnmarshalFailure, err.Error(), common_type.UnmarshalError.Error())
	//}
	return nil
}

func (d *SysDBOp) Exec(db, sql string) common_type.PluginError {
	return nil
}
