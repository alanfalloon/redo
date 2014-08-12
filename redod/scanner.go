package main

import (
	"sync"
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
	walker_c := make(chan bool, 1)
	scanner_wake = walker_c
	go scanner_main(scanner_c)
	go walker_main(walker_c, scanner_c)
}
func walker_main(wake <-chan bool, wake_scanner chan<- bool) {
	db := dbconn()
	for _ = range wake {
		for {
			drain(wake)
			if !walk_one(db) {
				break
			}
		}
		poke(wake_scanner)
	}
}
func scanner_main(wake <-chan bool) {
	db := dbconn()
	for _ = range wake {
		for {
			drain(wake)
			if !scan_one(db) {
				break
			}
		}
	}
}

func walk_one(db *db) bool {
	res := db.xExec(`
UPDATE files SET generation=?, step=?
WHERE id IN (
  SELECT c.id
  FROM files AS p JOIN deps ON p.id = deps.to_make JOIN files AS c ON deps.you_need = c.id
  WHERE p.generation = ?
  AND c.generation != p.generation);`,
		generation, NEEDS_SCAN, generation)
	n, err := res.RowsAffected()
	check(err)
	return n > 0
}
func scan_one(db *db) bool {
	return false
}
