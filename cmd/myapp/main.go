package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = ""
)

type HelloWorld struct {
	Words string
	Date  string
	Param string
}

func main() {
	if len(date) == 0 {
		date = time.Now().UTC().Format(time.RFC3339)
	}

	fsAssets := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fsAssets))

	fsStatic := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsStatic))
	http.HandleFunc("/", serveTemplate)

	log.Println("Listening on :3000...")
	log.Println("Have a look on http://localhost:3000?param=test123")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	helloWorld := HelloWorld{Words: "Hello world !", Date: date}
	if requestParam := r.FormValue("param"); requestParam != "" {
		helloWorld.Param = requestParam
	}
	templates := template.Must(template.ParseFiles("web/templates/hello-world-template.html"))
	if err := templates.ExecuteTemplate(w, "hello-world-template.html", helloWorld); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
