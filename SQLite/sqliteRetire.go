package SQLite

import (
	"database/sql"
	"log"
)

func RetireClient(db *sql.DB, clientID int) {
	// Supprimer les commentaires du client
	_, err := db.Exec(`DELETE FROM commentaires WHERE client_id = ?`, clientID)
	if err != nil {
		log.Println("Erreur suppression commentaires :", err)
		return
	}

	// Supprimer les likes du client
	_, err = db.Exec(`DELETE FROM likes WHERE client_id = ?`, clientID)
	if err != nil {
		log.Println("Erreur suppression likes :", err)
		return
	}

	// Supprimer les commentaires sur les publications du client
	_, err = db.Exec(`DELETE FROM commentaires WHERE publication_id IN ( SELECT id FROM publications WHERE client_id = ? )`, clientID)
	if err != nil {
		log.Println("Erreur suppression commentaires publications :", err)
		return
	}

	// Supprimer les likes sur les publications du client
	_, err = db.Exec(`DELETE FROM likes WHERE publication_id IN ( SELECT id FROM publications WHERE client_id = ? )`, clientID)
	if err != nil {
		log.Println("Erreur suppression likes publications :", err)
		return
	}

	// Supprimer les publications du client
	_, err = db.Exec(`DELETE FROM publications WHERE client_id = ?`, clientID)
	if err != nil {
		log.Println("Erreur suppression publications :", err)
		return
	}

	// Enfin supprimer le client
	_, err = db.Exec(`DELETE FROM clients WHERE id = ?`, clientID)
	if err != nil {
		log.Println("Erreur suppression client :", err)
		return
	}
}

func RetirePublication(db *sql.DB, publicationID int) {
	// Supprimer les commentaires liés
	_, err := db.Exec(`DELETE FROM commentaires WHERE publication_id = ?`, publicationID)
	if err != nil {
		log.Println("Erreur suppression commentaires :", err)
		return
	}

	// Supprimer les likes liés
	_, err = db.Exec(`DELETE FROM likes WHERE publication_id = ?`, publicationID)
	if err != nil {
		log.Println("Erreur suppression likes :", err)
		return
	}

	// Supprimer la publication
	_, err = db.Exec(`DELETE FROM publications WHERE id = ?`, publicationID)
	if err != nil {
		log.Println("Erreur suppression publication :", err)
		return
	}
}

func RetireCommentaire(db *sql.DB, commentaireID int) {
	_, err := db.Exec(`DELETE FROM commentaires WHERE id = ?`, commentaireID)
	if err != nil {
		log.Println("Erreur suppression commentaire :", err)
		return
	}
}

func RetireLike(db *sql.DB, clientID, publicationID int) {
	_, err := db.Exec(`DELETE FROM likes WHERE client_id = ? AND publication_id = ?`, clientID, publicationID)

	if err != nil {
		log.Println("Erreur suppression like :", err)
	}
}
