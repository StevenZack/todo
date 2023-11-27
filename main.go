package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	Todo struct {
		Id        string
		Name      string
		Age       int
		CreatedAt string
	}
	User struct {
		Email     string
		Password  string
		SessionId string
	}
)

const (
	PageSize = 2
)

var (
	todos       = []Todo{}
	users       = []User{}
	modifiedTag = strconv.FormatInt(time.Now().Unix(), 10)
	idCounter   = 0
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	t, e := template.ParseGlob("views/*.html")
	if e != nil {
		log.Println(e)
		return
	}

	http.HandleFunc("/nav", func(w http.ResponseWriter, r *http.Request) {
		var username string
		if sessionId, e := r.Cookie("at"); e == nil {
			for _, v := range users {
				if v.SessionId == sessionId.Value {
					username = v.Email
					break
				}
			}
		}

		t.ExecuteTemplate(w, "nav.html", username)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			email := r.FormValue("email")
			password := r.FormValue("password")
			if email == "" {
				http.Error(w, "invalid email", 400)
				return
			}
			if password == "" {
				http.Error(w, "invalid password", 400)
				return
			}
			for _, v := range users {
				if v.Email == email {
					http.Error(w, "account exists", 400)
					return
				}
			}
			v := User{
				Email:     email,
				Password:  password,
				SessionId: strconv.FormatInt(time.Now().Unix(), 10),
			}
			users = append(users, v)
			http.SetCookie(w, &http.Cookie{
				Name:     "at",
				Value:    v.SessionId,
				HttpOnly: true,
			})
			fmt.Fprint(w, "/")
			return
		}
		t.ExecuteTemplate(w, "register.html", nil)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			email := r.FormValue("email")
			password := r.FormValue("password")
			if email == "" {
				http.Error(w, "invalid email", 400)
				return
			}
			if password == "" {
				http.Error(w, "invalid password", 400)
				return
			}
			for i, v := range users {
				if v.Email == email {
					if v.Password != password {
						http.Error(w, "password incorrect", 400)
						return
					}
					v.SessionId = strconv.FormatInt(time.Now().Unix(), 10)
					users[i] = v
					http.SetCookie(w, &http.Cookie{
						Name:     "at",
						Value:    v.SessionId,
						HttpOnly: true,
					})
					fmt.Fprint(w, "/")
					return
				}
			}
			http.Error(w, "account not found", 400)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "at",
			Expires: time.Now(),
			MaxAge:  -1,
		})
		t.ExecuteTemplate(w, "login.html", nil)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		var page = 0
		pageStr := r.FormValue("page")
		if pageStr != "" {
			var e error
			page, e = strconv.Atoi(pageStr)
			if e != nil {
				http.Error(w, "invalid page", http.StatusBadRequest)
				return
			}
		}

		switch r.Method {
		case http.MethodPost:
			idCounter++
			todos = append(todos, Todo{
				Id:        strconv.Itoa(idCounter),
				Name:      r.FormValue("name"),
				CreatedAt: time.Now().Format(time.RFC3339),
			})
		case http.MethodDelete:
			var out []Todo
			id := r.FormValue("id")
			for _, v := range todos {
				if v.Id == id {
					continue
				}
				out = append(out, v)
			}
			todos = out
		case http.MethodPatch:
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
			w.Header().Set("ETag", modifiedTag)
			if modifiedTag == r.Header.Get("If-None-Match") {
				http.Redirect(w, r, "/", http.StatusNotModified)
				return
			}

			var out []Todo
			for i := PageSize * page; i < len(todos); i++ {
				out = append(out, todos[i])
				if len(out) >= PageSize {
					break
				}
			}
			e = t.ExecuteTemplate(w, "index.html", map[string]any{
				"List":      out,
				"Last":      page - 1,
				"Page":      page,
				"Next":      page + 1,
				"Total":     len(todos),
				"TotalPage": len(todos) / PageSize,
			})
			if e != nil {
				http.Error(w, e.Error(), 500)
			}
			return
		}
		modifiedTag = strconv.FormatInt(time.Now().Unix(), 10)
	})
	println("started at http://localhost:8080")
	e = http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Panic(e)
		return
	}
}
