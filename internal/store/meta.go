package store

const metaDDL = `CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);`

func (d *DB) WriteMeta(kv map[string]string) error {
	if _, err := d.db.Exec(metaDDL); err != nil {
		return err
	}
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO meta(key,value) VALUES(?,?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`)
	if err != nil {
		tx.Rollback()
		return err
	}
	for k, v := range kv {
		if _, err := stmt.Exec(k, v); err != nil {
			stmt.Close()
			tx.Rollback()
			return err
		}
	}
	stmt.Close()
	return tx.Commit()
}
