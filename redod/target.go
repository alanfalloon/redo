package main

type target int

type targetResult struct {
	name    string
	target  target
	outcome state
}

type state int

const (
	NEEDS_SCAN state = iota
	NEEDS_UPDATE
	// These are outcomes, previous states are transient
	NOTHING_TO_DO
	UPDATED
	ERROR
	MISSING
)

func (t target) demand(db *db) {
	db.xExec(`
UPDATE files
SET generation=?, step=?
WHERE id = ?
AND (generation < ? OR step < ?);`,
		db.generation, NEEDS_SCAN, t, db.generation, NEEDS_SCAN)
}

func (t target) setStep(db *db, step state) {
	db.xExec(`
UPDATE files SET generation = ?, step = ? WHERE id = ?;`, db.generation, NOTHING_TO_DO, t)
}
