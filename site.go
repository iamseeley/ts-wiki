package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"time"

	"github.com/go-yaml/yaml"
	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Title string
	Body  []byte
}

type Journal struct {
	Title     string
	Body      []byte
	Image     string
	UpdatedAt time.Time
}

func loadPageFromDirectory(directory, title string) (*Page, error) {
	filename := directory + title + ".md"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func loadJournalFromDirectory(directory, title string) (*Journal, error) {
	filename := directory + title + ".md"
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	frontMatter, body, err := parseFrontMatter(content)
	if err != nil {
		return nil, err
	}

	var journal Journal

	// Extract and set front matter data into the Journal struct
	if title, ok := frontMatter["title"].(string); ok {
		journal.Title = title
	}

	if image, ok := frontMatter["image"].(string); ok {
		journal.Image = image
	}

	// if tags, ok := frontMatter["tags"].([]interface{}); ok {
	// 	for _, tag := range tags {
	// 		if tagStr, ok := tag.(string); ok {
	// 			journal.Tags = append(journal.Tags, tagStr)
	// 		}
	// 	}
	// }

	if updatedAtStr, ok := frontMatter["updated_at"].(string); ok {
		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err == nil {
			journal.UpdatedAt = updatedAt
		}
	}

	journal.Body = body

	return &journal, nil
}

func parseFrontMatter(content []byte) (map[string]interface{}, []byte, error) {
	frontMatter := make(map[string]interface{})
	var contentStart int

	// Find the position of the first `---` delimiter
	delimiter := []byte("---")
	start := bytes.Index(content, delimiter)
	if start == -1 {
		return nil, nil, errors.New("Front matter delimiter not found")
	}

	// Find the position of the second `---` delimiter
	end := bytes.Index(content[start+len(delimiter):], delimiter)
	if end == -1 {
		return nil, nil, errors.New("Second front matter delimiter not found")
	}

	// Parse the front matter
	if err := yaml.Unmarshal(content[start+len(delimiter):start+len(delimiter)+end], &frontMatter); err != nil {
		return nil, nil, err
	}

	// Find the start of the actual content
	contentStart = start + len(delimiter) + end + len(delimiter)

	// Extract the front matter and content separately
	actualContent := content[contentStart:]

	return frontMatter, actualContent, nil
}

func pageHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := loadPageFromDirectory("pages/", title)
	if err != nil {
		http.Redirect(w, r, "/site/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "site", p)
}

func journalHandler(w http.ResponseWriter, r *http.Request, title string) {

	journal, err := loadJournalFromDirectory("journal/", title)
	if err != nil {
		http.Redirect(w, r, "/journal/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "journal", journal)
}

func markDowner(args ...interface{}) template.HTML {
	s := blackfriday.Run([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

var templates = template.Must(template.New("").Funcs(template.FuncMap{"markDown": markDowner}).ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, content interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(site|journal)/([a-zA-Z0-9]+)$")

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

func setCacheHeaders(w http.ResponseWriter, maxAge int) {
	// Set cache control headers to enable caching for the specified duration
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
}

func main() {
	http.Handle("/", http.RedirectHandler("/site/home", http.StatusSeeOther))
	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache headers to cache assets for 1 hour (you can adjust this duration)
		setCacheHeaders(w, 3600)
		fs.ServeHTTP(w, r)
	})))

	// http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/site/", makeHandler(pageHandler))
	http.HandleFunc("/journal/", makeHandler(journalHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
