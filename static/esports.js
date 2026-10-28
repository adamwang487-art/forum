(() => {
    const panel = document.querySelector('.esports-panel');
    if (!panel) return;
    const list = panel.querySelector('#esports-matches');
    const status = panel.querySelector('#esports-status');
    const refresh = panel.querySelector('#esports-refresh');
    let feed = null;
    let selected = 'past';
    let loading = false;
    const dateFormat = new Intl.DateTimeFormat('fr-FR', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' });
    function node(tag, className, text) {
        const element = document.createElement(tag);
        element.className = className;
        if (text !== undefined) element.textContent = text;
        return element;
    }
    function render() {
        list.replaceChildren();
        if (!feed) return;
        const matches = feed[selected] || [];
        if (!matches.length) {
            const messages = { past: 'Aucun résultat disponible.', running: 'Aucun match en cours.', upcoming: 'Aucune rencontre annoncée.' };
            list.append(node('p', 'esports-empty', messages[selected]));
        }
        for (const match of matches) {
            const card = node('article', 'esports-match');
            card.append(node('p', 'esports-league', match.league?.name || 'League of Legends'));
            const date = match.begin_at || match.scheduled_at;
            const time = node('time', 'esports-date');
            if (date && !Number.isNaN(Date.parse(date))) {
                time.dateTime = date;
                time.textContent = dateFormat.format(new Date(date));
            } else time.textContent = 'Horaire à confirmer';
            card.append(time);
            const label = { finished: 'Terminé', running: 'En cours', not_started: 'À venir', canceled: 'Annulé', postponed: 'Reporté' };
            card.append(node('span', `esports-state ${match.status === 'running' ? 'is-running' : ''}`, label[match.status] || 'À confirmer'));
            const opponents = (match.opponents || []).slice(0, 2);
            while (opponents.length < 2) opponents.push({ opponent: { name: 'Équipe à confirmer' } });
            for (const { opponent: team } of opponents) {
                const row = node('div', 'esports-team');
                const name = node('span', 'esports-team-name', team.acronym || team.name || 'Équipe à confirmer');
                name.title = team.name || '';
                row.append(name);
                const result = (match.results || []).find(result => result.team_id === team.id);
                row.append(node('strong', 'esports-score', ['finished', 'running'].includes(match.status) && result ? String(result.score) : '—'));
                card.append(row);
            }
            list.append(card);
        }
        status.textContent = feed.stale ? feed.message : `Actualisé à ${new Date(feed.updated_at).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })}`;
    }
    async function update() {
        if (loading) return;
        loading = true;
        refresh.disabled = true;
        try {
            const response = await fetch('/api/esports/matches', { signal: AbortSignal.timeout(15000) });
            const data = await response.json();
            if (!response.ok) throw new Error('unavailable');
            feed = data;
            render();
        } catch {
            if (feed && Date.now() - Date.parse(feed.updated_at) > 15 * 60 * 1000) {
                feed = null;
                list.replaceChildren();
            }
            status.textContent = feed ? 'Actualisation indisponible. Les derniers résultats reçus restent affichés.' : 'Résultats indisponibles pour le moment. Réessayez dans un instant.';
        } finally {
            loading = false;
            refresh.disabled = false;
        }
    }
    panel.querySelectorAll('[data-feed]').forEach(button => button.addEventListener('click', () => {
        selected = button.dataset.feed;
        panel.querySelectorAll('[data-feed]').forEach(item => item.setAttribute('aria-pressed', String(item === button)));
        render();
    }));
    refresh.addEventListener('click', update);
    update();
    setInterval(() => { if (!document.hidden) update(); }, 60000);
})();
