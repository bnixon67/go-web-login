package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrNoSuchUser      = errors.New("no such user")
	ErrInternalFailure = errors.New("login failed due to internal error")
)

// LoginHandler handles /login requests.
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := ExecTemplateOrError(app.tmpls, w, "login.html", nil)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}

	case http.MethodPost:
		app.loginPost(w, r)
	}
}

const (
	MsgMissingUserNameAndPassword = "Missing username and password"
	MsgMissingUserName            = "Missing username"
	MsgMissingPassword            = "Missing password"
	MsgLoginFailed                = "Login Failed"
)

// loginPost is called for the POST method of the LoginHandler.
func (app *App) loginPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))

	// check for missing values
	var msg string
	switch {
	case userName == "" && password == "":
		msg = MsgMissingUserNameAndPassword
	case userName == "":
		msg = MsgMissingUserName
	case password == "":
		msg = MsgMissingPassword
	}
	if msg != "" {
		log.Println(msg)
		err := ExecTemplateOrError(app.tmpls, w, "login.html", msg)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// attempt to login the given userName with the given password
	token, err := app.LoginUser(userName, password)
	if err != nil {
		err := ExecTemplateOrError(app.tmpls, w, "login.html", MsgLoginFailed)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// login successful, so create a cookie for the session Token
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   token.Value,
		Expires: token.Expires,
	})
	log.Printf("valid login for %q", userName)

	http.Redirect(w, r, "/hello", http.StatusSeeOther)
}

// LoginUser returns a session Token if userName and password is correct.
func (app *App) LoginUser(userName, password string) (Token, error) {
	err := CompareUserPassword(app.db, userName, password)
	if err != nil {
		log.Printf("invalid password for %q: %v", userName, err)
		return Token{}, err
	}

	// create and save a new session token
	token, err := SaveNewToken(app.db, "session", userName, 32, app.config.SessionExpiresHours)
	if err != nil {
		log.Printf("unable to save session token: %v", err)
		return Token{}, ErrInternalFailure
	}

	return token, nil
}
