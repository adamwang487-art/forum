package handler

import (
	"database/sql"
	"forum/serveur"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginUnknownAccount(t *testing.T) {
	t.Chdir("..")
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE clients (id INTEGER, gmail TEXT, hash TEXT)"); err != nil {
		t.Fatal(err)
	}
	previousDB := serveur.DB
	serveur.DB = db
	defer func() { serveur.DB = previousDB }()
	email := "inconnu@example.com"
	form := url.Values{"email": {email}, "password": {"secret-ne-pas-afficher"}}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	Login(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status: %d", res.Code)
	}
	if !strings.HasPrefix(res.Header().Get("Content-Type"), "text/html") {
		t.Fatal("expected HTML response")
	}
	for _, content := range []string{"Mot de passe ou adresse mail incorrect.", `role="alert"`, `/static/login.css`, `href="/register"`, "Réessayer la connexion", email} {
		if !strings.Contains(res.Body.String(), content) {
			t.Errorf("missing %q", content)
		}
	}
	if strings.Contains(res.Body.String(), "secret-ne-pas-afficher") {
		t.Fatal("password included in response")
	}
	res = httptest.NewRecorder()
	Login(res, httptest.NewRequest(http.MethodGet, "/login", nil))
	if res.Code != http.StatusOK || strings.Contains(res.Body.String(), `role="alert"`) {
		t.Fatal("unexpected initial login page")
	}
}
