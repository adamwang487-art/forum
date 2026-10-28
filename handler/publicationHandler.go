package handler

import (
	"forum/SQLite"
	"forum/serveur"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PublicationHandler(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	
	clientID, connecte := serveur.ClientID(r)
	nomClient := ""
	avatar := ""
	if connecte {
		nomClient, _ = SQLite.RecupererNomClient(serveur.DB, clientID)
		avatar, _ = SQLite.RecupererAvatar(serveur.DB, clientID)
	}

	if err != nil {
		http.Error(w, "Publication invalide", http.StatusBadRequest)
		return
	}

	publication, err := SQLite.RecupererPublication(serveur.DB, id, clientID)

	if err != nil {
		http.Error(w, "Publication introuvable", http.StatusNotFound)
		return
	}

	commentaires, err := serveur.RecupererCommentaires(serveur.DB, id)

	if err != nil {
		http.Error(w, "Erreur récupération commentaires", http.StatusInternalServerError)
		return
	}

	data := serveur.Data{
		Connecte:     connecte,
		ClientID:     clientID,
		NomClient:    nomClient,
		Avatar:       avatar,
		Publication:  publication,
		Commentaires: commentaires,
	}

	tmpl, err := template.ParseFiles(
		"templates/publication.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/espaces.html",
	)

	if err != nil {
		log.Println("Erreur ParseFiles :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "publication.html", data)

	if err != nil {
		log.Println("Erreur template publication :", err)
		return
	}
}

func CreationPublicationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	clientID, connecte := serveur.ClientID(r)

	// Seuls les utilisateurs connectés
	// peuvent publier.
	if !connecte {
		http.Error(
			w,
			"Vous devez être connecté pour publier",
			http.StatusUnauthorized,
		)
		return
	}

	switch r.Method {

	case http.MethodGet:

		categories, err :=
			serveur.RecupererCategories(serveur.DB)

		if err != nil {
			log.Println(err)

			http.Error(
				w,
				"Erreur serveur",
				http.StatusInternalServerError,
			)
			return
		}

		tmpl, err :=
			template.ParseFiles(
				"templates/create-post.html",
			)

		if err != nil {
			log.Println(err)

			http.Error(
				w,
				"Erreur serveur",
				http.StatusInternalServerError,
			)
			return
		}

		err = tmpl.Execute(
			w,
			categories,
		)

		if err != nil {
			log.Println(err)
		}

	case http.MethodPost:

		if err := r.ParseForm(); err != nil {
			http.Error(
				w,
				"Requête invalide",
				http.StatusBadRequest,
			)
			return
		}

		titre := strings.TrimSpace(
			r.FormValue("titre"),
		)

		if titre == "" {
			http.Error(
				w,
				"Le titre est obligatoire",
				http.StatusBadRequest,
			)
			return
		}
		contenu :=
			strings.TrimSpace(
				r.FormValue("contenu"),
			)

		if contenu == "" {
			http.Error(
				w,
				"La publication est vide",
				http.StatusBadRequest,
			)
			return
		}

		categories :=
			r.Form["categories"]

		if len(categories) == 0 {
			http.Error(
				w,
				"Choisissez au moins une catégorie",
				http.StatusBadRequest,
			)
			return
		}

		date :=
			time.Now().Format(
				"2006-01-02 15:04:05",
			)

		publicationID, err :=
			SQLite.AjoutePublication(
				serveur.DB,
				clientID,
				titre,
				contenu,
				date,
			)

		if err != nil {
			log.Println(err)

			http.Error(
				w,
				"Impossible de créer la publication",
				http.StatusInternalServerError,
			)
			return
		}

		for _, categorie := range categories {

			categorieID, err :=
				strconv.Atoi(categorie)

			if err != nil {
				continue
			}

			err =
				SQLite.AjouteCategoriePublication(
					serveur.DB,
					publicationID,
					categorieID,
				)

			if err != nil {
				log.Println(err)

				http.Error(
					w,
					"Erreur ajout catégorie",
					http.StatusInternalServerError,
				)
				return
			}
		}

		http.Redirect(
			w,
			r,
			"/",
			http.StatusSeeOther,
		)

	default:

		http.Error(
			w,
			"Méthode non autorisée",
			http.StatusMethodNotAllowed,
		)
	}
}
