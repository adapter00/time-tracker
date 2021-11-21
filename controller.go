package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Controller struct {
	db *sqlx.DB
}

func NewController(db *sqlx.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) Start(trackType TrackType) error {
	aw, err := c.ShowLatestTrack(attendance)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		switch trackType {
		case attendance:
			if aw.IsTracking() {
				return fmt.Errorf("latest attendant track is working.%s", aw.StartAt.String())
			}
		case rest:
			if !aw.IsTracking() {
				return fmt.Errorf("latest attendant track is not working.%s", aw.StartAt.String())
			}
		}
	}
	lt, err := c.ShowLatestTrack(rest)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		if lt.IsTracking() {
			return fmt.Errorf("latest track is woring %s, type:%d", lt.StartAt.String(), lt.Type)
		}
	}
	now := time.Now()
	track := Track{
		StartAt:   now,
		Type:      trackType,
		CreatedAt: now,
		UpdatedAt: now,
	}
	log.Printf("insert track")
	_, err = c.db.NamedExec("insert into tracks (start_at,type,created_at, updated_at) values (:start_at,:type,:created_at, :updated_at);", track)
	return err
}

func (c *Controller) Stop(finishAt time.Time) ([]Track, error) {
	tracks := make([]Track, 0)
	aw, err := c.ShowLatestTrack(attendance)
	if err != nil {
		if err != sql.ErrNoRows {
			return tracks, err
		}
	} else {
		log.Printf("show latest attedance:%v", aw)
		if aw.FinishedAt == nil || aw.FinishedAt.IsZero() {
			aw.FinishedAt = &finishAt
			tracks = append(tracks, aw)
		}
	}
	lt, err := c.ShowLatestTrack(rest)
	if err != nil {
		if err != sql.ErrNoRows {
			return tracks, err
		}
	} else {
		log.Printf("show latest attedance:%v", lt)
		if lt.FinishedAt == nil || lt.FinishedAt.IsZero() {
			lt.FinishedAt = &finishAt
			tracks = append(tracks, lt)
		}
	}
	if len(tracks) == 0 {
		return tracks, errors.New("not working tracks")
	}
	err = c.updateFinishAt(finishAt, tracks)
	return tracks, err
}

func (c *Controller) StopRest(finishAt time.Time) error {
	lt, err := c.ShowLatestTrack(rest)
	if err != nil {
		return err
	}

	if !lt.IsTracking() {
		return errors.New("not working rest track")

	}
	lt.FinishedAt = &finishAt
	return c.updateFinishAt(finishAt, []Track{lt})
}

func (c *Controller) updateFinishAt(finishAt time.Time, tracks []Track) error {
	ids := make([]uint, 0)
	duplicate := map[uint]struct{}{}
	for _, t := range tracks {
		if _, ok := duplicate[t.ID]; !ok {
			ids = append(ids, t.ID)
			duplicate[t.ID] = struct{}{}
		}
	}
	q := "update tracks set finish_at = ? WHERE id in (?);"
	sql, params, err := sqlx.In(q, finishAt, ids)
	if err != nil {
		return err
	}
	sql = c.db.Rebind(sql)
	if _, err := c.db.Exec(sql, params...); err != nil {
		return err
	}
	return nil
}
func (c *Controller) Delete(track Track) error {
	return nil
}

func (c *Controller) ShowLatestTrack(t TrackType) (Track, error) {
	var track Track
	err := c.db.Get(&track, "select id,type,start_at,finish_at from tracks where type=$1 order by id desc", t)
	return track, err
}

func (c *Controller) ShowIn(start time.Time, end time.Time) ([]Track, error) {
	q := "select type,start_at,finish_at from tracks where ( start_at >= $1 and start_at <= $2 ) or ( finish_at >= $3 and finish_at <= $4 )"
	tracks := []Track{}
	err := c.db.Select(&tracks, q, start, end, start, end)
	return tracks, err
}

func (c *Controller) ShowWorkTimeIn(start time.Time, end time.Time) (time.Duration, error) {
	tracks, err := c.ShowIn(start, end)
	if err != nil {
		return 0, err
	}
	totalDuration := c.calcWorkTime(start, end, tracks)
	return totalDuration, nil
}

func (c *Controller) calcWorkTime(start, end time.Time, tracks []Track) time.Duration {
	var totalDuration time.Duration
	//TODO: Cross-period calculations
	for _, t := range tracks {
		switch t.Type {
		case attendance:
			if t.IsTracking() {
				continue
			}
			if t.StartAt.After(*t.FinishedAt) {
				log.Printf("after finished at start at:%v", t)
				continue
			}
			finishedAt := *t.FinishedAt
			diff := finishedAt.Sub(t.StartAt)
			totalDuration += diff
		case rest:
			if t.IsTracking() {
				continue
			}
			if t.StartAt.After(*t.FinishedAt) {
				log.Printf("after finished at start at:%v", t)
				continue
			}
			finishedAt := *t.FinishedAt
			diff := finishedAt.Sub(t.StartAt)
			totalDuration -= diff
		}
	}
	log.Printf("tutaolDuration:%v", totalDuration)
	return totalDuration
}
