package main

import "time"

type TrackType uint

const (
	unknown TrackType = iota
	attendance
	rest
)

type Track struct {
	ID         uint       `db:"id"`
	StartAt    time.Time  `db:"start_at"`
	FinishedAt *time.Time `db:"finish_at"`
	Type       TrackType  `db:"type"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

//IsTracking  Whether the measurement is in progress or not
func (t Track) IsTracking() bool {
	return (t.StartAt.IsZero() || t.FinishedAt == nil || t.FinishedAt.IsZero())
}
