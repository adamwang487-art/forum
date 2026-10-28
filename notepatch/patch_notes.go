package notepatch

import (
	"io"
	"net/http"
	"regexp"
	"time"
)

const patchNotesURL = "https://www.leagueoflegends.com/fr-fr/news/tags/patch-notes/"

// La liste officielle est classée de la plus récente à la plus ancienne.
var patchLink = regexp.MustCompile(`href=["'](?:https://www\.leagueoflegends\.com)?(/fr-fr/news/game-updates/[^"'<>]*patch[^"'<>]*)["']`)

func Handler(w http.ResponseWriter, r *http.Request) {
	destination := patchNotesURL
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, patchNotesURL, nil)
	if err == nil {
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
				if err == nil {
					if match := patchLink.FindSubmatch(body); len(match) == 2 {
						destination = "https://www.leagueoflegends.com" + string(match[1])
					}
				}
			}
		}
	}
	// En cas d'indisponibilité, la rubrique officielle reste accessible.
	http.Redirect(w, r, destination, http.StatusFound)
}
