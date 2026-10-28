package handler

import (
	"forum/SQLite"
	"forum/serveur"
	"net/http"
	"strconv"
	"time"
)

func CommentaireHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	publicationID, err := strconv.Atoi(r.FormValue("publication_id"))
	if err != nil {
		http.Error(w, "Publication invalide", http.StatusBadRequest)
		return
	}

	clientID, connecte := serveur.ClientID(r)

	if !connecte {
		http.Error(w, "Vous devez être connecté", http.StatusUnauthorized)
		return
	}

	contenu := r.FormValue("contenu")

	if contenu == "" {
		http.Error(w, "Le commentaire est vide", http.StatusBadRequest)
		return
	}

	date := time.Now().Format("2006-01-02 15:04:05")

	SQLite.AjouteCommentaire(serveur.DB, clientID, publicationID, contenu, date)

	http.Redirect(w, r, "http://localhost:8080/publication?id="+strconv.Itoa(publicationID), http.StatusSeeOther)
}
