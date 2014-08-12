package main

import (
	"sync"
)

var tgt_done_pub chan<- tgt_done_cmd

type subscriber chan<- targetResult
type contacts []subscriber
type addressbook []contacts

type tgt_done_cmd interface {
	Exec(db *db, book addressbook) addressbook
}

type demand struct {
	path     string
	ret_to   subscriber
	at_least state
}

type notify_done struct{}

func demand_target(path string, ret_to chan<- targetResult, at_least state) {
	done_pub_once.Do(done_pub_start)
	tgt_done_pub <- demand{path, ret_to, at_least}
}

var done_pub_once sync.Once

func done_pub_start() {
	c := make(chan tgt_done_cmd)
	tgt_done_pub = c
	go done_pub_main(c)
}
func done_pub_main(demands <-chan tgt_done_cmd) {
	var subscribers addressbook
	db := dbconn()
	db.xExec(`
CREATE TEMPORARY TABLE demands (
   file INTEGER NOT NULL,
   idx INTEGER NOT NULL);
`)
	for dmnd := range demands {
		subscribers = dmnd.Exec(db, subscribers)
	}
}

func (dmnd demand) Exec(db *db, subscribers addressbook) addressbook {
	var file, gen int
	var step state
	var x string
	var index *int
	err := db.xQueryRow(`
SELECT id, idx, generation, step
FROM files LEFT JOIN demands ON file = id
WHERE path = ?`, dmnd.path).Scan(&file, &index, &gen, &x)
	check(err)
	step = ERROR
	if index == nil {
		if gen != generation || step < dmnd.at_least {
			db.xExec(`
INSERT OR REPLACE INTO files(id, path, generation, step) VALUES(?, ?, ? ,?);`,
				file, dmnd.path, generation, step)
			poke_scanner()
		}
		index = new(int)
		*index = len(subscribers)
		db.xExec(`
INSERT INTO demands(file, idx) VALUES (?, ?);`, file, *index)
		return append(subscribers, contacts{dmnd.ret_to})
	} else {
		subscribers[*index] = append(subscribers[*index], dmnd.ret_to)
		return subscribers
	}
}

func (_ notify_done) Exec(db *db, ab addressbook) addressbook {
	rows := db.xQuery(`
SELECT file, path, idx, step
FROM files JOIN demands ON file = id
WHERE step >= ?
ORDERBY index DESC`)
	for rows.Next() {
		var file target
		var path string
		var index int
		var step state
		err := rows.Scan(&file, &path, &index, &step)
		check(err)
		contacts := ab[index]
		// FIXME: memory leak: the new slice may still
		// refererence an underlying array referencing old
		// contacts.
		ab = append(ab[:index], ab[index+1:]...)
		res := targetResult{path, file, step}
		for _, subscriber := range contacts {
			subscriber <- res
		}
	}
	return ab
}
