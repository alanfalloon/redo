package main

import (
	"testing"
	"database/sql"
)

func tdbconn() *db {
	_conn, err := sql.Open("sqlite3", ":memory:")
	check(err)
	conn := db{_conn, 0}
	conn.xExec("PRAGMA foreign_keys = ON;")
	conn.xExec("PRAGMA temp_store = MEMORY;")
	dbinit(conn)
	return &conn
}

func TestDbinit(t *testing.T) {
	_ = tdbconn()
}

func TestDemand(t *testing.T) {
}
