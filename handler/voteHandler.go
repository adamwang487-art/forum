package handler

import (
	"forum/serveur"
	"net/http"
	"strconv"
)

// voteHandler ajoute, remplace ou annule le vote de l'utilisateur connecté.
// Un second clic sur le même bouton annule le vote ; un clic sur l'autre bouton
// remplace le vote existant grâce à la clé primaire (client_id, publication_id).
func VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Seul un utilisateur authentifié peut voter : son identifiant est enregistré
	// avec le vote pour empêcher plusieurs votes sur la même publication.
	clientID, connecte := serveur.ClientID(r)
	if !connecte {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Les deux valeurs cachées du formulaire sont vérifiées côté serveur,
	// car elles peuvent être modifiées depuis le navigateur.
	publicationID, err := strconv.Atoi(r.FormValue("publication_id"))
	if err != nil || publicationID < 1 {
		http.Error(w, "Publication invalide", http.StatusBadRequest)
		return
	}
	valeur, err := strconv.Atoi(r.FormValue("valeur"))
	if err != nil || (valeur != 1 && valeur != -1) {
		http.Error(w, "Vote invalide", http.StatusBadRequest)
		return
	}

	// On lit le vote actuel pour distinguer l'annulation d'un changement de vote.
	var voteActuel int
	err = serveur.DB.QueryRow("SELECT valeur FROM likes WHERE client_id = ? AND publication_id = ?", clientID, publicationID).Scan(&voteActuel)
	if err == nil && voteActuel == valeur {
		// Même vote : le clic agit comme un bouton bascule et supprime le vote.
		_, err = serveur.DB.Exec("DELETE FROM likes WHERE client_id = ? AND publication_id = ?", clientID, publicationID)
	} else {
		// Nouveau vote ou changement d'avis : INSERT crée le vote, ON CONFLICT
		// met à jour la ligne existante sans autoriser un second vote.
		_, err = serveur.DB.Exec(`INSERT INTO likes (client_id, publication_id, valeur) VALUES (?, ?, ?)
			ON CONFLICT(client_id, publication_id) DO UPDATE SET valeur = excluded.valeur`, clientID, publicationID, valeur)
	}
	if err != nil {
		http.Error(w, "Impossible d'enregistrer le vote", http.StatusInternalServerError)
		return
	}

	// Après le POST, on redirige vers la publication pour afficher les compteurs mis à jour.
	retour := "/publication?id=" + strconv.Itoa(publicationID)
	if r.FormValue("retour") == "accueil" {
		retour = "/"
	}
	http.Redirect(w, r, retour, http.StatusSeeOther)
}
