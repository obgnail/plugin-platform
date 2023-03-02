package db

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/blastrain/vitess-sqlparser/sqlparser"
	_ "github.com/go-sql-driver/mysql"
	"github.com/obgnail/plugin-platform/common/common_type"
	"github.com/obgnail/plugin-platform/common/config"
	"github.com/obgnail/plugin-platform/common/log"
	"github.com/obgnail/plugin-platform/common/protocol"
	"github.com/obgnail/plugin-platform/common/utils/file"
	"github.com/obgnail/plugin-platform/platform/service/utils"
	"io/ioutil"
	"math"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DB struct {
	DbConn *sql.DB
}

var Instance *DB
var Once sync.Once

func GetDB(dbName string) *DB {
	Once.Do(func() {
		addr := config.String("platform.mysql_address", "localhost:3306")
		user := config.String("platform.mysql_user", "root")
		pwd := config.String("platform.mysql_password", "root")
		maxIdle := config.Int("platform.mysql_db_max_idle", 10)
		maxOpen := config.Int("platform.mysql_db_max_open", 100)

		connString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
			user, pwd, addr, dbName)
		dbConn, err := sql.Open("mysql", connString)
		if err != nil {
			log.ErrorDetails(err)
		}
		dbConn.SetMaxIdleConns(maxIdle)
		dbConn.SetMaxOpenConns(maxOpen)
		dbConn.SetConnMaxLifetime(time.Hour)

		Instance = &DB{DbConn: dbConn}
	})
	return Instance
}

const (
	StmtSelect = "select"
	StmtRow    = "row"
)

// prepare 为了防止冲突,再建表时将tableName转为了instanceID_tableName
// 所以对执行的sql,也需要转化为真实的tableName
func (d *DB) prepare(rawSql, instanceID string, isLocalDB bool) (stmt string, realSql string, err common_type.PluginError) {
	statement, e := sqlparser.Parse(rawSql)
	if e != nil {
		return "", "", common_type.NewPluginError(common_type.DbSqlSyntaxErr, err.Error(), common_type.DbSqlSyntaxError.Error())
	}
	switch statement.(type) {
	case *sqlparser.Select:
		stmt = StmtSelect
	default:
		stmt = StmtRow
	}

	if !isLocalDB {
		return stmt, rawSql, nil
	}

	var tableNames []string
	tableNames = GetTableNames(reflect.Indirect(reflect.ValueOf(statement)), tableNames, 0)
	tableNames = UniqueNoNullSlice(tableNames...)
	realSql = ReplaceTableName(tableNames, rawSql, instanceID)
	return stmt, realSql, nil
}

func (d *DB) Exec(sqlStr string) common_type.PluginError {
	_, err := d.DbConn.Exec(sqlStr)
	if err != nil {
		return common_type.NewPluginError(common_type.DbExecFailure, err.Error(), common_type.DbExecFailureError.Error())
	}
	return nil
}

