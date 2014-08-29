package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

var dbinitonce sync.Once

type db struct {
	*sql.DB
	generation int
}

func (db *db) xExec(query string, args ...interface{}) sql.Result {
	r, err := db.DB.Exec(query, args...)
	check(err, query, args)
	return r
}

func (db *db) xQuery(query string, args ...interface{}) *sql.Rows {
	r, err := db.DB.Query(query, args...)
	check(err, query)
	return r
}

func dbconn() *db {
	_conn, err := sql.Open("sqlite3", "file:.redo.db?mode=rwc&vfs=unix-excl")
	check(err)
	conn := db{_conn, 0}
	conn.xExec("PRAGMA journal_mode = WAL;")
	conn.xExec("PRAGMA foreign_keys = ON;")
	conn.xExec("PRAGMA temp_store = MEMORY;")
	dbinitonce.Do(func() { dbinit(conn) })
	return &conn
}

const _DBVERSION = 1

func dbinit(conn db) {
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
	err = conn.QueryRow(`SELECT generation FROM config`).Scan(&conn.generation)
	check(err)
}

func dbcreate(conn db) {
	initcmd := []string{`
CREATE TABLE files (
id INTEGER PRIMARY KEY AUTOINCREMENT,
path TEXT UNIQUE,
generation INTEGER,
step INTEGER,
stat TEXT);`,
		`
CREATE TABLE deps (
to_make INTEGER REFERENCES files(id) ON UPDATE CASCADE ON DELETE CASCADE,
you_need INTEGER REFERENCES files(id) ON UPDATE CASCADE ON DELETE CASCADE,
relation TEXT,
generation INTEGER NOT NULL);`,
		`
CREATE TABLE config (
generation INTEGER NOT NULL);`,
		`
INSERT INTO config VALUES(0);`,
		`
PRAGMA user_version = 1;`}
	for _, cmd := range initcmd {
		conn.xExec(cmd)
	}
}

func (conn db) target(path string) target {
	res := conn.xExec(`
INSERT OR REPLACE INTO files(path, generation, step) VALUES(?, ? ,?);`,
		path, conn.generation, NEEDS_SCAN)
	id, err := res.LastInsertId()
	check(err)
	return target(id)
}

func update_target_error(tgtid target) {
	conn := dbconn()
	conn.xExec(`UPDATE files SET step=? WHERE id=?;`,
		ERROR, tgtid)
}

func update_target_done(tgtid target, st os.FileInfo) {
	conn := dbconn()
	if st != nil {
		conn.xExec(`UPDATE files SET step=?, stat=? WHERE id=?;`,
			UPDATED, st.ModTime(), tgtid)
	} else {
		conn.xExec(`UPDATE files SET step=?, stat=NULL WHERE id=?;`,
			UPDATED, tgtid)
	}
}

func insert_dep(tgt, dep target) {
	if tgt == -1 {
		return
	}
	conn := dbconn()
	conn.xExec(`
INSERT
 INTO deps(to_make, you_need, relation, generation)
 VALUES(?, ?, "ifchange", ?)
`,
		tgt, dep, conn.generation)
}
