package store

import "fmt"

type Report struct {
	Counts          map[string]int
	Orphans         map[string]int
	DiscogsCoverage map[string]int
}

func (d *DB) TableCounts(tables []string) (map[string]int, error) {
	out := map[string]int{}
	for _, t := range tables {
		var n int
		if err := d.db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", t)).Scan(&n); err != nil {
			return out, err
		}
		out[t] = n
	}
	return out, nil
}

func (d *DB) OrphanCount(childTable, childCol, parentTable string) (int, error) {
	q := fmt.Sprintf(
		`SELECT count(*) FROM %s c WHERE c.%s IS NOT NULL AND NOT EXISTS (SELECT 1 FROM %s p WHERE p.id = c.%s)`,
		childTable, childCol, parentTable, childCol)
	var n int
	err := d.db.QueryRow(q).Scan(&n)
	return n, err
}
