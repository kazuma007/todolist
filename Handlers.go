package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

type Handlers struct {
	todoitem *TodoItem
}

var (
	// キーの長さは 16, 24, 32 バイトのいずれかでなければならない。
	// (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

var formTmpl *template.Template
var topTmpl *template.Template

func init() {
	formTmpl = template.Must(template.ParseFiles("./html/form.html"))
	topTmpl = template.Must(template.ParseFiles("./html/top.html"))
}

func NewHandlers(todoitem *TodoItem) *Handlers {
	return &Handlers{todoitem: todoitem}
}

func (handlers *Handlers) TopPageHandler(w http.ResponseWriter, r *http.Request) {
	if err := topTmpl.Execute(w, nil); err != nil {
		return
	}
}

func (handlers *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusNotFound
		http.Error(w, http.StatusText(code), code)
		return
	}

	keyword := r.FormValue("keyword")
	SetSessionValue(w, r, "keyword", keyword)

	items, err := handlers.todoitem.GetItems(keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := formTmpl.Execute(w, items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handlers *Handlers) FormPageHandler(w http.ResponseWriter, r *http.Request) {

	keyword := GetSessionValue(w, r, "keyword")
	if keyword == "" {
		topTmpl.Execute(w, nil)
		return
	}

	items, err := handlers.todoitem.GetItems(keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := formTmpl.Execute(w, items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handlers *Handlers) SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusNotFound
		http.Error(w, http.StatusText(code), code)
		return
	}

	keyword := GetSessionValue(w, r, "keyword")
	if keyword == "" {
		topTmpl.Execute(w, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handlers.todoitem.DeleteItems(keyword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, values := range r.PostForm {
		for _, value := range values {
			if value == "" {
				continue
			}
			item := &Item{Content: value, Keyword: keyword}
			if err := handlers.todoitem.SaveItem(item); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/top", http.StatusFound)
}

func initSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, "_golang_cookie")
	if err != nil {
		panic(err)
	}
	return session
}

func SetSessionValue(w http.ResponseWriter, r *http.Request, key, value string) {
	session := initSession(r)
	session.Values[key] = value
	session.Save(r, w)
}

func GetSessionValue(w http.ResponseWriter, r *http.Request, key string) string {
	session := initSession(r)
	valWithOutType := session.Values[key]
	value, ok := valWithOutType.(string)
	if !ok {
		fmt.Println("cannot get session value by key: " + key)
	}

	return value
}
