package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"cloud.google.com/go/translate"
	"github.com/inhies/go-bytesize"
	"golang.org/x/text/language"
	"gopkg.in/gookit/color.v1"
)

type blogEntry struct {
	title   string
	content string
}

type allBlog struct {
	List []blogEntry
}

// Google API Golang snippet
func translateText(targetLanguage, text string) (string, error) {
	// text := "The Go Gopher is cute"
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
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

func createHTMLFile(buffer, filename string) {
	bytesToWrite := []byte(buffer)
	dir := "./html_SSG_files/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		check(err)
	}
	err := ioutil.WriteFile(dir+filename, bytesToWrite, 0644)
	check(err)
}

func makeMultipleHTMLfile(dir, lang string) (int, float64) {
	var numOfPages int
	var fileSizes float64

	allFiles, err := ioutil.ReadDir(dir)
	check(err)

	for _, file := range allFiles {
		if file.IsDir() {
			// recursive check for subdirectories
			return makeMultipleHTMLfile(dir+"/"+file.Name(), lang)
		}

		if filepath.Ext(file.Name()) == ".txt" {
			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			check(err)

			translatedFileContent, err := translateText(lang, string(fileContent))
			check(err)

			buffer := writeHTMLFile(translatedFileContent)

			fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

			createHTMLFile(buffer, fileName)

			fileSizes += float64(file.Size()) / float64(bytesize.KB)
			numOfPages = numOfPages + 1
		}
	}
	return numOfPages, fileSizes
}

func makeHTMLFile(fileName, lang string) (int, float64) {
	var fileSizes float64
	var numOfPages int

	file, err := os.Lstat(fileName)
	check(err)

	fileContent, err := ioutil.ReadFile(file.Name())
	check(err)

	fileSizes += float64(file.Size()) / float64(bytesize.KB)

	translatedFileContent, err := translateText(lang, string(fileContent))
	check(err)

	buffer := writeHTMLFile(translatedFileContent)

	fileName = strings.SplitN(fileName, ".", 2)[0] + ".html"

	createHTMLFile(buffer, fileName)

	numOfPages = numOfPages + 1

	return numOfPages, fileSizes
}

func main() {
	start := time.Now()
	var dir string
	var fileName string
	var lang string
	var numOfPages int
	var fileSizes float64

	flag.StringVar(&dir, "dir", "", "Name of the directory to grab and save the File")
	flag.StringVar(&fileName, "file", "", "name of file to write to html")
	flag.StringVar(&lang, "lang", "en", "Language to translate the text into(default english)")

	flag.Parse()

	if dir != "" {
		fmt.Println(lang)
		fmt.Println(dir)
		numOfPages, fileSizes = makeMultipleHTMLfile(dir, lang)
	} else if fileName != "" {
		numOfPages, fileSizes = makeHTMLFile(fileName, lang)
	} else if fileName == "" && dir == "" {
		fmt.Printf("%s You must provide either a directory or a file!\n", color.Danger.Render("ERROR:"))
		return
	}

	bold := color.Bold.Render
	success := color.Success.Render
	since := time.Since(start).Seconds()
	fmt.Printf("%s You generated %s pages in %.2f seconds. (%.1fkB total)\n", success("Success!"), bold(numOfPages), since, fileSizes)
}
