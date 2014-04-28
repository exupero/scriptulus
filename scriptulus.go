package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"fountain"
)

var root string

var addr = flag.String("addr", "8000", "HTTP service address")

func main() {
	flag.Parse()

	funcMap := template.FuncMap{
		"Mdashes": func(s string) string {
			return strings.Replace(s, "--", "—", -1)
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Println("Error reading file:", err)
			return
		}
		script := fountain.Parse(string(b))

		tmpl, err := template.New("script").Funcs(funcMap).ParseFiles(root + "/script.html")
		if err != nil {
			log.Println("Render error:", err)
			return
		}
		tmpl.ExecuteTemplate(w, "script", script)
	})

	log.Println("Starting server...")
	if err := http.ListenAndServe("0.0.0.0:" + *addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
