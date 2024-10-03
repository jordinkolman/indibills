package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

/*
	TODO: Implement the following handler list
	Account POST UPDATE DELETE
	Transaction GET POST UPDATE DELETE
	Sign Up GET POST
	Account Management GET
	- should have delete account button that issues a DELETE request to the USER record
*/

// Authentication Handlers
func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "indibills_session")
	if err != nil {
		log.Print(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	if session.Values["authenticated"] != nil && session.Values["authenticated"] != false {
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

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
		"../../ui/html/pages/signup.html",
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
	fmt.Println("made it here")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	hashed_password := fmt.Sprintf("%v", app.HashAndSalt([]byte(password)))
	firstName := r.PostForm.Get("firstName")
	lastName := r.PostForm.Get("lastName")

	user := struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{
		Email:     email,
		Password:  hashed_password,
		FirstName: firstName,
		LastName:  lastName,
	}

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%v/users/%v", app.user.Endpoint, email)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
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

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.loginGet(w, r)
	case http.MethodPost:
		app.loginPost(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) loginGet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "indibills_session")
	if err != nil {
		log.Print(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	if session.Values["authenticated"] != nil && session.Values["authenticated"] != false {
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
		return
	}

	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/partials/nav.html",
		"../../ui/html/pages/login.html",
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

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	db_user, err := app.user.Get(email)
	if err != nil || db_user == nil {
		fmt.Println("couldn't retrieve user")
		http.Error(w, fmt.Sprintf("%v", err), http.StatusNotFound)
		return
	}
	session, err := store.Get(r, "indibills_session")
	fmt.Println(session.IsNew)
	fmt.Println(session.Name())
	if err != nil {
		fmt.Println("couldn't establish session")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	session.Values["authenticated"] = true
	session.Values["user_id"] = db_user.Id
	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		fmt.Println("could not save session")
	}
	fmt.Println(session.Values)
	session.AddFlash("test flash")
	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "indibills_session")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// index handler (indibills.com/) -> redirects to indibills.com/accounts if authenticated or indibills.com/login if not
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if user logged in. if yes, redirect to account summary page
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, err := store.Get(r, "indibills_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	fmt.Println(session.Values)

	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		http.Redirect(w, r, "/accounts", http.StatusSeeOther)
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

// account handlers
func (app *application) accountsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "indibills_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	fmt.Println(session.Values)
	if (session.Values["authenticated"] == nil || session.Values["authenticated"] == false) {
		fmt.Println("invalid session")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	switch r.Method {
	case http.MethodGet:
		app.accountsView(w, r)
	case http.MethodPost:
		app.accountsPost(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) accountsView(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "indibills_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	fmt.Println(session.Values)
	if (session.Values["authenticated"] == nil) && session.Values["authenticated"] == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	user_id := session.Values["user_id"].(int64)
	fmt.Printf("%v\n", user_id)

	accounts, err := app.accounts.GetAll(user_id)

	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Accounts: %v\n", accounts)

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

func (app *application) accountsPost(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "indibills_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	}
	fmt.Println(&session.Values)
	if (session.Values["authenticated"] == nil) && session.Values["authenticated"] == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	account_name := r.PostForm.Get("name")
	account_type := r.PostForm.Get("type")
	balance, err := strconv.ParseFloat(r.PostForm.Get("balance"), 64)
	if err != nil {
		http.Error(w, "invalid value for balance; expected float", http.StatusBadRequest)
	}
	user_id := session.Values["user_id"].(int64)

	account := struct {
		Name    string  `json:"name"`
		Type    string  `json:"type"`
		Balance float64 `json:"balance"`
		User_id int64   `json:"user_id"`
	}{
		Name:    account_name,
		Type:    account_type,
		Balance: balance,
		User_id: user_id,
	}
	data, err := json.Marshal(account)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not convert account to JSON", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%v/%v/accounts", app.accounts.Endpoint, user_id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		http.Error(w, "could not create account", http.StatusBadRequest)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status: %s\n", resp.Status)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/accounts", http.StatusSeeOther)
}
