package blog

import (
	"html/template"
	"os"
	"strings"
)

func writeHTMLFile(filename, fileContent string) {
	newFile, err := os.Create("html_SSG_files/" + strings.SplitN(filename, ".", 2)[0] + ".html")
	check(err)

	Data := Entry{Content: template.HTML(fileContent)}

	tmpl, err := template.ParseFiles("template.tmpl")
	check(err)

	err = tmpl.Execute(newFile, Data)
	check(err)
}
