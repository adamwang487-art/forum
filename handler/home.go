package handler

import (
	"database/sql"
	"forum/SQLite"
	"forum/serveur"
	"html/template"
	"log"
	"net/http"
)

func Init(db *sql.DB) {
	serveur.DB = db
}

func Home(w http.ResponseWriter, r *http.Request) {

	clientID, connecte := serveur.ClientID(r)
	// Récupérer toutes les publications
	publications := SQLite.RecupererPublications(serveur.DB, clientID)
	nomClient := ""
	avatar := ""
	likes := []serveur.Like{}
	commentaires := []serveur.Commentaire{}

	if connecte {
		likes = SQLite.RecupererLikesClient(serveur.DB, clientID)
		commentaires = SQLite.RecupererCommentairesClient(serveur.DB, clientID)
		nomClient, _ = SQLite.RecupererNomClient(serveur.DB, clientID)
		avatar, _ = SQLite.RecupererAvatar(serveur.DB, clientID)
	}

	templates, err := template.ParseFiles(
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/espaces.html",
	)

	if err != nil {
		http.Error(w, "Erreur chargement templates", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := serveur.Data{
		Connecte:     connecte,
		ClientID:     clientID,
		NomClient:    nomClient,
		Avatar:       avatar,
		Publications: publications,
		Likes:        likes,
		Commentaires: commentaires,
	}

	err = templates.ExecuteTemplate(w, "index.html", data)

	if err != nil {
		http.Error(w, "Erreur affichage page", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
