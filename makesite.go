package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func writeHTMLFile(fileContent string) string {
	paths := []string{
		"template.tmpl",
	}

	buffer := new(bytes.Buffer)

	Data := data{string(fileContent)}

	tmpl := template.Must(template.New("template.tmpl").ParseFiles(paths...))

	err := tmpl.Execute(buffer, Data)
	check(err)

	return buffer.String()
}

func createHTMLFile(buffer, filename string) bool {
	bytesToWrite := []byte(buffer)
	err := ioutil.WriteFile(filename, bytesToWrite, 0644)
	check(err)

	return true
}

// helper function, help from Audi
func flagCheck(name string) bool {
	check := false
	fmt.Println(name)
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f)

		if f.Name == name {
			check = true
		}
	})

	return check
}

func main() {
	// dir, fileName := flagParse()

	dir := flag.String("dir", ".", "Name of the directory to save the File")
	fileName := flag.String("file", "first-post.txt", "name of file to write to html")

	flag.Parse()

	if flag.Args()[0] == *dir {
		allFiles, err := ioutil.ReadDir(*dir)
		check(err)

		for _, file := range allFiles {
			if filepath.Ext(file.Name()) == ".txt" {
				fileContent := readFile(file.Name())

				buffer := writeHTMLFile(string(fileContent))

				fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

				createHTMLFile(buffer, fileName)
			}
		}
	} else {
		fileContent := readFile(*fileName)
		fmt.Println("asga")
		buffer := writeHTMLFile(string(fileContent))
		fileName := strings.SplitN(*fileName, ".", 2)[0] + ".html"
		createHTMLFile(buffer, fileName)
	}
}
