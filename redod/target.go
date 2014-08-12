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
