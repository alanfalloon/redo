package main

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"sync"
)

var dbinitonce sync.Once

func dbconn() *sql.DB {
	conn, err := sql.Open("sqlite3", "file:.redo.db?mode=rwc&vfs=unix-excl")
	check(err)
	dbinitonce.Do(func (){ dbinit(conn) })
	return conn
}

const _DBVERSION = 1

func dbinit(conn *sql.DB) {
	rows, err := conn.Query("PRAGMA user_version;", nil)
	check(err)
	if !rows.Next() {panic("empty user_version")}
	var version int32
	err = rows.Scan(&version);
	check(err)
	rows.Close()
	switch version {
	case 0:
		dbcreate(conn)
	case _DBVERSION:
		// Current version, do nothing
	default:
		log.Fatal("Unrecognized DB schema version %d", version)
		panic(version)
	}
}

func dbcreate(conn *sql.DB) {
	_, err := conn.Exec(`
CREATE TABLE files (
id
);
`)
	check(err)
	_, err = conn.Exec(`PRAGMA user_version = ?;`, _DBVERSION)
}
