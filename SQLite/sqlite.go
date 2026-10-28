package SQLite

import (
	"database/sql"
	"fmt"
	"forum/serveur"
	"log"

	_ "modernc.org/sqlite"
)

func Ouvrirdata() *sql.DB {
	// Ouvre la base SQLite
	db, err := sql.Open("sqlite", "chat.db")
	if err != nil {
		log.Fatal(err)
	}

	// Crée la table pour les clients
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			gmail TEXT UNIQUE NOT NULL,
			name TEXT UNIQUE NOT NULL,
			hash TEXT NOT NULL,
			avatar TEXT NOT NULL DEFAULT 'Ahri'
		)
			
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Crée la table pour les publications
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS publications (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	client_id INTEGER NOT NULL,
    	titre TEXT NOT NULL,
    	contenu TEXT NOT NULL,
    	date TEXT NOT NULL,
    	FOREIGN KEY (client_id) REFERENCES clients(id)
	)
	`)

	// Crée la table des catégories
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nom TEXT NOT NULL UNIQUE
	)
`)
	if err != nil {
		log.Fatal(err)
	}

	// Crée la table de liaison entre publications et catégories
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS publication_categories (
		publication_id INTEGER NOT NULL,
		categorie_id INTEGER NOT NULL,

		PRIMARY KEY (publication_id, categorie_id),

		FOREIGN KEY (publication_id)
			REFERENCES publications(id)
			ON DELETE CASCADE,

		FOREIGN KEY (categorie_id)
			REFERENCES categories(id)
			ON DELETE CASCADE
	)
`)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	categories := []string{
		"Général",
		"Technologie",
		"Jeux",
		"Sport",
		"Musique",
		"Cinéma",
	}

	for _, categorie := range categories {
		_, err = db.Exec(`
		INSERT OR IGNORE INTO categories (nom)
		VALUES (?)
	`, categorie)

		if err != nil {
			log.Println("Erreur ajout catégorie :", err)
		}
	}
	// Crée la table pour les likes
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        client_id INTEGER NOT NULL,
        publication_id INTEGER NOT NULL,
		valeur INTEGER NOT NULL,
        FOREIGN KEY (client_id) REFERENCES clients(id),
        FOREIGN KEY (publication_id) REFERENCES publications(id),
        UNIQUE(client_id, publication_id)
    )
	`)
	if err != nil {
		log.Fatal(err)
	}
	// Crée la table pour les commentaires
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS commentaires (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        client_id INTEGER NOT NULL,
        publication_id INTEGER NOT NULL,
        contenu TEXT NOT NULL,
        date TEXT NOT NULL,

        FOREIGN KEY (client_id) REFERENCES clients(id),
        FOREIGN KEY (publication_id) REFERENCES publications(id))`)
	if err != nil {
		log.Fatal(err)
	}
	if err := migrerDonnees(db); err != nil {
		log.Fatal(err)
	}
	return db
}

func AfficherClients(db *sql.DB) {
	fmt.Println("Clients")

	rows, err := db.Query("SELECT id, gmail, name, hash FROM clients")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
    var id int
    var gmail string
    var name string
    var hash string

    err := rows.Scan(&id, &gmail, &name, &hash)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %d | Gmail : %s | Nom: %s | Hash: %s\n", id, gmail, name, hash)
}
}

func VerificationClient(db *sql.DB, name string) bool {
	var id int
	err := db.QueryRow(`SELECT id FROM clients WHERE name = ?`, name).Scan(&id)

	if err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

func RecupererPublications(db *sql.DB, clientID int) []serveur.Publication {

	rows, err := db.Query(`
		SELECT
			p.id,
			p.titre,
			p.contenu,
			c.name,
			c.avatar,
			p.date,

			(
				SELECT COUNT(*)
				FROM likes l
				WHERE l.publication_id = p.id
				AND l.valeur = 1
			),

			(
				SELECT COUNT(*)
				FROM likes l
				WHERE l.publication_id = p.id
				AND l.valeur = -1
			),

			(
				SELECT COUNT(*)
				FROM commentaires cmt
				WHERE cmt.publication_id = p.id
			)

		, COALESCE((SELECT valeur FROM likes WHERE publication_id = p.id AND client_id = ?), 0)
 FROM publications p

		JOIN clients c
			ON p.client_id = c.id

		ORDER BY p.id DESC
	`, clientID)

	if err != nil {
		log.Println("Erreur récupération publications :", err)
		return nil
	}

	defer rows.Close()

	var publications []serveur.Publication

	for rows.Next() {

		var publication serveur.Publication

		err := rows.Scan(
			&publication.ID,
			&publication.Titre,
			&publication.Contenu,
			&publication.NomClient,
			&publication.Avatar,
			&publication.Date,
			&publication.NombreLikes,
			&publication.NombreDislikes,
			&publication.NombreCommentaires,
			&publication.VoteUtilisateur,
		)

		if err != nil {
			log.Println("Erreur lecture publication :", err)
			continue
		}

		publications = append(publications, publication)
	}

	return publications
}

func ModifierAvatar(db *sql.DB, clientID int, avatar string) error {
	_, err := db.Exec(`
		UPDATE clients
		SET avatar = ?
		WHERE id = ?
	`, avatar, clientID)

	return err
}
