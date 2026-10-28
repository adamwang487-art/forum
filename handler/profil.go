package handler

import (
	"forum/SQLite"
	"forum/serveur"
	"html/template"
	"net/http"
	"strings"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {

	clientID, connecte := serveur.ClientID(r)
	nomClient, _ := SQLite.RecupererNomClient(serveur.DB, clientID)
	gmail, _ := SQLite.RecupererGmail(serveur.DB, clientID)
	avatar, _ := SQLite.RecupererAvatar(serveur.DB, clientID)
	publications := SQLite.RecupererPublicationsClient(serveur.DB, clientID)
	commentaires := SQLite.RecupererCommentairesClient(serveur.DB, clientID)
	data := serveur.Data{
		ClientID:     clientID,
		Connecte:     connecte,
		NomClient:    nomClient,
		Gmail:        gmail,
		Publications: publications,
		Commentaires: commentaires,
		Avatar:       avatar,
		Avatars:      serveur.ListeAvatars(),
	}

	tmpl, err := template.New("Profil.html").Funcs(template.FuncMap{"fondAvatar": serveur.FondAvatar}).ParseFiles("templates/Profil.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ModifierProfilHandler(w http.ResponseWriter, r *http.Request) {

	clientID, connecte := serveur.ClientID(r)

	if !connecte {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// =========================
	// AFFICHAGE DE LA PAGE
	// =========================

	if r.Method == http.MethodGet {

		nomClient, _ := SQLite.RecupererNomClient(serveur.DB, clientID)
		gmail, _ := SQLite.RecupererGmail(serveur.DB, clientID)

		data := serveur.Data{
			ClientID:  clientID,
			Connecte:  connecte,
			NomClient: nomClient,
			Gmail:     gmail,
			Avatars:   serveur.ListeAvatars(),
		}

		tmpl, err := template.ParseFiles(
			"templates/modifier-profil.html",
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	// =========================
	// MODIFICATION
	// =========================

	if r.Method == http.MethodPost {

		nom := strings.TrimSpace(r.FormValue("name"))
		gmail := strings.TrimSpace(r.FormValue("gmail"))
		motdepasse := r.FormValue("motdepasse")

		if nom == "" || gmail == "" {
			http.Error(
				w,
				"Tous les champs obligatoires doivent être remplis",
				http.StatusBadRequest,
			)
			return
		}

		err := serveur.ModifierProfil(
			serveur.DB,
			clientID,
			nom,
			gmail,
			motdepasse,
		)

		if err != nil {
			http.Error(
				w,
				"Erreur lors de la modification du profil",
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(w, r, "/profil", http.StatusSeeOther)
		return
	}

	http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
}

func ModifierAvatarHandler(w http.ResponseWriter, r *http.Request) {

	clientID, connecte := serveur.ClientID(r)

	if !connecte {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// =========================
	// AFFICHAGE
	// =========================

	if r.Method == http.MethodGet {

		avatar, _ := SQLite.RecupererAvatar(
			serveur.DB,
			clientID,
		)

		data := serveur.Data{
			ClientID: clientID,
			Connecte: connecte,
			Avatar:   avatar,
			Avatars:  serveur.ListeAvatars(),
		}

		tmpl, err := template.New("avatar.html").Funcs(template.FuncMap{"fondAvatar": serveur.FondAvatar}).ParseFiles(
			"templates/avatar.html",
		)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		return
	}

	// =========================
	// MODIFICATION
	// =========================

	if r.Method == http.MethodPost {

		avatar := r.FormValue("avatar")

		if avatar == "" {
			http.Error(
				w,
				"Avatar invalide",
				http.StatusBadRequest,
			)
			return
		}

		err := SQLite.ModifierAvatar(
			serveur.DB,
			clientID,
			avatar,
		)

		if err != nil {
			http.Error(
				w,
				"Erreur lors de la modification de l'avatar",
				http.StatusInternalServerError,
			)
			return
		}

		http.Redirect(
			w,
			r,
			"/profil",
			http.StatusSeeOther,
		)

		return
	}

	http.Error(
		w,
		"Méthode non autorisée",
		http.StatusMethodNotAllowed,
	)
}
