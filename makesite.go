package main

import (
	"html/template"
	"io/ioutil"
	"os"
)

type blogEntry struct {
	title   string
	content string
}

type allBlog struct {
	List []blogEntry
}

func readFile(file string) []byte {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	// for i := 0; string(fileContents[i]) != "\n"; i++ {
	// 	if string(fileContents[i+1]) == "\n" {
	// 		title = string(fileContents[:i+1])
	// 	}
	// }
	return fileContents
}

type data struct {
	Content string
}

func writeHTMLFile(file string) {
	s := data{}
	s.Content = string(file)
	paths := []string{
		"template.tmpl",
	}

	tmpl := template.Must(template.New("template.tmpl").ParseFiles(paths...))
	err := tmpl.Execute(os.Stdout, s)
	if err != nil {
		panic(err)
	}
}

func main() {

	FileContent := readFile("template.tmpl")

	error := ioutil.WriteFile("first-post.html", FileContent, 0644)
	if error != nil {
		panic(error)
	}
}
