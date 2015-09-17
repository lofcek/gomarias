package main

import (
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
)

var DEBUG = true


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

func renderTemplate(w http.ResponseWriter, tmpl string, t interface{}) {
	fmt.Printf("%#v\n", getTemplates().Name())
	err := getTemplates().ExecuteTemplate(w, tmpl, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var tour = Tournament{
			Basic: Basic{Name:"Chropov 2015"},
			Players: Players {
				1: Player {"Ferdinand Mravec", "Chropov"},
				2: Player {"Jan Hrach", ""},
				3: Player {"Princ Bajaja", "Tlačiarne Púchov"},
			},
		}

var validPath = regexp.MustCompile("^/players-edit/([A-Za-z0-9_ ]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Tournament)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	    //fmt.Println(r.URL.Path)
		if err := r.ParseForm(); err != nil {
			http.NotFound(w, r)
			return
		}
		m := validPath.FindStringSubmatch(r.URL.Path)
		tour,err := load_tournament(m[1])
		if err != nil {
			renderTemplate(w, "error-page.html", map[string]string{
				"Error": err.Error(),
				"Tournament": m[1],
			})
		}

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, tour)
	}
}

func main() {
	//fmt.Printf("%#v\n", tour)
	http.HandleFunc("/players-edit/", makeHandler(playersEditHandler))

	http.ListenAndServe(":2080", nil)
}
