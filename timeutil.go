package main

import "time"

func BeginAndLateDateInMonth(t time.Time) (time.Time, time.Time) {
	begin := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	end := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.Local).Add(-1 * time.Second)
	return begin, end
}
