package flex

import (
	"database/sql"
	"forum/SQLite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestGroupes(t *testing.T) {
	t.Chdir("../..")
	database, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	defer database.Close()
	db := database
	groupeHandler := Nouveau(db).Handler
	if _, err := db.Exec("CREATE TABLE clients (id INTEGER PRIMARY KEY, name TEXT); INSERT INTO clients VALUES (1, 'Joueur')"); err != nil {
		t.Fatal(err)
	}
	if err := Initialiser(db); err != nil {
		t.Fatal(err)
	}
	session := httptest.NewRecorder()
	SQLite.CreerSession(session, 1)
	cookie := session.Result().Cookies()[0]
	post := func(path, message string, authenticated bool) *httptest.ResponseRecorder {
		form := url.Values{"riot_id": {"Joueur#EUW"}, "message": {message}}
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if authenticated {
			req.AddCookie(cookie)
		}
		res := httptest.NewRecorder()
		groupeHandler(res, req)
		return res
	}
	if res := post("/duoq", "Refusée", false); res.Code != 303 || res.Header().Get("Location") != "/login" {
		t.Fatal("Connexion obligatoire")
	}
	if res := post("/duoq", "  ", true); res.Code != 400 {
		t.Fatal("Annonce vide acceptée")
	}
	if res := post("/duoq", "<script>alert(1)</script>", true); res.Code != 303 {
		t.Fatal(res.Body.String())
	}
	if res := post("/flex", "Équipe de cinq", true); res.Code != 303 {
		t.Fatal(res.Body.String())
	}
	for _, path := range []string{"/duoq", "/flex"} {
		res := httptest.NewRecorder()
		groupeHandler(res, httptest.NewRequest(http.MethodGet, path, nil))
		body := res.Body.String()
		if res.Code != 200 {
			t.Fatal(body)
		}
		if strings.Contains(body, "<script>") {
			t.Fatal("HTML non échappé")
		}
		if path == "/duoq" && (!strings.Contains(body, "&lt;script&gt;") || strings.Contains(body, "Équipe de cinq")) {
			t.Fatal("Annonces duoQ incorrectes")
		}
		if path == "/flex" && (!strings.Contains(body, "Équipe de cinq") || strings.Contains(body, "alert(1)")) {
			t.Fatal("Annonces flex incorrectes")
		}
	}
}
