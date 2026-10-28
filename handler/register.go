package handler

import (
	"forum/SQLite"
	"forum/serveur"
	"net/http"
	"strings"
	"text/template"
)

func Register(w http.ResponseWriter, r *http.Request) {

	template.ParseFiles("templates/header.html")

	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/register.html")
		return
	}

	if r.Method == "POST" {
		gmail := strings.TrimSpace(r.FormValue("email"))
		name := strings.TrimSpace(r.FormValue("name"))
		motdepasse := strings.TrimSpace(r.FormValue("password"))
		if gmail == "" || name == "" || motdepasse == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}
		hash, err := serveur.HashPassword(motdepasse)
		if err != nil {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		err = SQLite.Ajouteclients(serveur.DB, gmail, name, hash)

		if err != nil {
			http.Error(w, "Cette adresse mail ou ce nom d'utilisateur est déjà utilisée.", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func RegisterRetour(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