func (d *DB) Select(sqlStr string) (*protocol.TableMessage, common_type.PluginError) {
	rows, err := d.DbConn.Query(sqlStr)
	if err != nil {
		return nil, common_type.NewPluginError(common_type.DbSelectFailure, err.Error(), common_type.DbSelectFailureError.Error())
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, common_type.NewPluginError(common_type.DbSelectFailure, err.Error(), common_type.DbSelectFailureError.Error())
	}

	tableData := &protocol.TableMessage{
		RowData: make([]*protocol.RowMessage, 0),
		Column:  make([]*protocol.ColumnDesc, 0),
	}

	var built = make([]interface{}, 0)
	for index, ct := range colTypes {
		tmp := &protocol.ColumnDesc{}
		tmp.Name = ct.Name()
		tmp.Index = int64(index)
		tmp.Type = ct.DatabaseTypeName()
		tableData.Column = append(tableData.Column, tmp)

		switch ct.DatabaseTypeName() {
		case "CHAR", "VARCHAR", "TEXT", "BLOB", "DATA", "DATATIME", "JSON", "TIME", "TIMESTAMP":
			built = append(built, new(sql.NullString))
		case "BIT", "INT", "TINYINT", "BIGINT", "MEDIUMINT", "SMALLINT", "YEAR":
			built = append(built, new(sql.NullInt64))
		case "FLOAT", "DOUBLE", "DECIMAL":
			built = append(built, new(sql.NullFloat64))
		default:
			built = append(built, new(interface{}))
		}
	}

	for rows.Next() {
		err = rows.Scan(built...)
		if err != nil {
			return nil, common_type.NewPluginError(common_type.DbSelectFailure, err.Error(), common_type.DbSelectFailureError.Error())
		}
		tmp := &protocol.RowMessage{
			Cell: make([][]byte, 0),
		}
		for _, i := range built {
			switch reflect.TypeOf(i).String() {
			case "*string":
				tmp.Cell = append(tmp.Cell, []byte(reflect.ValueOf(i).Elem().String()))
			case "*int64":
				v := reflect.ValueOf(i).Elem().Int()
				bytebuf := bytes.NewBuffer([]byte{})
				err := binary.Write(bytebuf, binary.BigEndian, v)
				if err != nil {
					return nil, common_type.NewPluginError(common_type.DbSelectFailure, err.Error(), common_type.DbSelectFailureError.Error())
				}
				tmp.Cell = append(tmp.Cell, bytebuf.Bytes())
			case "*float":
				v := reflect.ValueOf(i).Elem().Float()
				bits := math.Float64bits(v)
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, bits)
				tmp.Cell = append(tmp.Cell, b)
			case "*sql.NullString":
				tmp.Cell = append(tmp.Cell, []byte(reflect.ValueOf(i).Elem().Field(0).String()))
			case "*sql.NullInt64":
				v := reflect.ValueOf(i).Elem().Field(0).Int()
				bytebuf := bytes.NewBuffer([]byte{})
				err := binary.Write(bytebuf, binary.BigEndian, v)
				if err != nil {
					return nil, common_type.NewPluginError(common_type.DbSelectFailure, err.Error(), common_type.DbSelectFailureError.Error())
				}
				tmp.Cell = append(tmp.Cell, bytebuf.Bytes())
			case "*sql.NullFloat64":
				v := reflect.ValueOf(i).Elem().Field(0).Float()
				bits := math.Float64bits(v)
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, bits)
				tmp.Cell = append(tmp.Cell, b)
			default:
				tmp.Cell = append(tmp.Cell, reflect.ValueOf(i).Elem().Bytes())
			}
		}
		tableData.RowData = append(tableData.RowData, tmp)
	}

	return tableData, nil
}

func (d *DB) ImportSql(sqlFileName, appUUID, version, instanceUUID string) common_type.PluginError {
	//判断是否是sql文件
	if !strings.HasSuffix(sqlFileName, ".sql") {
		err := fmt.Errorf("wrong file type")
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	pluginDir := utils.GetPluginDir(appUUID, version)
	sqlFilePath := file.JoinPath(pluginDir, sqlFileName)

	if ok, err := file.PathExists(sqlFilePath); !ok || err != nil {
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	sqlFilePath = filepath.Clean(sqlFilePath)

	fileByte, err := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}
	fileContent := string(fileByte)
	if fileContent, err = d.fixSql(fileContent, instanceUUID); err != nil {
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, err.Error(), common_type.SysDbImportSqlError.Error())
	}

	fmt.Println("import sql:\n", fileContent)

	if _, err := d.DbConn.Exec(fileContent); err != nil {
		return common_type.NewPluginError(common_type.SysDbImportSqlFailure, common_type.DbExecFailureError.Error(), err.Error())
	}
	return nil
}

func (d *DB) fixSql(content string, instanceUUID string) (string, common_type.PluginError) {
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
		newTableName = fmt.Sprintf("%s_%s", instanceUUID, newTableName) // table_name 前面添加 instanceUUID 防止冲突
		content = strings.Replace(content, tableName, newTableName, -1)
	}
	return content, nil
}
