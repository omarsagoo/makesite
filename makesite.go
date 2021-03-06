package main

import (
	"flag"
	"fmt"
	"makesite/blog"
	"time"

	"gopkg.in/gookit/color.v1"
)

func main() {
	start := time.Now()
	var dir string
	var fileName string
	var lang string
	var fileSizes float64
	var numOfPages int

	flag.StringVar(&dir, "dir", "", "Name of the directory to grab and save the File")
	flag.StringVar(&fileName, "file", "", "name of file to write to html")
	flag.StringVar(&lang, "lang", "en", "Language to translate the text into(default english)")

	flag.Parse()

	if dir != "" {
		fileSizes = blog.MakeMultipleHTMLfile(dir, lang, &numOfPages)
	} else if fileName != "" {
		fileSizes = blog.MakeHTMLFile(fileName, lang, &numOfPages)
	} else if fileName == "" && dir == "" {
		fmt.Printf("%s You must provide either a directory or a file!\n", color.Danger.Render("ERROR:"))
		return
	}

	bold := color.Bold.Render
	success := color.Success.Render
	since := time.Since(start).Seconds()
	fmt.Printf("%s You generated %s pages in %.2f seconds. (%.1fkB total)\n", success("Success!"), bold(numOfPages), since, fileSizes)
}
