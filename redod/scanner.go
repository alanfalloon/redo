package main

import (
	"sync"
	"fmt"
	"os"
	"time"
)

func poke_scanner() {
	scanner_once.Do(scanner_start)
	poke(scanner_wake)
}

func poke(c chan<- bool) {
	select {
	case c <- true:
	default:
	}
}
func drain(c <-chan bool) {
	select {
	case <-c:
	default:
	}
}

var scanner_wake chan<- bool
var scanner_once sync.Once

func scanner_start() {
	scanner_c := make(chan bool, 1)
	planner_c := make(chan bool, 1)
	scanner_wake = planner_c
	go scanner_main(scanner_c)
	go planner_main(planner_c, scanner_c)
}
func planner_main(wake <-chan bool, wake_scanner chan<- bool) {
	log := logWrap("planner_main: ", log)
	db := dbconn()
	defer db.Close()
	for _ = range wake {
		log.Printf("Wake")
		for {
			drain(wake)
			if !plan_one(db) {
				break
			}
			log.Printf("Again")
		}
		poke(wake_scanner)
	}
}
func scanner_main(wake <-chan bool) {
	log := logWrap("scanner_main: ", log)
	db := dbconn()
	defer db.Close()
	for _ = range wake {
		log.Printf("Wake")
		for {
			drain(wake)
			if !scan_one(db) {
				break
			}
			log.Printf("Again")
		}
		poke_builder()
	}
}

func plan_one(db *db) bool {
	log := logWrap("plan_one: ", log)
	// Any file that a file in the current generation depends on
	// is bumped to the current generation and set to NEED_SCAN.
	res := db.xExec(`
UPDATE files SET generation=?, step=?
WHERE id IN (
  SELECT c.id
  FROM files AS p
  JOIN deps ON p.id = deps.to_make
  JOIN files AS c ON deps.you_need = c.id
  WHERE p.generation = ?
  AND c.generation != p.generation);`,
		generation, NEEDS_SCAN, generation)
	n, err := res.RowsAffected()
	check(err)
	log.Printf("Pulled %d files in to generation %d", n, generation)
	return n > 0
}
func scan_one(db *db) bool {
	log := logWrap("scan_one: ", log)
	// Find the files in the current generation who have no
	// dependencies, or their dependencies have already been
	// scanned.
	rows := db.xQuery(`
SELECT parent.path, parent.id, parent.stat
FROM files AS parent
WHERE parent.generation = ?
AND parent.step = ?
AND NOT EXISTS (
  SELECT 1
  FROM files AS child JOIN deps ON child.id = deps.you_need
  WHERE deps.to_make = parent.id
  AND (child.generation != parent.generation OR child.step = ?));`,
		generation, NEEDS_SCAN, NEEDS_SCAN)
	var had_data bool
	for rows.Next() {
		had_data = true
		var path string
		var id target
		var stat *time.Time

		err := rows.Scan(&path, &id, &stat)
		check(err)
		log.Printf("Handling %v %v %v", path, id, stat)
		stat_file(db, path, id, stat)
	}
	return had_data
}

func stat_file(db *db, path string, id target, stat *time.Time) {
	st, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		db.xExec(`
UPDATE files SET step=?, stat=NULL
WHERE id=?;`, NEEDS_UPDATE, id)
	case err == nil && stat == nil:
		db.xExec(`
UPDATE files SET step=?, stat=?
WHERE id=?;`, NEEDS_UPDATE, st.ModTime(), id)
	case stat != nil && st != nil && st.ModTime().Equal(*stat):
		db.xExec(`
UPDATE files SET step=?
WHERE id=?;`, NOTHING_TO_DO, id)
	default:
		panic(fmt.Sprintf("Unhandled scan case %s %d %v %v %v", path, id, stat, st, err))
	}
}
