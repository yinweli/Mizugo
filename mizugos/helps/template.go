package helps

import (
	"fmt"
	"strings"
	"text/template"
)

// TemplateBuild 使用 Go 的 text/template 引擎, 將 blueprint 模板與 refer 參考物件套用後, 把產生的結果寫入 builder 的尾端
func TemplateBuild(blueprint string, builder *strings.Builder, refer any) (err error) {
	tmpl, err := template.New("blueprint").Parse(blueprint)

	if err != nil {
		return fmt.Errorf("templateBuild: %w", err)
	} // if

	if err = tmpl.Execute(builder, refer); err != nil {
		return fmt.Errorf("templateBuild: %w", err)
	} // if

	return nil
}
