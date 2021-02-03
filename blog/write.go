package blog

import (
	"bytes"
	"html/template"
)

func writeHTMLFile(fileContent string) string {
	buffer := new(bytes.Buffer)

	Data := Entry{Content: template.HTML(fileContent)}

	tmpl, err := template.ParseFiles("template.tmpl")
	check(err)

	err = tmpl.Execute(buffer, Data)
	check(err)

	return buffer.String()
}
