package file_utils

import (
	"fmt"
	"os"
)

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FindPath(path string) (string, error) {
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		} else {
			path = "../" + path
		}
	}
	return path, fmt.Errorf("no such path")
}
