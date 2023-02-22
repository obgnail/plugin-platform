package local

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	common "github.com/obgnail/plugin-platform/common_type"
	"github.com/obgnail/plugin-platform/host/config"
	"github.com/obgnail/plugin-platform/host/resource/utils"
	"io/ioutil"
	"math"
	"reflect"
	"regexp"
	"strings"
)

var _ common.LocalDB = (*LocalDB)(nil)

type LocalDB struct {
	db     *sql.DB
	plugin common.IPlugin
}

func NewLocalDB(plugin common.IPlugin) common.LocalDB {
	return &LocalDB{plugin: plugin}
}

func (localdb *LocalDB) initMysql() common.PluginError {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Asia%%2FShanghai&charset=utf8mb4&multiStatements=true",
		config.StringOrPanic("platform_mysql_user"),
		config.StringOrPanic("platform_mysql_password"),
		config.StringOrPanic("platform_mysql_host"),
		config.IntOrPanic("platform_mysql_port"),
		config.StringOrPanic("platform_mysql_db"),
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return common.NewPluginError(common.DataBaseNameFailure, err.Error(), err.Error())
	}
	if err = db.Ping(); err != nil {
		return common.NewPluginError(common.DataBaseNameFailure, err.Error(), err.Error())
	}
	localdb.db = db
	return nil
}

func (localdb *LocalDB) readSqlFile(sqlFilePath string) (string, common.PluginError) {
	if !strings.HasSuffix(sqlFilePath, ".sql") {
		err := fmt.Errorf("wrong file type")
		return "", common.NewPluginError(common.SysDbImportSqlFailure, err.Error(), common.SysDbImportSqlError.Error())
	}
	fileBytes, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return "", common.NewPluginError(common.SysDbImportSqlFailure, err.Error(), common.SysDbImportSqlError.Error())
	}
	return string(fileBytes), nil
}

func (localdb *LocalDB) fixSql(content string) (string, common.PluginError) {
	//通过正则表达式拿到对应的tableName
	re := regexp.MustCompile("{{([^{}]*)}}")
	tableNameList := re.FindAllString(content, -1)
	if len(tableNameList) == 0 {
		err := fmt.Errorf("not set tableName")
		return "", common.NewPluginError(common.SysDbImportSqlFailure, err.Error(), common.SysDbImportSqlError.Error())
	}
	for _, tableName := range tableNameList {
		newTableName := strings.TrimRight(tableName, "}}")
		newTableName = strings.TrimLeft(newTableName, "{{")
		content = strings.Replace(content, tableName, newTableName, -1)
	}
	return content, nil
}

func (localdb *LocalDB) ImportSQL(sqlFilePath string) common.PluginError {
	fileContent, err := localdb.readSqlFile(sqlFilePath)
	if err != nil {
		return err
	}
	if fileContent, err = localdb.fixSql(fileContent); err != nil {
		return err
	}

	fmt.Println("import sql:\n", fileContent)

	if err := localdb.Exec(fileContent); err != nil {
		return common.NewPluginError(common.SysDbImportSqlFailure, err.Error(), common.SysDbImportSqlError.Error())
	}
	return nil
}

