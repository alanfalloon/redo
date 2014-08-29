package main

import (
	"sync"
	"database/sql"
)

func poke_builder() {
	builder_once.Do(builder_start)
	poke(builder_wake)
}

var builder_wake chan<- bool
var builder_once sync.Once

func builder_start() {
	c := make(chan bool, 1)
	builder_wake = c
	go builder_main(c)
}
func builder_main(wake <-chan bool) {
	log := logWrap("builder_main: ", log)
	db := dbconn()
	defer db.Close()
	for _ = range wake {
		log.Printf("Wake")
		for {
			drain(wake)
			if !build_one(db) {
				break
			}
			log.Printf("Again")
		}
	}
}

func build_one(db *db) bool {
	var path string
	var id target
	err := db.QueryRow(`
SELECT parent.path, parent.id
FROM files AS parent
WHERE parent.generation = ?
AND parent.step = ?
AND NOT EXISTS (
  SELECT 1
  FROM files AS child JOIN deps ON child.id = deps.you_need
  WHERE deps.to_make = parent.id
  AND (child.generation != parent.generation OR child.step <= ?));`,
		db.generation, NEEDS_UPDATE, NEEDS_UPDATE).Scan(&path, &id)
	if err == sql.ErrNoRows {
		return false
	}
	dofile, cwd, tgt, base := find_dofile(path)
	err = run(dofile, cwd, tgt, base, id)
	check(err)
	pub.poke_publisher()
	return true
}
