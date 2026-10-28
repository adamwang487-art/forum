package serveur

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

func ModifierProfil(db *sql.DB, clientID int, nom string, gmail string, motdepasse string) error {

	if motdepasse != "" {

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(motdepasse),
			bcrypt.DefaultCost,
		)

		if err != nil {
			return err
		}

		_, err = db.Exec(`
			UPDATE clients
			SET name = ?,
			    gmail = ?,
			    hash = ?
			WHERE id = ?
		`, nom, gmail, string(hash), clientID)

		return err
	}

	// Aucun nouveau mot de passe :
	// on garde l'ancien.
	_, err := db.Exec(`
		UPDATE clients
		SET name = ?,
		    gmail = ?
		WHERE id = ?
	`, nom, gmail, clientID)

	return err
}
