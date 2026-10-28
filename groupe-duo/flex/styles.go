package flex

import (
	_ "embed"
	"net/http"
)

//go:embed groupes.css
var styles string

// Styles sert la feuille de style propre aux annonces duoQ et flexQ.
func Styles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write([]byte(styles))
}
