package main

/*
	These are the handlers for any HTTP requests to the API.

	!ENDPOINTS:
	/healthcheck:                 GET                     returns the current status of the API, including availability, environment, and version
	! This endpoint returns information for every user. Should only be reachable by site admin or authorized scripts (TODO)
	TODO /users                        GET                     endpoint for viewing a list of all users. Should only be accessible by site adminX
	! all endpoints below require an authenticated user and use the passed in user_id from the session variables
	TODO /users                        PUT/PATCH, DELETE       endpoint for updating or deleting a specific user record from the database
	/accounts                     GET, POST               endpoint for adding an account record to the database, or retrieving a list of all accounts for the specified user
	TODO /transactions                 GET, POST               endpoint for adding or retrieving transaction records for a specific user
	TODO /assets                       GET, POST               endpoint for adding or retrieving asset (property) records for a specific user
	TODO /liabilities                  GET, POST               endpoint for adding or retrieving liability (debt) records for a specific user
	TODO /goals                        GET, POST               endpoint for adding or retrieving all budget goal item records for a specific user
	/users/{email}                GET, POST               retrieving and creating a specific user record from the database. Email passed in via path parameters
		TODO - consider implementation of readinglist and decide if /users/ or /users is best for GET requests
	? /accounts/                    GET, PUT/PATCH, DELETE       modifying and deleting a specific account record from the database
	? /transactions/                GET, PUT/PATCH, DELETE       modifying and deleting a specific transaction record from the database
	? /assets/                      GET, PUT/PATCH, DELETE       modifying and deleting a specific account asset from the database
	? /liabilities/                 GET, PUT/PATCH, DELETE       modifying and deleting a specific account liability from the database
	? /goals/                       GET, PUT/PATCH, DELETE       modifying and deleting a specific budget goal item record from the database



*/

import (
	"encoding/json"
	"fmt"
	"indibills/internal/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
	TODO: Implement the following handlers
	get & create transactions
	get & create assets
	get & create liabilities
	get & create budget goals
*/

// an endpoint that can be pinged to check API status. Returns status: available, and the current environment and version
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// a string of the combined version, subversion and patch (i.e., 3.5.0)
	var versionString = fmt.Sprintf("%v.%v.%v", data.VERSION, data.SUBVERSION, data.PATCH)
	// the data to be returned on query
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     versionString,
	}
	// convert the data map into a JSON object for transmission
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// add a newline for readability
	js = append(js, '\n')
	// set the Content-Type header so the requesting party knows to expect JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}


//! Associated Endpoint: {api_path}/users/{email}    Methods: GET, POST
func (app *application) getCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	//! GET
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		// retrieve the email from the URL
		email, ok := vars["email"]
		if !ok {
			log.Print("email missing from path parameters")
		}

		user, err := app.models.Users.Get(email)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error 1", http.StatusInternalServerError)
			return
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil); err != nil {
			http.Error(w, "error 2", http.StatusInternalServerError)
			return
		}
	}
	//! POST
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
		headers.Set("Location", fmt.Sprintf("v%v/users/%v", data.VERSION, user.Email))

		err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}
//! Associated endpoint: /accounts
func (app *application) getCreateAccountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		user_id, err := strconv.ParseInt(vars["user_id"], 10, 64)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		accounts, err := app.models.Accounts.GetAccountsByUserId(user_id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error 1", http.StatusInternalServerError)
			return
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"accounts": accounts}, nil); err != nil {
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
