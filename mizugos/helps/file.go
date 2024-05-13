package helps

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteFile 寫入檔案, 如果有需要會建立目錄
func WriteFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("writeFile: %w", err)
	} // if

	if err := os.WriteFile(path, data, fs.ModePerm); err != nil {
		return fmt.Errorf("writeFile: %w", err)
	} // if

	return nil
}
