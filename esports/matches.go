package esports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Team struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Acronym string `json:"acronym"`
	Image   string `json:"image_url"`
}
type Match struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	BeginAt     *time.Time `json:"begin_at"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	League      struct {
		Name string `json:"name"`
	} `json:"league"`
	Opponents []struct {
		Opponent Team `json:"opponent"`
	} `json:"opponents"`
	Results []struct {
		TeamID int `json:"team_id"`
		Score  int `json:"score"`
	} `json:"results"`
}
type Feed struct {
	Past      []Match   `json:"past"`
	Running   []Match   `json:"running"`
	Upcoming  []Match   `json:"upcoming"`
	UpdatedAt time.Time `json:"updated_at"`
	Stale     bool      `json:"stale"`
	Message   string    `json:"message,omitempty"`
}

type Service struct {
	client  *http.Client
	baseURL string
	token   string
	mu      sync.Mutex
	cached  Feed
	retryAt time.Time
}

func New() *Service {
	return &Service{client: &http.Client{Timeout: 8 * time.Second}, baseURL: "https://api.pandascore.co", token: os.Getenv("PANDASCORE_TOKEN")}
}

func (s *Service) fetch(ctx context.Context, kind string) ([]Match, error) {
	sort := "begin_at"
	if kind == "past" {
		sort = "-begin_at"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/lol/matches/"+kind+"?per_page=5&sort="+sort, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PandaScore status %d", resp.StatusCode)
	}
	matches := []Match{}
	err = json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&matches)
	if matches == nil {
		matches = []Match{}
	}
	return matches, err
}

func (s *Service) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Méthode non autorisée."})
		return
	}
	if s.token == "" {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Les résultats sont momentanément indisponibles."})
		return
	}
	// Une seule mise à jour partagée entre les visiteurs, au maximum chaque minute.
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	if !now.Before(s.retryAt) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		type result struct {
			kind    string
			matches []Match
			err     error
		}
		results := make(chan result, 3)
		for _, kind := range []string{"past", "running", "upcoming"} {
			go func(kind string) { matches, err := s.fetch(ctx, kind); results <- result{kind, matches, err} }(kind)
		}
		next := Feed{}
		failed := false
		for range 3 {
			result := <-results
			if result.err != nil {
				failed = true
				continue
			}
			switch result.kind {
			case "past":
				next.Past = result.matches
			case "running":
				next.Running = result.matches
			case "upcoming":
				next.Upcoming = result.matches
			}
		}
		if failed {
			s.retryAt = now.Add(30 * time.Second)
			s.cached.Stale = true
			s.cached.Message = "Actualisation indisponible. Derniers résultats reçus affichés."
		} else {
			next.UpdatedAt = time.Now().UTC()
			s.cached = next
			s.retryAt = now.Add(time.Minute)
		}
	}
	if s.cached.UpdatedAt.IsZero() || time.Since(s.cached.UpdatedAt) > 15*time.Minute {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Les résultats sont momentanément indisponibles. Réessayez dans un instant."})
		return
	}
	_ = json.NewEncoder(w).Encode(s.cached)
}
