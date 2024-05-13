package helps

import (
	"fmt"
	"strings"
	"text/template"
)

// TemplateBuild 輸入模板字串(blueprint)與參考物件(refer), 利用go語言的 text/template 引擎在字串建造器(builder)的尾端產生字串
func TemplateBuild(blueprint string, builder *strings.Builder, refer any) (err error) {
	tmpl, err := template.New("template").Parse(blueprint)

	if err != nil {
		return fmt.Errorf("templateBuild: %w", err)
	} // if

	if err = tmpl.Execute(builder, refer); err != nil {
		return fmt.Errorf("templateBuild: %w", err)
	} // if

	return nil
}
