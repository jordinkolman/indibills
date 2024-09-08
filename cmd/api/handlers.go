package main

import (
	"encoding/json"
	"fmt"
	"indibills/internal/data"
	"net/http"
	"strconv"
	"strings"
)

/*
	TODO: Implement the following handlers
	get & create transactions
	get & create assets
	get & create liabilities
	get & create budget items
*/

type UserList []data.User

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var versionString = fmt.Sprintf("%v.%v.%v", data.VERSION, data.SUBVERSION, data.PATCH)

	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     versionString,
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func (app *application) getCreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		users, err := app.models.Users.GetAll()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error 1", http.StatusInternalServerError)
			return
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil); err != nil {
			http.Error(w, "error 2", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		var input struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user := &data.User{
			Email:     input.Email,
			Password:  input.Password,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		}

		err = app.models.Users.Insert(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v%v/users/%d", data.VERSION, user.Id))

		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}

// TODO
func (app *application) getCreateAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user_id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, fmt.Sprintf("v%v/accounts/users/", data.VERSION)), 10, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		users, err := app.models.Accounts.GetAccountById(user_id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error 1", http.StatusInternalServerError)
			return
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil); err != nil {
			http.Error(w, "error 2", http.StatusInternalServerError)
			return
		}
	}

	/* if r.Method == http.MethodPost {
		var input struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user := &data.User{
			Email:     input.Email,
			Password:  input.Password,
			FirstName: input.FirstName,
			LastName:  input.LastName,
		}

		err = app.models.Users.Insert(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v%v/users/%d", data.VERSION, user.Id))

		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}*/

}
