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

func TestCreatePostMissingCategory(t *testing.T) {
	t.Chdir("..")
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("CREATE TABLE categories (id INTEGER, nom TEXT); INSERT INTO categories VALUES (1, 'Discussion')"); err != nil {
		t.Fatal(err)
	}
	previous := serveur.DB
	serveur.DB = db
	defer func() { serveur.DB = previous }()
	session := serveur.CreateSession("1")
	defer serveur.DeleteSession(session)
	form := url.Values{"titre": {`Mon titre <test>`}, "contenu": {`Texte <script>alert(1)</script>`}}
	req := httptest.NewRequest(http.MethodPost, "/create-post", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session_id", Value: session})
	res := httptest.NewRecorder()
	CreationPublicationHandler(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("status %d: %s", res.Code, res.Body.String())
	}
	for _, content := range []string{`role="alert"`, "Choisissez au moins une catégorie", `href="#categories"`, "Discussion", "Mon titre &lt;test&gt;", "Texte &lt;script&gt;"} {
		if !strings.Contains(res.Body.String(), content) {
			t.Errorf("missing %q", content)
		}
	}
	if strings.Contains(res.Body.String(), "<script>") {
		t.Fatal("unescaped content")
	}
	if !strings.HasPrefix(res.Header().Get("Content-Type"), "text/html") {
		t.Fatal("expected HTML")
	}
	res = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/create-post", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: session})
	CreationPublicationHandler(res, req)
	if res.Code != http.StatusOK || strings.Contains(res.Body.String(), `role="alert"`) || !strings.Contains(res.Body.String(), "Discussion") {
		t.Fatal("unexpected initial form")
	}
}
