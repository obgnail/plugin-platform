package db

import (
	"fmt"
	"reflect"
	"strings"
)

func GetTableNames(v reflect.Value, tables []string, level int) []string {
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == "TableIdent" {
			// if this is a TableIdent struct, extract the table name
			tableName := v.FieldByName("v").String()
			if tableName != "" {
				tables = append(tables, tableName)
			}
		} else {
			// otherwise enumerate all fields of the struct and process further
			for i := 0; i < v.NumField(); i++ {
				tables = GetTableNames(reflect.Indirect(v.Field(i)), tables, level+1)
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			// enumerate all elements of an array/slice and process further
			tables = GetTableNames(reflect.Indirect(v.Index(i)), tables, level+1)
		}
	case reflect.Interface:

		// get the actual object that satisfies an interface and process further
		tables = GetTableNames(reflect.Indirect(reflect.ValueOf(v.Interface())), tables, level+1)
	}

	return tables

}

func ReplaceTableName(tableNames []string, sql, instanceID string) string {
	result := sql
	for _, oldTableName := range tableNames {
		newTableName := fmt.Sprintf("`%s_%s`", instanceID, oldTableName)
		result = strings.Replace(result, oldTableName, newTableName, -1)
	}
	return result
}

func UniqueNoNullSlice(slice ...string) (newSlice []string) {
	found := make(map[string]bool)
	for _, val := range slice {
		if val == "" {
			continue
		}
		if _, ok := found[val]; !ok {
			found[val] = true
			newSlice = append(newSlice, val)
		}
	}
	return
}