func (localdb *LocalDB) Select(sql string) ([]*common.RawData, []*common.ColumnDesc, common.PluginError) {
	if err := localdb.initMysql(); err != nil {
		return nil, nil, common.NewPluginError(common.SysDbSelectFailure, err.Error(), err.Msg())
	}
	defer localdb.db.Close()

	rows, err := localdb.db.Query(sql)
	if err != nil {
		return nil, nil, common.NewPluginError(common.SysDbSelectFailure, common.SysDbSelectFailureError.Error(), err.Error())
	}
	defer rows.Close()

	built := make([]interface{}, 0)
	colDesc := make([]*common.ColumnDesc, 0)

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, common.NewPluginError(common.SysDbSelectFailure, common.SysDbSelectFailureError.Error(), err.Error())
	}
	for idx, ct := range colTypes {
		desc := &common.ColumnDesc{
			Index: int64(idx),
			Name:  ct.Name(),
			Type:  ct.DatabaseTypeName(),
		}
		colDesc = append(colDesc, desc)

		switch desc.Type {
		case "CHAR", "VARCHAR", "TEXT", "BLOB", "DATA", "DATATIME", "JSON", "TIME", "TIMESTAMP":
			built = append(built, new(string))
		case "BIT", "INT", "TINYINT", "BIGINT", "MEDIUMINT", "SMALLINT", "YEAR":
			built = append(built, new(int64))
		case "FLOAT", "DOUBLE", "DECIMAL":
			built = append(built, new(float64))
		default:
			built = append(built, new(interface{}))
		}
	}

	rawData := make([]*common.RawData, 0)
	for rows.Next() {
		if err = rows.Scan(built...); err != nil {
			return nil, nil, common.NewPluginError(common.SysDbSelectFailure, common.SysDbSelectFailureError.Error(), err.Error())
		}
		cells := make([][]byte, 0)
		for _, i := range built {
			switch reflect.TypeOf(i).String() {
			case "*string":
				cells = append(cells, []byte(reflect.ValueOf(i).Elem().String()))
			case "*int64":
				v := reflect.ValueOf(i).Elem().Int()
				byteBuf := bytes.NewBuffer([]byte{})
				if err := binary.Write(byteBuf, binary.BigEndian, v); err != nil {
					return nil, nil, common.NewPluginError(common.SysDbSelectFailure, common.SysDbSelectFailureError.Error(), err.Error())
				}
				cells = append(cells, byteBuf.Bytes())
			case "*float":
				v := reflect.ValueOf(i).Elem().Float()
				bits := math.Float64bits(v)
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, bits)
				cells = append(cells, b)
			case "*sql.NullString":
				cells = append(cells, []byte(reflect.ValueOf(i).Elem().Field(0).String()))
			case "*sql.NullInt64":
				v := reflect.ValueOf(i).Elem().Field(0).Int()
				byteBuf := bytes.NewBuffer([]byte{})
				if err := binary.Write(byteBuf, binary.BigEndian, v); err != nil {
					return nil, nil, common.NewPluginError(common.SysDbSelectFailure, err.Error(), common.SysDbSelectFailureError.Error())
				}
				cells = append(cells, byteBuf.Bytes())
			case "*sql.NullFloat64":
				v := reflect.ValueOf(i).Elem().Field(0).Float()
				bits := math.Float64bits(v)
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, bits)
				cells = append(cells, b)
			default:
				cells = append(cells, reflect.ValueOf(i).Elem().Bytes())
			}
		}
		rawData = append(rawData, &common.RawData{Cell: cells})
	}
	return rawData, colDesc, nil
}

func (localdb *LocalDB) AsyncSelect(sql string, asyncObj interface{}, callback common.SysDBCallBack) {
	rawData, columnDesc, err := localdb.Select(sql)
	callback(rawData, columnDesc, err, asyncObj)
}

func (localdb *LocalDB) Exec(sql string) common.PluginError {
	if err := localdb.initMysql(); err != nil {
		return common.NewPluginError(common.SysDbExecFailure, err.Error(), err.Msg())
	}
	defer localdb.db.Close()

	if _, err := localdb.db.Exec(sql); err != nil {
		return common.NewPluginError(common.SysDbExecFailure, common.SysDbExecFailureError.Error(), err.Error())
	}
	return nil
}

func (localdb *LocalDB) Unmarshal(rawData []*common.RawData, columnDesc []*common.ColumnDesc, v interface{}) common.PluginError {
	if len(rawData) == 0 {
		return nil
	}
	if err := utils.Unmarshal(rawData, columnDesc, v); err != nil {
		return common.NewPluginError(common.UnmarshalFailure, err.Error(), common.UnmarshalError.Error())
	}
	return nil
}
