package helps

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FileExist 檢查檔案是否存在
func FileExist(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat != nil && stat.IsDir() == false
}

// FileCompare 比對檔案內容與位元陣列資料
func FileCompare(path string, expected []byte) bool {
	actual, err := os.ReadFile(path)
	return err == nil && bytes.Equal(expected, actual)
}

// FileWrite 寫入檔案, 如果有需要會建立目錄
func FileWrite(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("writeFile: %w", err)
	} // if

	if err := os.WriteFile(path, data, fs.ModePerm); err != nil {
		return fmt.Errorf("writeFile: %w", err)
	} // if

	return nil
}
