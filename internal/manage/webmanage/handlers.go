package webmanage

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Kwynto/GracefulDB/internal/base/basicsystem/gauth"
	"github.com/Kwynto/gosession"
)

// Handler after authorization
func homeDefault(w http.ResponseWriter, r *http.Request, login string) {
	// -
	ts, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		slog.Debug("Internal Server Error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Authorization Handler
func homeAuth(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			slog.Debug("Bad request", slog.String("err", err.Error()))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")
		isAuth := gauth.CheckUser(username, password)
		if isAuth {
			sesID := gosession.Start(&w, r)
			sesID.Set("auth", username)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		ts, err := template.ParseFiles("./ui/html/auth.html")
		if err != nil {
			slog.Debug("Internal Server Error", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			slog.Debug("Internal Server Error", slog.String("err", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// Handler for the main route
func home(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	sesID := gosession.Start(&w, r)
	auth := sesID.Get("auth")
	if auth == nil {
		homeAuth(w, r)
	} else {
		login := fmt.Sprint(auth)
		homeDefault(w, r, login)
	}
}

// Exit handler
func logout(w http.ResponseWriter, r *http.Request) {
	// This function is complete
	sesID := gosession.Start(&w, r)
	sesID.Remove("auth")
	http.Redirect(w, r, "/", http.StatusFound)
}