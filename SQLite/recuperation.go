package SQLite

import (
	"database/sql"
	"forum/serveur"
	"log"
)

func RecupererPublicationsClient(db *sql.DB, clientID int) []serveur.Publication {

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
        FROM publications p
        JOIN clients c ON c.id = p.client_id
        WHERE p.client_id = ?
        ORDER BY p.date DESC
    `, clientID)

	if err != nil {
		return nil
	}

	defer rows.Close()

	var publications []serveur.Publication

	for rows.Next() {

		var p serveur.Publication

		err := rows.Scan(
			&p.ID,
			&p.Titre,
			&p.Contenu,
			&p.NomClient,
			&p.Avatar,
			&p.Date,
			&p.NombreLikes,
			&p.NombreDislikes,
			&p.NombreCommentaires,
		)

		if err != nil {
			continue
		}

		publications = append(publications, p)
	}

	return publications
}

func RecupererLikesClient(db *sql.DB, clientID int) []serveur.Like {

	rows, err := db.Query(`
		SELECT publication_id, valeur
		FROM likes
		WHERE client_id = ?
	`, clientID)

	if err != nil {
		log.Println("Erreur récupération likes :", err)
		return nil
	}

	defer rows.Close()

	var likes []serveur.Like

	for rows.Next() {

		var like serveur.Like

		err := rows.Scan(
			&like.PublicationID,
			&like.Valeur,
		)

		if err != nil {
			log.Println("Erreur lecture like :", err)
			continue
		}

		likes = append(likes, like)
	}

	return likes
}

func RecupererCommentairesClient(db *sql.DB, clientID int) []serveur.Commentaire {

	rows, err := db.Query(`
        SELECT
            cmt.id,
            cmt.client_id,
            c.name,
            c.avatar,
            cmt.publication_id,
            cmt.contenu,
            cmt.date
        FROM commentaires cmt
        JOIN clients c ON c.id = cmt.client_id
        WHERE cmt.client_id = ?
        ORDER BY cmt.date DESC
    `, clientID)

	if err != nil {
		return nil
	}

	defer rows.Close()

	var commentaires []serveur.Commentaire

	for rows.Next() {

		var c serveur.Commentaire

		err := rows.Scan(
			&c.ID,
			&c.ClientID,
			&c.NomClient,
			&c.Avatar,
			&c.PublicationID,
			&c.Contenu,
			&c.Date,
		)

		if err != nil {
			continue
		}

		commentaires = append(commentaires, c)
	}

	return commentaires
}


func RecupererNomClient(db *sql.DB, clientID int) (string, error) {

	var nom string

	err := db.QueryRow(`
		SELECT name
		FROM clients
		WHERE id = ?
	`, clientID).Scan(&nom)

	if err != nil {
		return "", err
	}

	return nom, nil
}

func RecupererClientParGmail(db *sql.DB, gmail string) (int, string, error) {
    var id int
    var hash string

    err := db.QueryRow(`SELECT id, hash FROM clients WHERE gmail = ?`, gmail).Scan(&id, &hash)
    if err != nil {
        return 0, "", err
    }

    return id, hash, nil
}
func RecupererAvatar(db *sql.DB, clientID int) (string, error) {

	var avatar string

	err := db.QueryRow(`
        SELECT avatar
        FROM clients
        WHERE id = ?
    `, clientID).Scan(&avatar)

	return avatar, err
}

func RecupererPublication(db *sql.DB, publicationID, clientID int) (serveur.Publication, error) {

	var publication serveur.Publication

	err := db.QueryRow(`
        SELECT
            publications.id,
            publications.client_id,
            publications.titre,
            publications.contenu,
            clients.name,
            clients.avatar,
            publications.date,

            (
                SELECT COUNT(*)
                FROM likes
                WHERE likes.publication_id = publications.id
                AND likes.valeur = 1
            ),

            (
                SELECT COUNT(*)
                FROM likes
                WHERE likes.publication_id = publications.id
                AND likes.valeur = -1
            ),

            (
                SELECT COUNT(*)
                FROM commentaires
                WHERE commentaires.publication_id = publications.id
            )

        , COALESCE((SELECT valeur FROM likes WHERE publication_id = publications.id AND client_id = ?), 0)
 FROM publications

        JOIN clients
            ON publications.client_id = clients.id

        WHERE publications.id = ?
    `, clientID, publicationID).Scan(
		&publication.ID,
		&publication.ClientID,
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
		return publication, err
	}

	return publication, nil
}

func RecupererGmail(db *sql.DB, clientID int) (string, error) {

	var gmail string

	err := db.QueryRow(`
		SELECT gmail
		FROM clients
		WHERE id = ?
	`, clientID).Scan(&gmail)

	return gmail, err
}
