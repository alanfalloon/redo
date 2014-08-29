package main

import (
	"sync"
)

type subscriber chan<- targetResult
type contacts []subscriber
type addressbook []contacts

type tgt_done_cmd interface {
	Exec(db *db, book addressbook) addressbook
}

type publisher struct {
	*db
	tgt_done_pub chan<- tgt_done_cmd
	done_pub_once sync.Once
}

var pub *publisher

type demand struct {
	path     string
	ret_to   subscriber
}

type notify_done struct{}

func (pub *publisher) demand_target(path string, ret_to chan<- targetResult) {
	pub.done_pub_once.Do(pub.done_pub_start)
	pub.tgt_done_pub <- demand{path, ret_to}
}
func (pub *publisher) poke_publisher() {
	pub.done_pub_once.Do(pub.done_pub_start)
	pub.tgt_done_pub <- notify_done{}
}

func (pub *publisher) done_pub_start() {
	c := make(chan tgt_done_cmd)
	pub.tgt_done_pub = c
	go pub.done_pub_main(c)
}
func (pub publisher) done_pub_main(demands <-chan tgt_done_cmd) {
	var subscribers addressbook
	pub.db.xExec(`
CREATE TEMPORARY TABLE demands (
   file INTEGER NOT NULL,
   idx INTEGER NOT NULL);
`)
	for dmnd := range demands {
		subscribers = dmnd.Exec(pub.db, subscribers)
	}
}

func (dmnd demand) Exec(db *db, subscribers addressbook) addressbook {
	var file target
	var gen int
	var step state
	var index *int
	err := db.QueryRow(`
SELECT id, idx, generation, step
FROM files LEFT JOIN demands ON file = id
WHERE path = ?`, dmnd.path).Scan(&file, &index, &gen, &step)
	check(err)
	if index == nil {
		file.demand(db)
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
ORDER BY idx DESC`, NOTHING_TO_DO)
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
