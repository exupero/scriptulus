package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"fountain"
)

var root string

var addr = flag.String("addr", "8000", "HTTP service address")

var regexWord *regexp.Regexp

func init() {
	regexWord = regexp.MustCompile("([A-Za-z0-9']+)")
}

type Script fountain.Document

func (script *Script) Statistics() map[string]string {
	lines := ""
	parts := make(map[string]bool)
	for _, paragraph := range script.Body {
		if paragraph.IsDialogue() {
			parts[paragraph.Speaker()] = true
			lines += paragraph.Dialogue()
		}
	}
	return map[string]string{
		"parts": fmt.Sprintf("%d", len(parts)),
		"words": fmt.Sprintf("%d", len(regexWord.FindAllString(lines, -1))),
	}
}

type Castmember struct {
	Name string
	WordCount int
}

type byWords []Castmember

func (a byWords) Len() int { return len(a) }
func (a byWords) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byWords) Less(i, j int) bool { return a[i].WordCount > a[j].WordCount }

func (script *Script) Cast() []Castmember {
	lines := make(map[string]string)
	for _, paragraph := range script.Body {
		if paragraph.IsDialogue() {
			speaker := paragraph.Speaker()
			_, ok := lines[speaker]
			if ok {
				lines[speaker] += paragraph.Dialogue()
			} else {
				lines[speaker] = paragraph.Dialogue()
			}
		}
	}

	cast := []Castmember{}
	for name, allLines := range lines {
		wordCount := len(regexWord.FindAllString(allLines, -1))
		cast = append(cast, Castmember{name, wordCount})
	}

	sort.Sort(byWords(cast))
	return cast
}

func main() {
	flag.Parse()

	funcMapHtml := template.FuncMap{
		"Mdashes": func(s string) string {
			return strings.Replace(s, "---", "—", -1)
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Println("Error reading file:", err)
			return
		}
		script := Script(*fountain.Parse(string(b)))

		tmpl, err := template.New("script").Funcs(funcMapHtml).ParseFiles(root + "/index.html")
		if err != nil {
			log.Println("Render error:", err)
			return
		}
		tmpl.ExecuteTemplate(w, "script", &script)
	})

	log.Println("Starting server...")
	if err := http.ListenAndServe("0.0.0.0:" + *addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
