package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	controller := NewController(testDB)
	t.Run("start", func(t *testing.T) {
		testDB.Exec("truncate table tracks;")
		require.NoError(t, controller.Start(attendance))
		require.NoError(t, controller.Start(rest))
	})
}
func TestCal(t *testing.T) {
	controller := NewController(testDB)
	type testCase struct {
		name string
		in   []Track
		want time.Duration
	}
	now := time.Now()
	testCases := []testCase{
		{
			name: "not track times",
			in:   []Track{},
			want: time.Duration(0),
		},
		{
			name: "1hour",
			in: []Track{
				Track{
					StartAt:    now.Add(-1 * time.Hour),
					FinishedAt: &now,
					Type:       attendance,
				},
			},
			want: time.Duration(1 * time.Hour),
		},
		{
			name: "1hour and 15 rest",
			in: []Track{
				Track{
					StartAt:    now.Add(-1 * time.Hour),
					FinishedAt: &now,
					Type:       attendance,
				},
				Track{
					StartAt:    now.Add(-30 * time.Minute),
					FinishedAt: toP(now.Add(-15 * time.Minute)),
					Type:       rest,
				},
			},
			want: time.Duration(45 * time.Minute),
		},
		{
			name: "1hour and 15 rest",
			in: []Track{
				Track{
					StartAt:    now.Add(-1 * time.Hour),
					FinishedAt: &now,
					Type:       attendance,
				},
				Track{
					StartAt:    now.Add(-30 * time.Minute),
					FinishedAt: toP(now.Add(-15 * time.Minute)),
					Type:       rest,
				},
				Track{
					StartAt:    now.Add(-1*time.Hour + 24*time.Hour),
					FinishedAt: toP(now.Add(24 * time.Hour)),
					Type:       attendance,
				},
			},
			want: time.Duration(1*time.Hour + 45*time.Minute),
		},
		{
			name: "skip not finished",
			in: []Track{
				Track{
					StartAt: now.Add(-1*time.Hour + -24*time.Hour),
					Type:    attendance,
				},
				Track{
					StartAt:    now.Add(-1 * time.Hour),
					FinishedAt: &now,
					Type:       attendance,
				},
				Track{
					StartAt: now.Add(-30 * time.Minute),
					Type:    rest,
				},
			},
			want: time.Duration(1 * time.Hour),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := controller.calcWorkTime(now, now, tc.in)
			assert.Equal(t, tc.want, result)
		})
	}
}

func toP(t time.Time) *time.Time {
	return &t
}
