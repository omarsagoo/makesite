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

//
func timeTrack(start time.Time) time.Duration {
	elapsed := time.Since(start)
	return elapsed
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

func main() {
	// dir, fileName := flagParse()
	start := time.Now()

	dir := flag.String("dir", ".", "Name of the directory to save the File")
	fileName := flag.String("file", "first-post.txt", "name of file to write to html")

	flag.Parse()
	numOfPages := 0
	var fileSizes float64
	if _, err := os.Stat(*dir); os.IsNotExist(err) == false {
		allFiles, err := ioutil.ReadDir(*dir)
		check(err)

		for _, file := range allFiles {
			if filepath.Ext(file.Name()) == ".txt" {
				fileContent := readFile(file.Name())
				fileSizes += float64(file.Size()) / float64(bytesize.KB)
				translatedFileContent, err := translateText("en", string(fileContent))
				check(err)

				buffer := writeHTMLFile(translatedFileContent)

				fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

				createHTMLFile(buffer, fileName)
				numOfPages = numOfPages + 1
			}
		}
	} else {
		fileContent := readFile(*fileName)
		buffer := writeHTMLFile(string(fileContent))
		fileName := strings.SplitN(*fileName, ".", 2)[0] + ".html"
		createHTMLFile(buffer, fileName)
		numOfPages = numOfPages + 1
	}
	bold := color.Bold.Render
	success := color.Success.Render
	since := time.Since(start).Seconds()
	fmt.Printf("%s You generated %s pages in %.2f seconds. (%.1fkB total)\n", success("Success!"), bold(numOfPages), since, fileSizes)
}
