package release

import (
	"database/sql"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/host/resource/common"
)

var _ common_type.LocalDB = (*LocalDB)(nil)
var _ common_type.SysDB = (*SysDB)(nil)

type CommonDB struct {
	db     *sql.DB
	plugin common_type.IPlugin
	sender common.Sender
}

func NewCommonDB(plugin common_type.IPlugin, sender common.Sender) *CommonDB {
	return &CommonDB{plugin: plugin, sender: sender}
}

func (d *CommonDB) buildMessage(databaseRequestMessage *protocol.DatabaseMessage_DatabaseRequestMessage) *protocol.PlatformMessage {
	msg := message_utils.GetInitMessage(nil, nil)
	msg.Resource = &protocol.ResourceMessage{
		Database: &protocol.DatabaseMessage{DBRequest: databaseRequestMessage},
	}
	return msg
}

func (d *CommonDB) sendMsgToHost(platformMessage *protocol.PlatformMessage) (*protocol.PlatformMessage, common_type.PluginError) {
	return d.sender.Send(d.plugin, platformMessage)
}

func (d *CommonDB) sendToHostAsync(platformMessage *protocol.PlatformMessage, callback common_type.DBCallBack) {
	cb := &dbCallbackWrapper{Func: callback}
	d.sender.SendAsync(d.plugin, platformMessage, cb.callBack)
}

func (d *CommonDB) CommonSelect(db, sql string) ([]*common_type.RawData, []*common_type.ColumnDesc, common_type.PluginError) {
	dbMsg := &protocol.DatabaseMessage_DatabaseRequestMessage{DB: db, Statement: sql}
	msg, err := d.sendMsgToHost(d.buildMessage(dbMsg))
	if err != nil {
		return nil, nil, err
	}
	if e := d.checkError(msg); e != nil {
		return nil, nil, e
	}
	rawDataset, columns := common.ParseTableData(msg.GetResource().GetDatabase().GetDBResponse().GetData())
	return rawDataset, columns, nil
}

func (d *CommonDB) CommonAsyncSelect(db, sql string, callback common_type.DBCallBack) {
	dbMsg := &protocol.DatabaseMessage_DatabaseRequestMessage{DB: db, Statement: sql}
	d.sendToHostAsync(d.buildMessage(dbMsg), callback)
}

func (d *CommonDB) CommonExec(db, sql string) common_type.PluginError {
	dbMsg := &protocol.DatabaseMessage_DatabaseRequestMessage{DB: db, Statement: sql}
	msg, err := d.sendMsgToHost(d.buildMessage(dbMsg))
	if err != nil {
		return err
	}
	if e := d.checkError(msg); e != nil {
		return e
	}
	return nil
}

func (d *CommonDB) Unmarshal(rawData []*common_type.RawData, columnDesc []*common_type.ColumnDesc, v interface{}) common_type.PluginError {
	if len(rawData) == 0 {
		return nil
	}
	if err := common.Unmarshal(rawData, columnDesc, v); err != nil {
		return common_type.NewPluginError(common_type.UnmarshalFailure, err.Error(), common_type.UnmarshalError.Error())
	}
	return nil
}

func (d *CommonDB) checkError(msg *protocol.PlatformMessage) common_type.PluginError {
	err := msg.GetResource().GetDatabase().GetDBResponse().GetError()
	if err != nil {
		return common_type.NewPluginError(int(err.GetCode()), err.GetError(), err.GetMsg())
	}
	err = msg.GetResource().GetDatabase().GetDBResponse().GetDBError()
	if err != nil {
		return common_type.NewPluginError(int(err.GetCode()), err.GetError(), err.GetMsg())
	}
	return nil
}

type dbCallbackWrapper struct {
	Func common_type.DBCallBack
}

func (w *dbCallbackWrapper) callBack(input, result *protocol.PlatformMessage, err common_type.PluginError) {
	dbResponse := result.GetResource().GetDatabase().GetDBResponse()
	data := dbResponse.GetData()
	rawDataset, columns := common.ParseTableData(data)
	w.Func(rawDataset, columns, err)
}

type LocalDB struct {
	db     string
	common *CommonDB
}

func NewLocalDB(plugin common_type.IPlugin, sender common.Sender) common_type.LocalDB {
	db := config.String("platform.mysql_user_plugin_db_name", "plugins")
	return &LocalDB{common: NewCommonDB(plugin, sender), db: db}
}

func (d *LocalDB) Select(sql string) ([]*common_type.RawData, []*common_type.ColumnDesc, common_type.PluginError) {
	return d.common.CommonSelect(d.db, sql)
}

func (d *LocalDB) AsyncSelect(sql string, callback common_type.DBCallBack) {
	d.common.CommonAsyncSelect(d.db, sql, callback)
}

func (d *LocalDB) Exec(sql string) common_type.PluginError {
	return d.common.CommonExec(d.db, sql)
}

func (d *LocalDB) ImportSQL(sqlFilePath string) common_type.PluginError {
	dbMsg := &protocol.DatabaseMessage_DatabaseRequestMessage{SqlFileName: sqlFilePath}
	msg, err := d.common.sendMsgToHost(d.common.buildMessage(dbMsg))
	if err != nil {
		return err
	}
	if e := d.common.checkError(msg); e != nil {
		return e
	}

	return nil
}

func (d *LocalDB) Unmarshal(rawData []*common_type.RawData, columnDesc []*common_type.ColumnDesc, v interface{}) common_type.PluginError {
	return d.common.Unmarshal(rawData, columnDesc, v)
}

type SysDB struct {
	common *CommonDB
}

func NewSysDB(plugin common_type.IPlugin, sender common.Sender) *SysDB {
	return &SysDB{common: NewCommonDB(plugin, sender)}
}

func (d *SysDB) Select(db, sql string) ([]*common_type.RawData, []*common_type.ColumnDesc, common_type.PluginError) {
	return d.common.CommonSelect(db, sql)
}

func (d *SysDB) AsyncSelect(db, sql string, callback common_type.DBCallBack) {
	d.common.CommonAsyncSelect(db, sql, callback)
}

func (d *SysDB) Exec(db, sql string) common_type.PluginError {
	return d.common.CommonExec(db, sql)
}

func (d *SysDB) Unmarshal(rawData []*common_type.RawData, columnDesc []*common_type.ColumnDesc, v interface{}) common_type.PluginError {
	return d.common.Unmarshal(rawData, columnDesc, v)
}
