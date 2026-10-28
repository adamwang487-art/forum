package notepatch

import "testing"

func TestPatchLink(t *testing.T) {
	for _, prefix := range []string{"", "https://www.leagueoflegends.com"} {
		body := `<a href="` + prefix + `/fr-fr/news/game-updates/league-of-legends-patch-26-18-notes">Dernier patch</a><a href="/fr-fr/news/game-updates/patch-26-17-notes/">Ancien patch</a>`
		match := patchLink.FindStringSubmatch(body)
		if len(match) != 2 || match[1] != "/fr-fr/news/game-updates/league-of-legends-patch-26-18-notes" {
			t.Fatal(match)
		}
	}
}
