package main

import (
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
)

var DEBUG = true

type Basic struct {
	Name string
}

type Player struct {
	Name string
	Club string	 `json:",omitempty`
}

type Players map[int] Player

type Tournament struct {
	Basic Basic
	Players Players
}

/*type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}*/

func playersEditHandler(w http.ResponseWriter, r *http.Request, t *Tournament) {
	renderTemplate(w, "players-edit.html", t)
	if r.FormValue("submit") != "" {
		// User send a request, we should override a data
	}
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func getTemplates() *template.Template {
	if DEBUG || templates == nil {
		templates = template.Must(template.ParseGlob("templates/*.html"))
	}
	return templates
}

func renderTemplate(w http.ResponseWriter, tmpl string, t *Tournament) {
	fmt.Printf("%#v\n", getTemplates().Name())
	err := getTemplates().ExecuteTemplate(w, tmpl, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/players-edit/$")
var tour = Tournament{
			Basic: Basic{Name:"Chropov 2015"},
			Players: Players {
				1: Player {"Ferdinand Mravec", "Chropov"},
				2: Player {"Jan Hrach", ""},
				3: Player {"Princ Bajaja", "Tlačiarne Púchov"},
			},
		}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Tournament)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    //fmt.Println(r.URL.Path)
		if err := r.ParseForm(); err != nil {
			http.NotFound(w, r)
			return
		}
		fmt.Printf("%#v\n", r.Form)
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, &tour)
	}
}

func main() {
	//fmt.Printf("%#v\n", tour)
	http.HandleFunc("/players-edit/", makeHandler(playersEditHandler))

	http.ListenAndServe(":2080", nil)
}

