package local

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"io/ioutil"
	"math"
	"reflect"
	"regexp"
	"strings"
)

var _ common_type.LocalDB = (*LocalDB)(nil)

type LocalDB struct {
	db     *sql.DB
	plugin common_type.IPlugin
}

func NewLocalDB(plugin common_type.IPlugin) common_type.LocalDB {
	return &LocalDB{plugin: plugin}
}

func (localdb *LocalDB) initMysql() common_type.PluginError {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Asia%%2FShanghai&charset=utf8mb4&multiStatements=true",
		config.StringOrPanic("platform.mysql_user"),
		config.StringOrPanic("platform.mysql_password"),
		config.StringOrPanic("platform.mysql_host"),
		config.IntOrPanic("platform.mysql_port"),
		config.StringOrPanic("platform.mysql_db_name"),
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return common_type.NewPluginError(common_type.DataBaseNameFailure, err.Error(), err.Error())
	}
	if err = db.Ping(); err != nil {
		return common_type.NewPluginError(common_type.DataBaseNameFailure, err.Error(), err.Error())
	}
	localdb.db = db
	return nil
}

func (localdb *LocalDB) readSqlFile(sqlFilePath string) (string, common_type.PluginError) {
	if !strings.HasSuffix(sqlFilePath, ".sql") {
		err := fmt.Errorf("wrong file type")
		return "", common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	fileBytes, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return "", common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	return string(fileBytes), nil
}

func (localdb *LocalDB) fixSql(content string) (string, common_type.PluginError) {
	//通过正则表达式拿到对应的tableName
	re := regexp.MustCompile("{{([^{}]*)}}")
	tableNameList := re.FindAllString(content, -1)
	if len(tableNameList) == 0 {
		err := fmt.Errorf("not set tableName")
		return "", common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	for _, tableName := range tableNameList {
		newTableName := strings.TrimRight(tableName, "}}")
		newTableName = strings.TrimLeft(newTableName, "{{")
		content = strings.Replace(content, tableName, newTableName, -1)
	}
	return content, nil
}

func (localdb *LocalDB) ImportSQL(sqlFilePath string) common_type.PluginError {
	fileContent, err := localdb.readSqlFile(sqlFilePath)
	if err != nil {
		return err
	}
	if fileContent, err = localdb.fixSql(fileContent); err != nil {
		return err
	}

	fmt.Println("import sql:\n", fileContent)

	if err := localdb.Exec(fileContent); err != nil {
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	return nil
}

func (localdb *LocalDB) Select(sql string) ([]*common_type.RawData, []*common_type.ColumnDesc, common_type.PluginError) {
	if err := localdb.initMysql(); err != nil {
		return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, err.Error(), err.Msg())
	}
	defer localdb.db.Close()

	rows, err := localdb.db.Query(sql)
	if err != nil {
		return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, common_type.SysDbSelectFailureError.Error(), err.Error())
	}
	defer rows.Close()

	built := make([]interface{}, 0)
	colDesc := make([]*common_type.ColumnDesc, 0)

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, common_type.SysDbSelectFailureError.Error(), err.Error())
	}
	for idx, ct := range colTypes {
		desc := &common_type.ColumnDesc{
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

	rawData := make([]*common_type.RawData, 0)
	for rows.Next() {
		if err = rows.Scan(built...); err != nil {
			return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, common_type.SysDbSelectFailureError.Error(), err.Error())
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
					return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, common_type.SysDbSelectFailureError.Error(), err.Error())
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
					return nil, nil, common_type.NewPluginError(common_type.SysDbSelectFailure, err.Error(), common_type.SysDbSelectFailureError.Error())
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
		rawData = append(rawData, &common_type.RawData{Cell: cells})
	}
	return rawData, colDesc, nil
}

func (localdb *LocalDB) AsyncSelect(sql string, callback common_type.SysDBCallBack) {
	rawData, columnDesc, err := localdb.Select(sql)
	callback(rawData, columnDesc, err)
}

func (localdb *LocalDB) Exec(sql string) common_type.PluginError {
	if err := localdb.initMysql(); err != nil {
		return common_type.NewPluginError(common_type.SysDbExecFailure, err.Error(), err.Msg())
	}
	defer localdb.db.Close()

	if _, err := localdb.db.Exec(sql); err != nil {
		return common_type.NewPluginError(common_type.SysDbExecFailure, common_type.SysDbExecFailureError.Error(), err.Error())
	}
	return nil
}

func (localdb *LocalDB) Unmarshal(rawData []*common_type.RawData, columnDesc []*common_type.ColumnDesc, v interface{}) common_type.PluginError {
	if len(rawData) == 0 {
		return nil
	}
	if err := common_type.Unmarshal(rawData, columnDesc, v); err != nil {
		return common_type.NewPluginError(common_type.UnmarshalFailure, err.Error(), common_type.UnmarshalError.Error())
	}
	return nil
}
