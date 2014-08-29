package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"time"
)

func TestAlreadyDone(t *testing.T) {
	db := tdbconn()
	require.NotNil(t, db)
	defer db.Close()
	db.generation = 1
	pub := publisher{db: db}
	b := make(chan targetResult)
	f := db.target("a")
	pub.demand_target("a", b)
	select {
	case r := <-b:
		t.Errorf("Unexpectedly ready with %v", r)
	case <-time.After(time.Second / 2):
	}
	f.setStep(db, NOTHING_TO_DO)
	pub.poke_publisher()
	select {
	case r := <-b:
		assert.Equal(t, "a", r.name)
		assert.Equal(t, NOTHING_TO_DO, r.outcome)
	case <-time.After(5 * time.Second):
		t.Error("Timeout after 5 seconds")
	}
	
}

type testing_db testing.T
func (t *testing_db) dbcreate() *db {
	db := tdbconn()
	require.NotNil((*testing.T)(t), db)
	db.generation = 1
	return db
}
