package file_utils

import "os"

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
			break
		} else {
			path = "../" + path
		}
	}
	return path, nil
}
