package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	common "github.com/obgnail/plugin-platform/common_type"
	"math"
	"reflect"
	"strconv"
)

func Unmarshal(rawData []*common.RawData, columnDesc []*common.ColumnDesc, v interface{}) error {
	val := reflect.Indirect(reflect.ValueOf(v))
	typ := val.Type()
	for _, r := range rawData {
		mVal := reflect.Indirect(reflect.New(typ.Elem().Elem())).Addr()
		for key, val := range r.Cell {
			var value interface{}
			switch columnDesc[key].Type {
			case "CHAR", "VARCHAR", "TEXT", "BLOB", "DATA", "DATATIME", "JSON", "TIME", "TIMESTAMP":
				value = string(val)
			case "BIT", "INT", "TINYINT", "BIGINT", "MEDIUMINT", "SMALLINT", "YEAR":
				value = int64(binary.BigEndian.Uint64(val))
			case "FLOAT", "DOUBLE", "DECIMAL":
				bits := binary.LittleEndian.Uint64(val)
				value = math.Float64frombits(bits)
			}
			err := setField(mVal.Interface(), columnDesc[key].Name, value)
			if err != nil {
				fmt.Println("setField", err.Error())
				return err
			}
		}
		val = reflect.Append(val, mVal)
	}

	err := deepCopy(v, val.Interface())
	if err != nil {
		fmt.Println("deepCopy err", err.Error())
	}
	return err
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.TypeOf(obj).Elem()
	var structFieldValue reflect.Value
	var j int
	for i := 0; i < structValue.NumField(); i++ {
		if structValue.Field(i).Tag.Get("orm") == name {
			structFieldValue = reflect.ValueOf(obj).Elem().Field(i)
			j++
		}
	}
	if j == 0 {
		return nil
	}
	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such structfield name: %s in obj", name)
	}
	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set structfield name %s value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	s := toString(val)

	switch structFieldType.String() {
	case "int", "int8", "int16", "int32", "int64":
		i64, err := strconv.ParseInt(s, 10, val.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", obj, s, val.Kind(), err)
		}
		structFieldValue.SetInt(i64)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		u64, err := strconv.ParseUint(s, 10, val.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", obj, s, val.Kind(), err)
		}
		structFieldValue.SetUint(u64)
	case "float32", "float64":
		f64, err := strconv.ParseFloat(s, val.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", obj, s, val.Kind(), err)
		}
		structFieldValue.SetFloat(f64)
	case "string":
		structFieldValue.Set(reflect.ValueOf(s))
	}
	return nil
}

func toString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}
