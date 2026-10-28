package SQLite

import "database/sql"

// Adapte les bases existantes aux profils et au stockage commun des votes.
func migrerDonnees(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, column := range []struct{ table, name, definition string }{
		{"clients", "avatar", "TEXT NOT NULL DEFAULT 'Ahri'"},
		{"likes", "valeur", "INTEGER NOT NULL DEFAULT 1"},
	} {
		var count int
		if err := tx.QueryRow("SELECT count(*) FROM pragma_table_info(?) WHERE name = ?", column.table, column.name).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			if _, err := tx.Exec("ALTER TABLE " + column.table + " ADD COLUMN " + column.name + " " + column.definition); err != nil {
				return err
			}
		}
	}
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS migrations (nom TEXT PRIMARY KEY)"); err != nil {
		return err
	}
	var done, exists int
	if err := tx.QueryRow("SELECT count(*) FROM migrations WHERE nom = 'votes_vers_likes'").Scan(&done); err != nil {
		return err
	}
	if done == 0 {
		if err := tx.QueryRow("SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'votes'").Scan(&exists); err != nil {
			return err
		}
		if exists != 0 {
			// La table votes représente l'état courant, y compris les likes annulés.
			if _, err := tx.Exec("DELETE FROM likes; INSERT INTO likes (client_id, publication_id, valeur) SELECT client_id, publication_id, valeur FROM votes"); err != nil {
				return err
			}
		}
		if _, err := tx.Exec("INSERT INTO migrations (nom) VALUES ('votes_vers_likes')"); err != nil {
			return err
		}
	}
	return tx.Commit()
}
