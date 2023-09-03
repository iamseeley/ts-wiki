package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Title string

	Body []byte
}

func (p *Page) save(dir string) error {
	filename := dir + p.Title + ".md"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPageFromDirectory(directory, title string) (*Page, error) {
	filename := directory + title + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPageFromDirectory("pages/", title)
	if err != nil {
		http.Redirect(w, r, "/site/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "site", p)
}

// func editHandler(w http.ResponseWriter, r *http.Request, title string) {
// 	p, err := loadPageFromDirectory("data/", title)
// 	if err != nil {
// 		p = &Page{Title: title}
// 	}

// 	renderTemplate(w, "edit", p)
// }

// func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
// 	body := r.FormValue("body")
// 	p := &Page{Title: title, Body: []byte(body)}
// 	err := p.save("data/")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/site/"+title, http.StatusFound)
// }

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

var templates = template.Must(template.New("").Funcs(template.FuncMap{"markDown": markDowner}).ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(site)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
	}
}

func main() {
	http.Handle("/", http.RedirectHandler("/site/home", http.StatusSeeOther))
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/site/", makeHandler(viewHandler))
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
