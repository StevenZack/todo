package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Todo struct {
	Id        string
	Name      string
	Age       int
	CreatedAt string
}

//go:embed views
var fs embed.FS
var (
	todos     = []Todo{}
	idCounter = 0
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	t, e := template.ParseFS(fs, "views/*.html")
	if e != nil {
		log.Println(e)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.FormValue("method") {
		case "post":
			idCounter++
			todos = append(todos, Todo{
				Id:        strconv.Itoa(idCounter),
				Name:      r.FormValue("name"),
				CreatedAt: time.Now().Format(time.RFC3339),
			})
		case "delete":
			var out []Todo
			id := r.FormValue("id")
			for _, v := range todos {
				if v.Id == id {
					continue
				}
				out = append(out, v)
			}
			todos = out
		case "patch":
			id := r.FormValue("id")
			name := r.FormValue("name")
			age := r.FormValue("age")
			agei, e := strconv.Atoi(age)
			if e != nil {
				http.Error(w, "invalid age", http.StatusBadRequest)
				return
			}
			for i, v := range todos {
				if v.Id == id {
					todos[i].Name = name
					todos[i].Age = agei
					break
				}
			}
		default:
			e = t.ExecuteTemplate(w, "index.html", todos)
			if e != nil {
				http.Error(w, e.Error(), 500)
			}
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
	println("started at http://localhost:8080")
	e = http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Panic(e)
		return
	}

}
