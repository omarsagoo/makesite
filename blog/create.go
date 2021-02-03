package blog

import (
	"io/ioutil"
	"makesite/translate"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/inhies/go-bytesize"
)

// CreateHTMLFile creates an html file
func CreateHTMLFile(buffer, filename string) {
	bytesToWrite := []byte(buffer)
	dir := "./html_SSG_files/"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0777)
		check(err)
	}
	err := ioutil.WriteFile(dir+filename, bytesToWrite, 0644)
	check(err)
}

// MakeMultipleHTMLfile will make multiple html files
func MakeMultipleHTMLfile(dir, lang string) float64 {
	var fileSizes float64

	allFiles, err := ioutil.ReadDir(dir)
	check(err)

	for _, file := range allFiles {
		if file.IsDir() {
			// recursive check for subdirectories
			MakeMultipleHTMLfile(dir+"/"+file.Name(), lang)
		}

		if filepath.Ext(file.Name()) == ".txt" {
			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			check(err)

			translatedFileContent, err := translate.Translate(lang, string(fileContent))
			check(err)

			buffer := writeHTMLFile(translatedFileContent)

			fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

			CreateHTMLFile(buffer, fileName)

			fileSizes += float64(file.Size()) / float64(bytesize.KB)
		} else if filepath.Ext(file.Name()) == ".md" {
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs
			p := parser.NewWithExtensions(extensions)

			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)

			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			check(err)

			translatedFileContent, err := translate.Translate(lang, string(fileContent))
			check(err)

			html := markdown.ToHTML([]byte(translatedFileContent), p, renderer)

			buffer := writeHTMLFile(string(html))

			fileName := strings.SplitN(file.Name(), ".", 2)[0] + ".html"

			CreateHTMLFile(buffer, fileName)

			fileSizes += float64(file.Size()) / float64(bytesize.KB)

		}
	}
	return fileSizes
}

// MakeHTMLFile will make a new HTML file
func MakeHTMLFile(fileName, lang string) float64 {
	var fileSizes float64

	file, err := os.Lstat(fileName)
	check(err)

	fileContent, err := ioutil.ReadFile(fileName)
	check(err)

	fileSizes += float64(file.Size()) / float64(bytesize.KB)

	translatedFileContent, err := translate.Translate(lang, string(fileContent))
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

		CreateHTMLFile(buffer, fileName)

		return fileSizes
	}

	buffer := writeHTMLFile(translatedFileContent)

	fileName = strings.Replace(file.Name(), ".txt", ".html", 1)

	CreateHTMLFile(buffer, fileName)

	return fileSizes
}
