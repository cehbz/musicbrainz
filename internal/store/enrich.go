package store

import (
	"fmt"
	"regexp"
	"strconv"
)

var reDiscogs = regexp.MustCompile(`discogs\.com/(?:[a-z]{2}/)?(artist|label|master|release)/(\d+)`)

func ParseDiscogsID(url, entityType string) (int, bool) {
	m := reDiscogs.FindStringSubmatch(url)
	if m == nil || m[1] != entityType {
		return 0, false
	}
	id, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, false
	}
	return id, true
}

// enrichSpec wires one MB entity to its Discogs id space.
type enrichSpec struct {
	entity     string // mb table
	linkTable  string // l_<entity>_url
	column     string // discogs_*_id
	discogsTyp string // artist|label|master|release
}

var enrichSpecs = []enrichSpec{
	{"artist", "l_artist_url", "discogs_artist_id", "artist"},
	{"label", "l_label_url", "discogs_label_id", "label"},
	{"release_group", "l_release_group_url", "discogs_master_id", "master"},
	{"release", "l_release_url", "discogs_release_id", "release"},
}

func (d *DB) EnrichDiscogs() (map[string]int, error) {
	cov := map[string]int{}
	for _, s := range enrichSpecs {
		q := fmt.Sprintf(
			`SELECT lk.entity0, u.url FROM %s lk JOIN url u ON u.id = lk.entity1 WHERE u.url LIKE '%%discogs.com/%%'`,
			s.linkTable)
		rows, err := d.db.Query(q)
		if err != nil {
			return cov, fmt.Errorf("enrich query %s: %w", s.entity, err)
		}
		type pair struct {
			id   int
			disc int
		}
		var pairs []pair
		for rows.Next() {
			var entityID int
			var url string
			if err := rows.Scan(&entityID, &url); err != nil {
				rows.Close()
				return cov, err
			}
			if disc, ok := ParseDiscogsID(url, s.discogsTyp); ok {
				pairs = append(pairs, pair{entityID, disc})
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return cov, fmt.Errorf("enrich scan %s: %w", s.entity, err)
		}

		// If an entity has multiple Discogs links of the same type, the last one in the
		// result set wins (the UPDATE overwrites). MusicBrainz entities rarely have more
		// than one, so a deterministic-but-arbitrary choice is acceptable.
		tx, err := d.db.Begin()
		if err != nil {
			return cov, err
		}
		stmt, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET %s=? WHERE id=?", s.entity, s.column))
		if err != nil {
			tx.Rollback()
			return cov, err
		}
		for _, p := range pairs {
			if _, err := stmt.Exec(p.disc, p.id); err != nil {
				stmt.Close()
				tx.Rollback()
				return cov, err
			}
		}
		stmt.Close()
		if err := tx.Commit(); err != nil {
			return cov, err
		}
		cov[s.entity] = len(pairs)
	}
	return cov, nil
}
