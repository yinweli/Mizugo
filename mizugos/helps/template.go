package helps

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

// WriteTemplate 輸入模板字串(blueprint)與參考物件(refer), 利用go語言的 text/template 引擎在指定路徑(path)產生文字檔案, 如果有需要會建立目錄
func WriteTemplate(path, blueprint string, refer any) (err error) {
	tmpl, err := template.New(path).Parse(blueprint)

	if err != nil {
		return fmt.Errorf("writeTemplate: %w", err)
	} // if

	buffer := &bytes.Buffer{}

	if err = tmpl.Execute(buffer, refer); err != nil {
		return fmt.Errorf("writeTemplate: %w", err)
	} // if

	if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("writeTemplate: %w", err)
	} // if

	if err = os.WriteFile(path, buffer.Bytes(), fs.ModePerm); err != nil {
		return fmt.Errorf("writeTemplate: %w", err)
	} // if

	return nil
}
