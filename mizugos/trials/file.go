package trials

import (
	"bytes"
	"os"
)

// FileExist 檢查檔案是否存在
func FileExist(path string) bool {
	stat, err := os.Stat(path)
	return os.IsNotExist(err) == false && stat != nil && stat.IsDir() == false
}

// FileCompare 比對檔案內容與位元陣列資料
func FileCompare(path string, expected []byte) bool {
	if actual, err := os.ReadFile(path); err == nil {
		if bytes.Equal(expected, actual) {
			return true
		} // if
	} // if

	return false
}
