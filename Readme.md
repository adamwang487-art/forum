# SQLite – Forum

## Installation

Installer SQLite pour Go :

```bash
go get modernc.org/sqlite
```

La base de données `chat.db` est créée automatiquement.

## Base de données

### Table `clients`

La table `clients` contient 4 colonnes :

* `id`
* `gmail`
* `name`
* `motdepasse`

Description :

* `id` : identifiant automatique
* `gmail` : adresse Gmail unique
* `name` : nom unique
* `motdepasse` : mot de passe du client

## Fonctions

### Ouvrir la base de données

```go
db := Ouvrirdata()
defer db.Close()
```

### Ajouter un client

```go
AjouteClient(db, gmail, name, motdepasse)
```

### Afficher les clients

```go
afficherClients(db)
```

## Commandes SQL principales

```sql
INSERT  -- ajouter
SELECT  -- récupérer
UPDATE  -- modifier
DELETE  -- supprimer
```

---

## `db.Query`

`db.Query` permet de récupérer **plusieurs lignes** dans la base de données.

Exemple :

```go
rows, err := db.Query("SELECT id, gmail, name, motdepasse FROM clients")
```

Explication :

```sql
SELECT id, gmail, name, motdepasse
```

→ On récupère les colonnes `id`, `gmail`, `name` et `motdepasse`.

```sql
FROM clients
```

→ On récupère ces informations dans la table `clients`.

---

## `db.QueryRow`

`db.QueryRow` permet de chercher **une seule ligne** dans la base de données.

Exemple :

```go
err := db.QueryRow(
    `SELECT id FROM clients WHERE name = ?`,
    "adam",
).Scan(&id)
```

Explication :

```sql
SELECT id
```

→ On récupère la colonne `id`.

```sql
FROM clients
```

→ On recherche dans la table `clients`.

```sql
WHERE name = ?
```

→ On recherche une ligne dont la colonne `name` correspond à la valeur donnée.

```go
"adam"
```

→ Le `?` est remplacé par `"adam"`.

```go
Scan(&id)
```

→ Si `adam` est présent dans la table, son `id` est placé dans la variable `id`.

---

## `db.Exec`

`db.Exec` permet d'exécuter une commande qui **modifie la base de données**.

Cela peut être par exemple :

* ajouter une ligne avec `INSERT`
* modifier une ligne avec `UPDATE`
* supprimer une ligne avec `DELETE`
* créer une table avec `CREATE TABLE`

Exemple :

```go
_, err := db.Exec(
    `INSERT INTO clients (gmail, name, motdepasse)
     VALUES (?, ?, ?)`,
    gmail, name, mot,
)
```

Explication :

```sql
INSERT INTO clients (gmail, name, motdepasse)
```

→ On ajoute un nouveau client dans la table `clients`.

La colonne `id` n'est pas indiquée car elle est générée automatiquement.

```sql
VALUES (?, ?, ?)
```

→ Ces valeurs remplacent les `?` dans le même ordre :

```text
? → gmail
? → name
? → mot
```
package serveur

func ListeAvatars() []string {
	return []string{
		"Aatrox",
		"Ahri",
		"Akali",
		"Akshan",
		"Alistar",
		"Ambessa",
		"Amumu",
		"Anivia",
		"Annie",
		"Aphelios",
		"Ashe",
		"Aurelion Sol",
		"Aurora",
		"Azir",
		"Bard",
		"Bel'veth",
		"Blitzcrank",
		"Brand",
		"Braum",
		"Briar",
		"Caitlyn",
		"Camille",
		"Cassiopeia",
		"Cho'gath",
		"Corki",
		"Darius",
		"Diana",
		"DrMundo",
		"Draven",
		"Ekko",
		"Elise",
		"Evelynn",
		"Ezreal",
		"Fiddlesticks",
		"Fiora",
		"Fizz",
		"Galio",
		"Gangplank",
		"Garen",
		"Gnar",
		"Gragas",
		"Graves",
		"Gwen",
		"Hecarim",
		"Heimerdinger",
		"Hwei",
		"Illaoi",
		"Irelia",
		"Ivern",
		"Janna",
		"JarvanIV",
		"Jax",
		"Jayce",
		"Jhin",
		"Jinx",
		"KSante",
		"Kaisa",
		"Kalista",
		"Karma",
		"Karthus",
		"Kassadin",
		"Katarina",
		"Kayle",
		"Kayn",
		"Kennen",
		"Khazix",
		"Kindred",
		"Kled",
		"KogMaw",
		"Leblanc",
		"LeeSin",
		"Leona",
		"Lillia",
		"Lissandra",
		"Locke",
		"Lucian",
		"Lulu",
		"Lux",
		"Malphite",
		"Malzahar",
		"Maokai",
		"MasterYi",
		"Mel",
		"Milio",
		"MissFortune",
		"Mordekaiser",
		"Morgana",
		"Naafiri",
		"Nami",
		"Nasus",
		"Nautilus",
		"Neeko",
		"Nidalee",
		"Nilah",
		"Nocturne",
		"Nunu",
		"Olaf",
		"Orianna",
		"Ornn",
		"Pantheon",
		"Poppy",
		"Pyke",
		"Qiyana",
		"Quinn",
		"Rakan",
		"Rammus",
		"RekSai",
		"Rell",
		"RenataGlasc",
		"Renekton",
		"Rengar",
		"Riven",
		"Rumble",
		"Ryze",
		"Samira",
		"Sejuani",
		"Senna",
		"Seraphine",
		"Sett",
		"Shaco",
		"Shen",
		"Shyvana",
		"Singed",
		"Sion",
		"Sivir",
		"Skarner",
		"Smolder",
		"Sona",
		"Soraka",
		"Swain",
		"Sylas",
		"Syndra",
		"TahmKench",
		"Taliyah",
		"Talon",
		"Taric",
		"Teemo",
		"Thresh",
		"Tristana",
		"Trundle",
		"Tryndamere",
		"TwistedFate",
		"Twitch",
		"Udyr",
		"Urgot",
		"Varus",
		"Vayne",
		"Veigar",
		"Velkoz",
		"Vex",
		"Vi",
		"Viego",
		"Viktor",
		"Vladimir",
		"Volibear",
		"Warwick",
		"Wukong",
		"Xayah",
		"Xerath",
		"XinZhao",
		"Yasuo",
		"Yone",
		"Yorick",
		"Yunara",
		"Yuumi",
		"Zaahen",
		"Zac",
		"Zed",
		"Zeri",
		"Ziggs",
		"Zilean",
		"Zoe",
		"Zyra",
	}
<<<<<<< HEAD
}
=======
}
>>>>>>> origin/data
