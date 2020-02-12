package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"github.com/inhies/go-bytesize"
	"golang.org/x/text/language"
	"gopkg.in/gookit/color.v1"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var numOfPages int

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
	Content template.HTML
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func writeHTMLFile(fileContent string) string {
	buffer := new(bytes.Buffer)

	Data := data{Content: template.HTML(fileContent)}

	tmpl, err := template.ParseFiles("template.tmpl")
	check(err)

	err = tmpl.Execute(buffer, Data)
	check(err)

	return buffer.String()
}

func createHTMLFile(buffer, filename string) {
	numOfPages = numOfPages + 1
	bytesToWrite := []byte(buffer)
	dir := "./html_SSG_files/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		check(err)
	}
	err := ioutil.WriteFile(dir+filename, bytesToWrite, 0644)
	check(err)
}

func makeMultipleHTMLfile(dir, lang string) float64 {
	var fileSizes float64

	allFiles, err := ioutil.ReadDir(dir)
	check(err)

	for _, file := range allFiles {
		if file.IsDir() {
			// recursive check for subdirectories
			makeMultipleHTMLfile(dir+"/"+file.Name(), lang)
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
		} else if filepath.Ext(file.Name()) == ".md" {
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs
			p := parser.NewWithExtensions(extensions)

			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)

			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			check(err)

			translatedFileContent, err := translateText(lang, string(fileContent))
			check(err)

			html := markdown.ToHTML([]byte(translatedFileContent), p, renderer)

			buffer := writeHTMLFile(string(html))

			fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

			createHTMLFile(buffer, fileName)

			fileSizes += float64(file.Size()) / float64(bytesize.KB)

		}
	}
	return fileSizes
}

func makeHTMLFile(fileName, lang string) float64 {
	var fileSizes float64

	file, err := os.Lstat(fileName)
	check(err)

	fileContent, err := ioutil.ReadFile(fileName)
	check(err)

	fileSizes += float64(file.Size()) / float64(bytesize.KB)

	translatedFileContent, err := translateText(lang, string(fileContent))
	check(err)

	if filepath.Ext(file.Name()) == ".md" {
		extensions := parser.CommonExtensions | parser.AutoHeadingIDs
		p := parser.NewWithExtensions(extensions)

		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		opts := html.RendererOptions{Flags: htmlFlags}
		renderer := html.NewRenderer(opts)

		translatedFileContent := markdown.ToHTML([]byte(translatedFileContent), p, renderer)
		buffer := writeHTMLFile(string(translatedFileContent))

		fileName = strings.Replace(file.Name(), ".md", ".html", 1)

		createHTMLFile(buffer, fileName)

		return fileSizes
	}

	buffer := writeHTMLFile(translatedFileContent)

	fileName = strings.Replace(file.Name(), ".txt", ".html", 1)

	createHTMLFile(buffer, fileName)

	return fileSizes
}

func main() {
	start := time.Now()
	var dir string
	var fileName string
	var lang string
	var fileSizes float64

	flag.StringVar(&dir, "dir", "", "Name of the directory to grab and save the File")
	flag.StringVar(&fileName, "file", "", "name of file to write to html")
	flag.StringVar(&lang, "lang", "en", "Language to translate the text into(default english)")

	flag.Parse()

	if dir != "" {
		fileSizes = makeMultipleHTMLfile(dir, lang)
	} else if fileName != "" {
		fileSizes = makeHTMLFile(fileName, lang)
	} else if fileName == "" && dir == "" {
		fmt.Printf("%s You must provide either a directory or a file!\n", color.Danger.Render("ERROR:"))
		return
	}

	bold := color.Bold.Render
	success := color.Success.Render
	since := time.Since(start).Seconds()
	fmt.Printf("%s You generated %s pages in %.2f seconds. (%.1fkB total)\n", success("Success!"), bold(numOfPages), since, fileSizes)
}
