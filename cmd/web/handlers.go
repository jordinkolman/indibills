package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)
/*
	TODO: Implement the following handler list
	Login GET
	Logout GET
	Account GET POST UPDATE DELETE
	Transaction GET POST UPDATE DELETE
	Sign Up GET POST
	Account Management GET
	- should have delete account button that issues a DELETE request to the USER record
*/

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func (app *application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	hashed_password := app.HashAndSalt([]byte(password))

	db_user, err := app.user.Get(email)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	session, err := store.Get(r, "session.id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}

	if hashed_password == db_user.Password {
		session.Values["authenticated"] = true
		session.Values["user_id"] = db_user.Id
		session.Save(r, w)
	} else {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if user logged in. if yes, redirect to account summary page
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, err := store.Get(r, "session.id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		http.Redirect(w, r, fmt.Sprintf("/users/%v/accounts", session.Values["user_id"]), http.StatusSeeOther)
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

func (app *application) accountsView(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session.id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	user_id, err := strconv.ParseInt(session.Values["user_id"].(string), 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	accounts, err := app.accountList.GetAll(user_id)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/partials/nav.html",
		"../../ui/html/pages/accounts.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", accounts)
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
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	data, err := json.Marshal(user)
	if err != nil {
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
