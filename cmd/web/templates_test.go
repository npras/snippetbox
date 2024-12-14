package main

import (
	"testing"
	"time"

	"github.com/npras/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		desc string
		tm   time.Time
		want string
	}{
		{
			desc: "check if utc",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			desc: "check if empty",
			tm:   time.Time{},
			want: "",
		},
		{
			desc: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := humanDate(test.tm)
			assert.Equal(t, got, test.want)
		})
	}
}
