package blog

import (
	"io/ioutil"
	"makesite/translate"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/inhies/go-bytesize"
)

// MakeMultipleHTMLfile will make multiple html files
func MakeMultipleHTMLfile(dir, lang string, numOfPages *int) float64 {
	var fileSizes float64

	allFiles, err := ioutil.ReadDir(dir)
	check(err)

	for _, file := range allFiles {
		if file.IsDir() {
			// recursive check for subdirectories
			MakeMultipleHTMLfile(dir+"/"+file.Name(), lang, numOfPages)
		}

		if filepath.Ext(file.Name()) == ".txt" {
			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			check(err)

			translatedFileContent, err := translate.Translate(lang, string(fileContent))
			check(err)

			writeHTMLFile(file.Name(), translatedFileContent)

			fileSizes += float64(file.Size()) / float64(bytesize.KB)
			*numOfPages++
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

			writeHTMLFile(file.Name(), string(html))

			fileSizes += float64(file.Size()) / float64(bytesize.KB)
			*numOfPages++
		}
	}
	return fileSizes
}

// MakeHTMLFile will make a new HTML file
func MakeHTMLFile(fileName, lang string, numOfPages *int) float64 {
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

		translatedFileContent = string(markdown.ToHTML([]byte(translatedFileContent), p, renderer))
	}

	writeHTMLFile(file.Name(), translatedFileContent)
	*numOfPages++
	return fileSizes
}
