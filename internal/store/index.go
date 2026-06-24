// internal/store/index.go
package store

const ftsDDL = `
CREATE VIRTUAL TABLE artist_fts        USING fts5(name,                content='', tokenize='unicode61 remove_diacritics 2');
CREATE VIRTUAL TABLE label_fts         USING fts5(name,                content='', tokenize='unicode61 remove_diacritics 2');
CREATE VIRTUAL TABLE work_fts          USING fts5(title,               content='', tokenize='unicode61 remove_diacritics 2');
CREATE VIRTUAL TABLE release_group_fts USING fts5(title, artist_names, content='', tokenize='unicode61 remove_diacritics 2');
CREATE VIRTUAL TABLE release_fts       USING fts5(title, artist_names, content='', tokenize='unicode61 remove_diacritics 2');
CREATE VIRTUAL TABLE recording_fts     USING fts5(title, artist_names, content='', tokenize='unicode61 remove_diacritics 2');
`

// BuildIndexes executes the embedded generated index DDL.
func (d *DB) BuildIndexes() error {
	_, err := d.db.Exec(indexesSQL)
	return err
}

// BuildFTS creates the six contentless FTS5 virtual tables and populates them.
// artist_names for release_group/release/recording is denormalized from artist_credit_name,
// concatenated in artist_credit_name.position order (in-aggregate ORDER BY, SQLite 3.44+).
// Schema-verified column names (schema_seq31.sql):
//   - artist.name, label.name, work.name, release_group.name, release.name, recording.name
//   - release_group.artist_credit, release.artist_credit, recording.artist_credit
//   - artist_credit_name: artist_credit, position, name
func (d *DB) BuildFTS() error {
	if _, err := d.db.Exec(ftsDDL); err != nil {
		return err
	}
	stmts := []string{
		`INSERT INTO artist_fts(rowid, name) SELECT id, name FROM artist`,
		`INSERT INTO label_fts(rowid, name) SELECT id, name FROM label`,
		`INSERT INTO work_fts(rowid, title) SELECT id, name FROM work`,
		`INSERT INTO release_group_fts(rowid, title, artist_names) SELECT rg.id, rg.name, ` +
			`(SELECT group_concat(acn.name,' ' ORDER BY acn.position) FROM artist_credit_name acn WHERE acn.artist_credit=rg.artist_credit) FROM release_group rg`,
		`INSERT INTO release_fts(rowid, title, artist_names) SELECT r.id, r.name, ` +
			`(SELECT group_concat(acn.name,' ' ORDER BY acn.position) FROM artist_credit_name acn WHERE acn.artist_credit=r.artist_credit) FROM release r`,
		`INSERT INTO recording_fts(rowid, title, artist_names) SELECT rec.id, rec.name, ` +
			`(SELECT group_concat(acn.name,' ' ORDER BY acn.position) FROM artist_credit_name acn WHERE acn.artist_credit=rec.artist_credit) FROM recording rec`,
	}
	for _, s := range stmts {
		if _, err := d.db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
