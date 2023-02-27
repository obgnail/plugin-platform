package db

import (
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/message_utils"
	"github.com/obgnail/plugin-platform/common/protocol"
)

type DataBase struct {
	Source   *protocol.PlatformMessage
	Distinct *protocol.PlatformMessage

	db            string
	sql           string
	importSqlPath string

	instanceID string
	isLocalDB  bool
}

func NewDataBase(sourceMessage *protocol.PlatformMessage, distinctMessage *protocol.PlatformMessage) *DataBase {
	dataBaseReqMsg := sourceMessage.GetResource().GetDatabase().GetDBRequest()

	isLocalDB := message_utils.IsLocalDB(dataBaseReqMsg.GetDB())
	db := dataBaseReqMsg.GetDB()
	if isLocalDB {
		db = config.String("platform.mysql_user_plugin_db_name", "plugins")
	}

	dataBase := &DataBase{
		Source:        sourceMessage,
		Distinct:      distinctMessage,
		db:            db,
		importSqlPath: dataBaseReqMsg.GetSqlFileName(),
		sql:           dataBaseReqMsg.GetStatement(),
		instanceID:    sourceMessage.GetResource().GetSender().GetInstanceID(),
		isLocalDB:     isLocalDB,
	}
	return dataBase
}

func (d *DataBase) Execute() {
	if d.importSqlPath != "" {
		d.importSql()
	} else {
		d.onSql()
	}
}

func (d *DataBase) onSql() {
	var data *protocol.TableMessage
	var err common_type.PluginError

	defer d.buildMsg(data, err)

	db := GetDB(d.db)
	stmt, realSql, err := db.prepare(d.sql, d.instanceID, d.isLocalDB)
	if err != nil {
		return
	}

	if stmt == StmtSelect {
		data, err = db.Select(realSql)
	} else if stmt == StmtRow {
		err = db.Exec(realSql)
	}
}

func (d *DataBase) importSql() {
	version := message_utils.VersionPb2String(d.Source.GetResource().GetSender().GetApplication().GetApplicationVersion())
	appId := d.Source.GetResource().GetSender().GetApplication().GetApplicationID()
	db := GetDB(d.db)
	err := db.ImportSql(d.importSqlPath, appId, version, d.instanceID)
	d.buildMsg(nil, err)
}

func (d *DataBase) buildMsg(data *protocol.TableMessage, err common_type.PluginError) {
	resp := &protocol.DatabaseMessage_DatabaseResponseMessage{
		Data:    data,
		DBError: message_utils.BuildErrorMessage(err),
	}
	message_utils.BuildResourceDbMessage(d.Distinct, resp)
}
