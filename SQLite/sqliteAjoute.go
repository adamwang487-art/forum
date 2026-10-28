package SQLite

import (
	"database/sql"
	"log"
)

func Ajouteclients(db *sql.DB, gmail string, name string, hash string) error {
    _, err := db.Exec(`INSERT INTO clients (gmail, name, hash) VALUES (?, ?, ?)`, gmail, name, hash)
    return err
}

func AjoutePublication(db *sql.DB, clientID int, titre, contenu, date string) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO publications (
			client_id,
			titre,
			contenu,
			date
		)
		VALUES (?, ?, ?, ?)
	`, clientID, titre, contenu, date)

	if err != nil {
		return 0, err
	}

	publicationID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return publicationID, nil
}

func AjouteLike(db *sql.DB, clientID int, publicationsID int) {
	_, err := db.Exec(`INSERT INTO likes (client_id, publication_id) VALUES (?, ?)`, clientID, publicationsID)
	if err != nil {
		log.Fatal(err)
	}
}

func AjouteCommentaire(db *sql.DB, clientID, publicationID int, contenu, date string) error {
	_, err := db.Exec(`
		INSERT INTO commentaires
		(client_id, publication_id, contenu, date)
		VALUES (?, ?, ?, ?)
	`, clientID, publicationID, contenu, date)

	if err != nil {
		log.Println("Erreur ajout commentaire :", err)
		return err
	}

	return nil
}

func AjouteCategoriePublication(db *sql.DB, publicationID int64, categorieID int) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO publication_categories (
			publication_id,
			categorie_id
		)
		VALUES (?, ?)
	`, publicationID, categorieID)

	return err
}
