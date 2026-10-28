package handler

import (
	"net/http"
	"strconv"
	"strings"

	"forum/SQLite"
	"forum/serveur"
)

func Login(w http.ResponseWriter, r *http.Request) {

	// 1. Afficher la page de connexion
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	// 2. Traiter le formulaire envoyé
	if r.Method == "POST" {

		// 2a. Lire les champs
		gmail := strings.TrimSpace(r.FormValue("email"))
		motdepasse := strings.TrimSpace(r.FormValue("password"))

		if gmail == "" || motdepasse == "" {
			http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
			return
		}

		// 2b. Vérifier l'identité
		clientID, hash, err := SQLite.RecupererClientParGmail(serveur.DB, gmail)
		if err != nil || !serveur.CheckPassword(motdepasse, hash) {
			http.Error(w, "Mot de passe ou adresse mail incorrect.", http.StatusUnauthorized)
			return
		}

		// 2c. Créer la session (le "bracelet" côté serveur)
		sessionID := serveur.CreateSession(strconv.Itoa(clientID))
		if sessionID == "" {
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}

		// 2d. Donner le bracelet au navigateur
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(serveur.SessionDuration.Seconds()),
		})

		// 2e. Aller à l'accueil
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// 3. Autre méthode HTTP : refusée
	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}
