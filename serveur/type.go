package serveur

import "database/sql"

type Publication struct {
	ID                 int
	ClientID           int
	Titre              string
	Contenu            string
	NomClient          string
	Avatar             string
	Date               string
	NombreLikes        int
	NombreDislikes     int
	NombreCommentaires int
	VoteUtilisateur    int
}

type Like struct {
	PublicationID int
	Valeur        int
}

type Commentaire struct {
	ID            int
	ClientID      int
	NomClient     string
	Avatar        string
	PublicationID int
	Contenu       string
	Date          string
}

type Data struct {
	Connecte  bool
	ClientID  int
	Gmail     string
	NomClient string

	Publications []Publication
	Publication  Publication

	Likes        []Like
	Commentaires []Commentaire

	Avatar  string
	Avatars []string

	Erreur string
}

type Categorie struct {
	ID  int
	Nom string
}

var DB *sql.DB
