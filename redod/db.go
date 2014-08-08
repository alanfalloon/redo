package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

var dbinitonce sync.Once

type db sql.DB

func (db *db) xExec(query string, args ...interface{}) sql.Result {
	r, err := (*sql.DB)(db).Exec(query, args...)
	check(err, query, args)
	return r
}

func (db *db) xQuery(query string, args ...interface{}) *sql.Rows {
	r, err := (*sql.DB)(db).Query(query, args...)
	check(err, query)
	return r
}

func dbconn() *db {
	_conn, err := sql.Open("sqlite3", "file:.redo.db?mode=rwc&vfs=unix-excl")
	check(err)
	conn := (*db)(_conn)
	conn.xExec("PRAGMA journal_mode = WAL;")
	conn.xExec("PRAGMA foreign_keys = ON;")
	dbinitonce.Do(func() { dbinit(conn) })
	return conn
}

const _DBVERSION = 1

var generation int

func dbinit(conn *db) {
	rows := conn.xQuery("PRAGMA user_version;", nil)
	if !rows.Next() {
		panic("empty user_version")
	}
	var version int32
	err := rows.Scan(&version)
	check(err)
	rows.Close()
	switch version {
	case 0:
		dbcreate(conn)
	case _DBVERSION:
		// Current version, increment the generation counter
		conn.xExec(`UPDATE config SET generation=generation+1;`)
	default:
		log.Fatal("Unrecognized DB schema version %d", version)
		panic(version)
	}
	// Read the generation counter
	rows = conn.xQuery(`SELECT generation FROM config`)
	if !rows.Next() {
		panic("empty config")
	}
	err = rows.Scan(&generation)
	check(err)
	rows.Close()
}

func dbcreate(conn *db) {
	initcmd := []string{`
CREATE TABLE files (
id INTEGER PRIMARY KEY AUTOINCREMENT,
path TEXT UNIQUE,
generation INTEGER,
step TEXT,
stat TEXT
);
`,
		`
CREATE TABLE deps (
to_make INTEGER REFERENCES files(id) ON UPDATE CASCADE ON DELETE CASCADE,
you_need INTEGER REFERENCES files(id) ON UPDATE CASCADE ON DELETE CASCADE,
relation TEXT,
generation INTEGER NOT NULL
);
`,
		`
CREATE TABLE config (
generation INTEGER NOT NULL
);
`,
		`
INSERT INTO config VALUES(0);
`,
		`
PRAGMA user_version = 1;
`}
	for _, cmd := range initcmd {
		conn.xExec(cmd)
	}
}

func insert_build_target(path string) int {
	conn := dbconn()
	r := conn.xExec(`
INSERT OR REPLACE
 INTO files(path, generation, step)
 VALUES(?, ?, ?)
`,
		path, generation, "build")
	id, err := r.LastInsertId()
	check(err, path)
	return int(id)
}

func update_target_error(tgtid int, err error) {
	conn := dbconn()
	conn.xExec(`UPDATE files SET step=? WHERE id=?;`,
		"error: "+err.Error(), tgtid)
}

func update_target_done(tgtid int, st os.FileInfo) {
	conn := dbconn()
	var j []byte
	if st != nil {
		stat := map[string]interface{}{
			"size":  st.Size(),
			"mode":  st.Mode(),
			"mtime": st.ModTime()}
		var err error
		j, err = json.Marshal(stat)
		check(err, tgtid, st)
	}
	conn.xExec(`UPDATE files SET step=?, stat=? WHERE id=?;`,
		"updated", string(j), tgtid)
}

func insert_dep(tgt, dep int) {
	if tgt == -1 {
		return
	}
	conn := dbconn()
	conn.xExec(`
INSERT
 INTO deps(to_make, you_need, relation, generation)
 VALUES(?, ?, "ifchange", ?)
`,
		tgt, dep, generation)
}
