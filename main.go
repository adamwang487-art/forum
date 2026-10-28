package main

import (
    "database/sql"
    "forum/SQLite"
    "forum/groupe-duo/flex"
    "forum/handler"
    "forum/notepatch"
    "forum/serveur"
    "log"
    "net/http"

    _ "modernc.org/sqlite"
)
var db *sql.DB

func main() {
	db = SQLite.Ouvrirdata()
	defer db.Close()
	if err := flex.Initialiser(db); err != nil {
		log.Fatal(err)
	}
	groupes := flex.Nouveau(db)
	handler.Init(db)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/scores-et-matchs/", http.StripPrefix("/scores-et-matchs/", http.FileServer(http.Dir("scores et matchs"))))
	http.HandleFunc("/groupes/groupes.css", flex.Styles)
	http.HandleFunc("/duoq", groupes.Handler)
	http.HandleFunc("/flex", groupes.Handler)
	http.HandleFunc("/patch-notes", notepatch.Handler)
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/register/retour", handler.RegisterRetour)
	http.HandleFunc("/register", handler.Register)
	http.Handle("/create-post", serveur.RequireAuth(http.HandlerFunc(handler.CreationPublicationHandler)))
	http.Handle("/commentaire", serveur.RequireAuth(http.HandlerFunc(handler.CommentaireHandler)))
	http.HandleFunc("/publication", handler.PublicationHandler)
	http.HandleFunc("/profil", handler.ProfilHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.Handle("/profil/modifier", serveur.RequireAuth(http.HandlerFunc(handler.ModifierProfilHandler)))
	http.Handle("/profil/avatar", serveur.RequireAuth(http.HandlerFunc(handler.ModifierAvatarHandler)))
	http.Handle("/vote", serveur.RequireAuth(http.HandlerFunc(handler.VoteHandler)))

	log.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
