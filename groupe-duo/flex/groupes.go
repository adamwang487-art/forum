package flex

import (
	"forum/SQLite"
    "forum/serveur"
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func Initialiser(database *sql.DB) error {
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS recherches_groupe (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 client_id INTEGER NOT NULL REFERENCES clients(id),
 mode TEXT NOT NULL CHECK(mode IN ('duoq', 'flex')),
 riot_id TEXT NOT NULL,
 message TEXT NOT NULL,
 date TEXT NOT NULL
 )`)
	return err
}

type annonceGroupe struct {
	RiotID, Message, Auteur, Date string
}

// Serveur gère les annonces duoQ et flexQ avec la base du forum.
type Serveur struct {
	db *sql.DB
}

func Nouveau(database *sql.DB) *Serveur {
	return &Serveur{db: database}
}

func (s *Serveur) Handler(w http.ResponseWriter, r *http.Request) {
	mode := strings.TrimPrefix(r.URL.Path, "/")
	if mode != "duoq" && mode != "flex" {
		http.NotFound(w, r)
		return
	}
	clientID, connecte := serveur.ClientID(r)
	if r.Method == http.MethodPost {
		if !connecte {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulaire invalide", http.StatusBadRequest)
			return
		}
		riotID, message := strings.TrimSpace(r.FormValue("riot_id")), strings.TrimSpace(r.FormValue("message"))
		if riotID == "" || message == "" || len([]rune(riotID)) > 80 || len([]rune(message)) > 1000 {
			http.Error(w, "Indiquez votre Riot ID (80 caractères maximum) et votre annonce (1000 caractères maximum).", http.StatusBadRequest)
			return
		}
		_, err := s.db.Exec("INSERT INTO recherches_groupe (client_id, mode, riot_id, message, date) VALUES (?, ?, ?, ?, ?)",
			clientID, mode, riotID, message, time.Now().Format("2006-01-02 15:04"))
		if err != nil {
			http.Error(w, "Impossible de publier l'annonce", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	rows, err := s.db.Query(`SELECT g.riot_id, g.message, c.name, g.date
 FROM recherches_groupe g JOIN clients c ON c.id = g.client_id
 WHERE g.mode = ? ORDER BY g.id DESC LIMIT 100`, mode)
	if err != nil {
		http.Error(w, "Impossible de charger les annonces", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	annonces := []annonceGroupe{}
	for rows.Next() {
		var annonce annonceGroupe
		if err := rows.Scan(&annonce.RiotID, &annonce.Message, &annonce.Auteur, &annonce.Date); err != nil {
			http.Error(w, "Impossible de lire les annonces", http.StatusInternalServerError)
			return
		}
		annonces = append(annonces, annonce)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Impossible de lire les annonces", http.StatusInternalServerError)
		return
	}
	nom, description := "duoQ", "Trouvez un partenaire pour jouer en groupe de 2."
	if mode == "flex" {
		nom, description = "flex", "Trouvez des coéquipiers pour former un groupe de 5."
	}
	page := struct {
		Nom, Description  string
		NomClient, Avatar string
		Connecte          bool
		Annonces          []annonceGroupe
	}{Nom: nom, Description: description, Connecte: connecte, Annonces: annonces}
	if connecte {
		page.NomClient, _ = SQLite.RecupererNomClient(s.db, clientID)
		page.Avatar, _ = SQLite.RecupererAvatar(s.db, clientID)
	}
	tmpl, err := template.ParseFiles("groupe-duo/flex/groupes.html", "templates/header.html", "templates/footer.html", "templates/espaces.html")
	var output bytes.Buffer
	if err == nil {
		err = tmpl.ExecuteTemplate(&output, "groupes.html", page)
	}
	if err != nil {
		http.Error(w, "Impossible d'afficher la page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	output.WriteTo(w)
}
