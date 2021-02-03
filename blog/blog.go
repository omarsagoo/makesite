package blog

import "html/template"

// Entry stores the data for a blog entry
type Entry struct {
	Content template.HTML
}

type allBlog struct {
	List []Entry
}

// check is a wrapper function to check for an error
func check(err error) {
	if err != nil {
		panic(err)
	}
}
