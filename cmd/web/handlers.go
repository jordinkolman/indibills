package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/partials/nav.html",
		"../../ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) usersView(w http.ResponseWriter, r *http.Request) {
	users, err := app.userList.GetAll()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/partials/nav.html",
		"../../ui/html/pages/users.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", users)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.userCreateForm(w)
	case http.MethodPost:
		app.userCreateProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) userCreateForm(w http.ResponseWriter) {
	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/partials/nav.html",
		"../../ui/html/pages/create.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}


func (app *application) userCreateProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := app.HashAndSalt([]byte(r.PostForm.Get("password")))
	firstName := r.PostForm.Get("firstName")
	lastName := r.PostForm.Get("lastName")

	user := struct {
		Email string `json:"email"`
		Password  string    `json:"password"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
	}{
		Email: email,
		Password: password,
		FirstName: firstName,
		LastName: lastName,
	}

	data, err := json.Marshal(user)
	if err != nil{
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", app.userList.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status: %s", resp.Status)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
