package main

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	MsgMissingRequired    = "Please provide all the required values"
	MsgUserNameExists     = "Your desired User Name already exists."
	MsgEmailExists        = "A User Name already exists for this Email Address."
	MsgPasswordsDifferent = "Password do not match."
)

// RegisterHandler handles /register requests.
func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := app.tmpls.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		app.registerPost(w, r)
	}
}

// registerPost is called for the POST method of the RegisterHandler.
func (app *App) registerPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("userName"))
	fullName := strings.TrimSpace(r.PostFormValue("fullName"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password1 := strings.TrimSpace(r.PostFormValue("password1"))
	password2 := strings.TrimSpace(r.PostFormValue("password2"))

	// check for missing values
	// redundant given client side required fields, but good practice
	if userName == "" || password1 == "" || password2 == "" || fullName == "" || email == "" {
		msg := MsgMissingRequired
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// check that userName doesn't already exist
	userExists, err := UserExists(app.db, userName)
	if err != nil {
		log.Printf("error in UserExists for %q: %v", userName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Printf("userName %q already exists", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", MsgUserNameExists)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// check that email doesn't already exist
	emailExists, err := EmailExists(app.db, email)
	if err != nil {
		log.Printf("error in EmailExists for %q: %v", email, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if emailExists {
		log.Printf("email %q already exists", email)
		err := app.tmpls.ExecuteTemplate(w, "register.html", MsgEmailExists)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		msg := "Cannot hash password"
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// store the user and hashed password
	_, err = app.db.Exec("INSERT INTO users(username, hashedPassword, fullName, email) VALUES (?, ?, ?, ?)",
		userName, hashedPassword, fullName, email)
	if err != nil {
		msg := "Unable to register user"
		log.Println(msg, "for", userName, err)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// register successful
	log.Printf("Username %q registered", userName)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
