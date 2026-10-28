package serveur

import "database/sql"

func RecupererCommentaires(db *sql.DB, publicationID int) ([]Commentaire, error) {

	rows, err := db.Query(`
        SELECT
            commentaires.id,
            commentaires.client_id,
            commentaires.publication_id,
            commentaires.contenu,
            commentaires.date,
            clients.name,
			clients.avatar
        FROM commentaires
        JOIN clients
            ON commentaires.client_id = clients.id
        WHERE commentaires.publication_id = ?
        ORDER BY commentaires.date ASC
    `, publicationID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var commentaires []Commentaire

	for rows.Next() {

		var commentaire Commentaire

		err := rows.Scan(
			&commentaire.ID,
			&commentaire.ClientID,
			&commentaire.PublicationID,
			&commentaire.Contenu,
			&commentaire.Date,
			&commentaire.NomClient,
			&commentaire.Avatar,
		)

		if err != nil {
			return nil, err
		}

		commentaires = append(commentaires, commentaire)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentaires, nil
}

func SupprimerCommentaire(db *sql.DB, commentaireID int, clientID int) error {
	_, err := db.Exec(`
		DELETE FROM commentaires
		WHERE id = ? AND client_id = ?
	`, commentaireID, clientID)

	return err
}
