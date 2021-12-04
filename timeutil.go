package main

import "time"

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func BeginAndLateDateInMonth(t time.Time) (time.Time, time.Time) {
	begin := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location()).Add(-1 * time.Second)
	return begin, end
}

func toJST(t time.Time) time.Time {
	return t.In(jst)
}
