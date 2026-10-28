package main

import (
	"forum/SQLite"
	"forum/groupe-duo/flex"
	"forum/handler"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFusionPagesEtVotes(t *testing.T) {
	// Ouvre une copie de la base existante pour tester la migration sans la modifier.
	content, err := os.ReadFile("chat.db")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "chat.db"), content, 0600); err != nil {
		t.Fatal(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	database := SQLite.Ouvrirdata()
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	oldDB := db
	db = database
	handler.Init(db)
	defer func() { db = oldDB; handler.Init(oldDB) }()
	if err := flex.Initialiser(db); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO clients (id, gmail, name, motdepasse) VALUES (99999, 'fusion@example.test', 'Fusion', 'test');
 INSERT INTO publications (id, client_id, titre, contenu, date) VALUES (99999, 99999, 'Fusion', 'Test fusion', '2026-09-28');`)
	if err != nil {
		t.Fatal(err)
	}
	groupes := flex.Nouveau(db)
	session := httptest.NewRecorder()
	SQLite.CreerSession(session, 99999)
	cookie := session.Result().Cookies()[0]
	for _, page := range []struct {
		path    string
		handler http.HandlerFunc
	}{
		{"/", handler.Home}, {"/publication?id=99999", handler.PublicationHandler},
		{"/duoq", groupes.Handler}, {"/flex", groupes.Handler}, {"/profil", handler.ProfilHandler},
	} {
		req := httptest.NewRequest("GET", page.path, nil)
		req.AddCookie(cookie)
		res := httptest.NewRecorder()
		page.handler(res, req)
		if res.Code != 200 || !strings.Contains(res.Body.String(), "</html>") {
			t.Fatalf("%s: %d %s", page.path, res.Code, res.Body.String())
		}
	}
	for _, value := range []string{"1", "-1", "-1"} {
		req := httptest.NewRequest("POST", "/vote", strings.NewReader("publication_id=99999&retour=accueil&valeur="+value))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie)
		res := httptest.NewRecorder()
		handler.VoteHandler(res, req)
		if res.Code != 303 || res.Header().Get("Location") != "/" {
			t.Fatalf("Vote: %d %s", res.Code, res.Body.String())
		}
		publication, err := SQLite.RecupererPublication(db, 99999, 99999)
		if err != nil {
			t.Fatal(err)
		}
		if value == "1" && (publication.VoteUtilisateur != 1 || publication.NombreLikes != 1) {
			t.Fatal(publication)
		}
	}
	publication, err := SQLite.RecupererPublication(db, 99999, 99999)
	if err != nil || publication.VoteUtilisateur != 0 || publication.NombreDislikes != 0 {
		t.Fatalf("Annulation: %+v %v", publication, err)
	}
}
