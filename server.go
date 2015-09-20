package main

import (
	"html/template"
	//"io/ioutil"
	"net/http"
	"regexp"
	"fmt"
	"log"
	"encoding/json"
)

var DEBUG = true

func gettext(s string) string {
	return s
}

func tournamentSelectHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.NotFound(w, r)
		return
	}
	tour,err := getTournaments()
	if err != nil {
		renderError(w, err, "")
		return
	}
	renderTemplate(w, "tournament-select.html", tour)
}

func playersEditHandler(w http.ResponseWriter, r *http.Request, t *Tournament) {
	if r.FormValue("submit") != "" {
		// User send a request, we should override a data
		var players Players= make(Players)
		var names, clubs []string
		var ok bool
		if names,ok = r.Form["name[]"]; !ok {renderError(w, fmt.Errorf("No names in page"), t.Basic.FileName); return}
		if clubs,ok = r.Form["club[]"]; !ok {renderError(w, fmt.Errorf("No clubs in page"), t.Basic.FileName); return}
		if len(names) != len(clubs) {
			log.Panicf("names and clus has different lenght - %d, %d", len(names), len(clubs))
		}
		for i,name:= range(names) {
			players[i] = Player{name, clubs[i]}
		}
		// and remove last empty players
		for i:=len(names); i >= 0; i-- {
			if(players[i].empty()) {
				delete(players, i)
			} else {
				break
			}
		}
		//t.Players = players
		b,err := json.MarshalIndent(t, "", " ")
		if err != nil {
			renderError(w, err, t.Basic.LongName)
			return
		}
		fmt.Println("b = ", string(b))
	}
	renderTemplate(w, "players-edit.html", t)
}

func getTemplatesInit() *template.Template {
	funcMap := template.FuncMap{
		"gettext": gettext,
	}
	return template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
}
var templates = getTemplatesInit()

func getTemplates() *template.Template {
	if DEBUG || templates == nil {
		templates = getTemplatesInit()
	}
	return templates
}

func renderTemplate(w http.ResponseWriter, tmpl string, t interface{}) {
	//fmt.Printf("%#v\n", getTemplates().Name())
	err := getTemplates().ExecuteTemplate(w, tmpl, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderError(w http.ResponseWriter, err error, tournament string) {
	renderTemplate(w, "error-page.html", map[string]string{
		"Error": err.Error(),
		"Tournament": tournament,
	})
}
//var tour = Tournament{}

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
			renderError(w, err, m[1])
			return
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
	http.HandleFunc("/", tournamentSelectHandler)
	http.ListenAndServe(":2080", nil)
}
