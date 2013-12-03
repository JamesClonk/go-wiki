package main

import "os"
import "log"
import "regexp"
import "errors"
import "io/ioutil"
import "net/http"
import "html/template"

type Page struct {
	Title string
	Body  template.HTML
}

var isLogging = true
var dataDirectory = "data"
var viewDirectory = "view"
var defaultPage = "FrontPage"
var templates = template.Must(template.ParseGlob(viewDirectory + "/*"))
var validPath = regexp.MustCompile("^/(view|edit|save)/([a-zA-Z0-9_]+)$")
var pageLinking = regexp.MustCompile(`\[([a-zA-Z0-9_]+)\]`)
var newLines = regexp.MustCompile(`\n`)

func main() {
	checkDirectories()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	http.ListenAndServe(":8008", nil)
}

func checkDirectories() {
	_, err := os.Open(dataDirectory)
	if err != nil {
		log.Fatalf("There is a problem with the data directory! [%v]", err)
	}
	_, err = os.Open(viewDirectory)
	if err != nil {
		log.Fatalf("There is a problem with the view directory! [%v]", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if isLogging {
		log.Printf("%v - %v[%v] *** redirect to /view/"+defaultPage, r.RemoteAddr, r.Method, r.RequestURI)
	}
	http.Redirect(w, r, "/view/"+defaultPage, http.StatusFound)
}

func makeHandler(handler func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isLogging {
			log.Printf("%v - %v[%v]", r.RemoteAddr, r.Method, r.RequestURI)
		}

		title, err := getTitle(w, r)
		if err != nil {
			return
		}

		handler(w, r, title)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	pH := []byte(p.Body)
	pH = pageLinking.ReplaceAll(pH, []byte("<a href=\"/wiki/view/${1}\">${1}</a>"))
	pH = newLines.ReplaceAll(pH, []byte("<br/>"))
	p.Body = template.HTML(pH)

	_renderTemplate(w, "view.html", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	_renderTemplate(w, "edit.html", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: template.HTML(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func _renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl, p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}

func getFilename(title string) string {
	return dataDirectory + "/" + title + ".txt"
}

func (p *Page) save() error {
	filename := getFilename(p.Title)
	return ioutil.WriteFile(filename, []byte(template.HTMLEscapeString(string(p.Body))), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := getFilename(title)
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: template.HTML(body)}, nil
}
