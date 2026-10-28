package SQLite

import (
	"database/sql"
	"testing"
)

func TestMigrationVotesEtAvatar(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE clients (id INTEGER PRIMARY KEY);
 INSERT INTO clients VALUES (1);
 CREATE TABLE likes (client_id INTEGER, publication_id INTEGER, UNIQUE(client_id, publication_id));
 INSERT INTO likes VALUES (1, 2);
 CREATE TABLE votes (client_id INTEGER, publication_id INTEGER, valeur INTEGER);
 INSERT INTO votes VALUES (1, 3, -1);`)
	if err != nil {
		t.Fatal(err)
	}
	if err := migrerDonnees(db); err != nil {
		t.Fatal(err)
	}
	var avatar string
	if err := db.QueryRow("SELECT avatar FROM clients WHERE id = 1").Scan(&avatar); err != nil || avatar != "Ahri" {
		t.Fatalf("avatar=%q, err=%v", avatar, err)
	}
	var count, value int
	if err := db.QueryRow("SELECT count(*), valeur FROM likes").Scan(&count, &value); err != nil || count != 1 || value != -1 {
		t.Fatalf("votes: count=%d value=%d err=%v", count, value, err)
	}
	if _, err := db.Exec("DELETE FROM likes"); err != nil {
		t.Fatal(err)
	}
	if err := migrerDonnees(db); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRow("SELECT count(*) FROM likes").Scan(&count); err != nil || count != 0 {
		t.Fatalf("Vote annulé réapparu: count=%d err=%v", count, err)
	}
}
