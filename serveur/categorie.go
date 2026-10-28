package serveur

import "database/sql"

func RecupererCategories(
	db *sql.DB,
) ([]Categorie, error) {

	rows, err := db.Query(`
		SELECT id, nom
		FROM categories
		ORDER BY nom
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Categorie

	for rows.Next() {
		var categorie Categorie

		err := rows.Scan(
			&categorie.ID,
			&categorie.Nom,
		)
		if err != nil {
			return nil, err
		}

		categories = append(
			categories,
			categorie,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func RecupererCategoriesPublication(
	db *sql.DB,
	publicationID int,
) ([]Categorie, error) {

	rows, err := db.Query(`
		SELECT
			categories.id,
			categories.nom

		FROM categories

		INNER JOIN publication_categories
			ON categories.id =
			   publication_categories.categorie_id

		WHERE publication_categories.publication_id = ?

		ORDER BY categories.nom
	`,
		publicationID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Categorie

	for rows.Next() {
		var categorie Categorie

		err := rows.Scan(
			&categorie.ID,
			&categorie.Nom,
		)
		if err != nil {
			return nil, err
		}

		categories = append(
			categories,
			categorie,
		)
	}

	return categories, rows.Err()
}
