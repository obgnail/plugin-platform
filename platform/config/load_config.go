package config

import (
	"encoding/json"
	"fmt"
	"github.com/obgnail/plugin-platform/utils/errors"
	"io/ioutil"
	"os"
)

func InitConfig() error {
	configPath := "./platform/config/config.json"
	if err := LoadConfigs(configPath); err != nil {
		return errors.Trace(err)
	}
	return nil
}

var (
	config map[string]interface{}
)

func LoadConfigs(configPath string) (err error) {
	// 向上找5层，满足在一些单元测试中加载不了配置文件的问题
	for i := 0; i < 5; i++ {
		if _, err = os.Stat(configPath); err == nil {
			break
		} else {
			configPath = "../" + configPath
		}
	}

	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		panic(err)
		return
	}
	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
		return
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	return
}

func StringSlice(key string) []string {
	rawValues, _ := getArray(key, nil)
	if len(rawValues) == 0 {
		return nil
	}
	values := make([]string, 0, len(rawValues))
	for _, rawValue := range rawValues {
		value, ok := rawValue.(string)
		if ok && value != "" {
			values = append(values, value)
		}
	}
	return values
}

func getArray(key string, defaultValue []interface{}) (value []interface{}, found bool) {
	v, found := config[key]
	if !found {
		return defaultValue, false
	}
	if value, ok := v.([]interface{}); ok {
		return value, ok
	} else {
		return defaultValue, false
	}
}

func ArrayOrError(key string) (value []interface{}, err error) {
	value, found := getArray(key, []interface{}{})
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func ArrayOrPanic(key string) []interface{} {
	v, err := ArrayOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func String(key string, defaultValue string) (value string) {
	value, _ = getString(key, defaultValue)
	return
}

func getString(key string, defaultValue string) (value string, found bool) {
	v, found := config[key]
	if !found {
		return defaultValue, false
	}
	if value, ok := v.(string); ok {
		return value, ok
	} else {
		return defaultValue, false
	}
}

func StringOrError(key string) (value string, err error) {
	value, found := getString(key, "")
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func StringOrPanic(key string) string {
	v, err := StringOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Int(key string, defaultValue int) (value int) {
	value, _ = getInt(key, defaultValue)
	return
}

func getInt(key string, defaultValue int) (value int, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if v64, found := v.(float64); found {
		return int(v64), found
	} else {
		return defaultValue, false
	}
}

func IntOrError(key string) (value int, err error) {
	value, found := getInt(key, 0)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func IntOrPanic(key string) int {
	v, err := IntOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Int64(key string, defaultValue int64) (value int64, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if v64, found := v.(float64); found {
		return int64(v64), found
	} else {
		return defaultValue, false
	}
}

func Int64OrError(key string) (value int64, err error) {
	value, found := Int64(key, 0)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func Int64OrPanic(key string) int64 {
	v, err := Int64OrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Bool(key string, defaultValue bool) (value bool) {
	value, _ = getBool(key, defaultValue)
	return
}

func getBool(key string, defaultValue bool) (value bool, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if b, ok := v.(bool); ok {
		return b, ok
	} else if i, ok := v.(int); ok {
		if i == 1 {
			return true, true
		} else if i == 0 {
			return false, true
		} else {
			return defaultValue, false
		}
	} else {
		return defaultValue, false
	}
}

func BoolOrError(key string) (value bool, err error) {
	value, found := getBool(key, false)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func BoolOrPanic(key string) bool {
	v, err := BoolOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}
