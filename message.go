package main

import (
	"fmt"
	"time"
)

const (
	workTimeLayout = "2006-01-02 15:04:05"
)

func stopMultipleMessage(tracks []Track) string {
	message := "multiple stop"
	for _, t := range tracks {
		m := fmt.Sprintf("day:%s %s", toJST(t.CreatedAt), showWorkTime(t))
		message = fmt.Sprintf("%s\n%s", message, m)
	}
	return message
}

func showWorkTimes(startAt time.Time, tracks []Track) string {
	message := fmt.Sprintf("work time in :%s", toJST(startAt).Format(workTimeLayout))
	for _, t := range tracks {
		message = fmt.Sprintf("%s\n%s", message, showWorkTime(t))
	}
	return message
}

func showWorkTime(track Track) string {
	if track.IsTracking() {
		return fmt.Sprintf(" - *%s*  \n 🕰: none`", toJST(track.StartAt).Format(workTimeLayout))
	}
	return fmt.Sprintf(" - *%s~%s*  \n 🕰: `%s`", toJST(track.StartAt).Format(workTimeLayout), toJST(*track.FinishedAt).Format(workTimeLayout), track.FinishedAt.Sub(track.StartAt))
}
