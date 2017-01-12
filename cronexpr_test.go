/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_test.go
 * Version: 1.0
 * License: pick the one which suits you best:
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr_test

/******************************************************************************/

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorhill/cronexpr"
)

/******************************************************************************/

type crontimes struct {
	from string
	next string
}

type crontest struct {
	expr   string
	layout string
	times  []crontimes
}

var crontests = []crontest{
	// Seconds
	{
		"* * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:01"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// every 5 Second
	{
		"*/5 * * * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:00:05"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:59", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:59", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:59", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:59", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:59", "2013-01-01 00:00:00"},
		},
	},

	// Minutes
	{
		"* * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:01:00"},
			{"2013-01-01 00:00:59", "2013-01-01 00:01:00"},
			{"2013-01-01 00:59:00", "2013-01-01 01:00:00"},
			{"2013-01-01 23:59:00", "2013-01-02 00:00:00"},
			{"2013-02-28 23:59:00", "2013-03-01 00:00:00"},
			{"2016-02-28 23:59:00", "2016-02-29 00:00:00"},
			{"2012-12-31 23:59:00", "2013-01-01 00:00:00"},
		},
	},

	// Minutes with interval
	{
		"17-43/5 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:17:00"},
			{"2013-01-01 00:16:59", "2013-01-01 00:17:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:32:00"},
			{"2013-01-01 00:50:00", "2013-01-01 01:17:00"},
			{"2013-01-01 23:50:00", "2013-01-02 00:17:00"},
			{"2013-02-28 23:50:00", "2013-03-01 00:17:00"},
			{"2016-02-28 23:50:00", "2016-02-29 00:17:00"},
			{"2012-12-31 23:50:00", "2013-01-01 00:17:00"},
		},
	},

	// Minutes interval, list
	{
		"15-30/4,55 * * * *",
		"2006-01-02 15:04:05",
		[]crontimes{
			{"2013-01-01 00:00:00", "2013-01-01 00:15:00"},
			{"2013-01-01 00:16:00", "2013-01-01 00:19:00"},
			{"2013-01-01 00:30:00", "2013-01-01 00:55:00"},
			{"2013-01-01 00:55:00", "2013-01-01 01:15:00"},
			{"2013-01-01 23:55:00", "2013-01-02 00:15:00"},
			{"2013-02-28 23:55:00", "2013-03-01 00:15:00"},
			{"2016-02-28 23:55:00", "2016-02-29 00:15:00"},
			{"2012-12-31 23:54:00", "2012-12-31 23:55:00"},
			{"2012-12-31 23:55:00", "2013-01-01 00:15:00"},
		},
	},

	// Days of week
	{
		"0 0 * * MON",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Mon 2013-01-07 00:00"},
			{"2013-01-28 00:00:00", "Mon 2013-02-04 00:00"},
			{"2013-12-30 00:30:00", "Mon 2014-01-06 00:00"},
		},
	},
	{
		"0 0 * * friday",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Fri 2013-01-04 00:00"},
			{"2013-01-28 00:00:00", "Fri 2013-02-01 00:00"},
			{"2013-12-30 00:30:00", "Fri 2014-01-03 00:00"},
		},
	},
	{
		"0 0 * * 6,7",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-01-01 00:00:00", "Sat 2013-01-05 00:00"},
			{"2013-01-28 00:00:00", "Sat 2013-02-02 00:00"},
			{"2013-12-30 00:30:00", "Sat 2014-01-04 00:00"},
		},
	},

	// Specific days of week
	{
		"0 0 * * 6#5",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Sat 2013-11-30 00:00"},
		},
	},

	// Work day of month
	{
		"0 0 14W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-31 00:00:00", "Mon 2013-04-15 00:00"},
			{"2013-08-31 00:00:00", "Fri 2013-09-13 00:00"},
		},
	},

	// Work day of month -- end of month
	{
		"0 0 30W * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-03-02 00:00:00", "Fri 2013-03-29 00:00"},
			{"2013-06-02 00:00:00", "Fri 2013-06-28 00:00"},
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
		},
	},

	// Last day of month
	{
		"0 0 L * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2014-01-01 00:00:00", "Fri 2014-01-31 00:00"},
			{"2014-02-01 00:00:00", "Fri 2014-02-28 00:00"},
			{"2016-02-15 00:00:00", "Mon 2016-02-29 00:00"},
		},
	},

	// Last work day of month
	{
		"0 0 LW * *",
		"Mon 2006-01-02 15:04",
		[]crontimes{
			{"2013-09-02 00:00:00", "Mon 2013-09-30 00:00"},
			{"2013-11-02 00:00:00", "Fri 2013-11-29 00:00"},
			{"2014-08-15 00:00:00", "Fri 2014-08-29 00:00"},
		},
	},

	// TODO: more tests
}

func TestExpressions(t *testing.T) {
	for _, test := range crontests {
		for _, times := range test.times {
			from, _ := time.Parse("2006-01-02 15:04:05", times.from)
			expr, err := cronexpr.Parse(test.expr)
			if err != nil {
				t.Errorf(`cronexpr.Parse("%s") returned "%s"`, test.expr, err.Error())
			}
			next := expr.Next(from)
			nextstr := next.Format(test.layout)
			if nextstr != times.next {
				t.Errorf(`("%s").Next("%s") = "%s", got "%s"`, test.expr, times.from, times.next, nextstr)
			}
		}
	}
}

/******************************************************************************/

func TestZero(t *testing.T) {
	from, _ := time.Parse("2006-01-02", "2013-08-31")
	next := cronexpr.MustParse("* * * * * 1980").Next(from)
	if next.IsZero() == false {
		t.Error(`("* * * * * 1980").Next("2013-08-31").IsZero() returned 'false', expected 'true'`)
	}

	next = cronexpr.MustParse("* * * * * 2050").Next(from)
	if next.IsZero() == true {
		t.Error(`("* * * * * 2050").Next("2013-08-31").IsZero() returned 'true', expected 'false'`)
	}

	next = cronexpr.MustParse("* * * * * 2099").Next(time.Time{})
	if next.IsZero() == false {
		t.Error(`("* * * * * 2014").Next(time.Time{}).IsZero() returned 'true', expected 'false'`)
	}
}

/******************************************************************************/

func TestNextN(t *testing.T) {
	expected := []string{
		"Sat, 30 Nov 2013 00:00:00",
		"Sat, 29 Mar 2014 00:00:00",
		"Sat, 31 May 2014 00:00:00",
		"Sat, 30 Aug 2014 00:00:00",
		"Sat, 29 Nov 2014 00:00:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:30")
	result := cronexpr.MustParse("0 0 * * 6#5").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:15")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("0 0 * * 6#5").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

func TestNextN_every5min(t *testing.T) {
	expected := []string{
		"Mon, 2 Sep 2013 08:45:00",
		"Mon, 2 Sep 2013 08:50:00",
		"Mon, 2 Sep 2013 08:55:00",
		"Mon, 2 Sep 2013 09:00:00",
		"Mon, 2 Sep 2013 09:05:00",
	}
	from, _ := time.Parse("2006-01-02 15:04:05", "2013-09-02 08:44:32")
	result := cronexpr.MustParse("*/5 * * * *").NextN(from, uint(len(expected)))
	if len(result) != len(expected) {
		t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
		t.Errorf(`  Expected %d returned time values but got %d instead`, len(expected), len(result))
	}
	for i, next := range result {
		nextStr := next.Format("Mon, 2 Jan 2006 15:04:05")
		if nextStr != expected[i] {
			t.Errorf(`MustParse("*/5 * * * *").NextN("2013-09-02 08:44:30", 5):\n"`)
			t.Errorf(`  result[%d]: expected "%s" but got "%s"`, i, expected[i], nextStr)
		}
	}
}

func TestDST(t *testing.T) {
	var locs [3]*time.Location

	// 1 hour DST, negative UTC offset
	// time.Date(2014, 3, 9, 2, 0, 0, 0, locs[0]) Leap PST -> PDT
	// time.Date(2014, 11, 2, 1, 59, 59, 0, locs[0]) Fall PDT -> PST
	locs[0], _ = time.LoadLocation("America/Los_Angeles")

	// biggest tz leap ever (3 hours), occurred from YAKT to MAGST
	// at time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
	locs[1], _ = time.LoadLocation("Asia/Ust-Nera")

	// 30 mins DST, positive UTC offset
	// time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]) Leap LHST -> LHDT
	// time.Date(2015, 4, 5, 1, 30, 0, 0, locs[2]) Fall LHDT -> LHST
	locs[2], _ = time.LoadLocation("Australia/LHI")

	cases := []struct {
		name     string
		expr     string
		opts     cronexpr.Options
		from     time.Time
		expected []time.Time
	}{
		{
			fmt.Sprintf("%s daily leap skip", locs[0]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 10, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 11, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 12, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 13, 2, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily leap unskip", locs[0]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 9, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 10, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 11, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 12, 2, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap skip", locs[0]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2016, 3, 12, 14, 6, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2016, 3, 13, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 14, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 15, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 16, 14, 5, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap unskip", locs[0]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireLate},
			time.Date(2016, 3, 12, 14, 6, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2016, 3, 13, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 14, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 15, 14, 5, 0, 0, locs[0]),
				time.Date(2016, 3, 16, 14, 5, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap skip", locs[0]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 9, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 4, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 5, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 6, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap unskip", locs[0]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 9, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 4, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 5, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 6, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily quarter-hourly leap skip", locs[0]),
			"0 */15 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 10, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 10, 2, 15, 0, 0, locs[0]),
				time.Date(2014, 3, 10, 2, 30, 0, 0, locs[0]),
				time.Date(2014, 3, 10, 2, 45, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily quarter-hourly leap unskip", locs[0]),
			"0 */15 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 3, 9, 1, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 3, 9, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 3, 15, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 3, 30, 0, 0, locs[0]),
				time.Date(2014, 3, 9, 3, 45, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire early", locs[0]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 11, 1, 2, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]).Add(1 * time.Hour),
				time.Date(2014, 11, 3, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 4, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 5, 2, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire late", locs[0]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2014, 11, 1, 2, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 3, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 4, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 5, 2, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire both", locs[0]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate},
			time.Date(2014, 11, 1, 2, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]).Add(1 * time.Hour),
				time.Date(2014, 11, 2, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 3, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 4, 2, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire early", locs[0]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 11, 2, 0, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 4, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire late", locs[0]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2014, 11, 2, 0, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]).Add(1 * time.Hour),
				time.Date(2014, 11, 2, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 3, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 4, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire twice", locs[0]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate},
			time.Date(2014, 11, 2, 0, 0, 0, 0, locs[0]),
			[]time.Time{
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 1, 0, 0, 0, locs[0]).Add(1 * time.Hour),
				time.Date(2014, 11, 2, 2, 0, 0, 0, locs[0]),
				time.Date(2014, 11, 2, 3, 0, 0, 0, locs[0]),
			},
		},
		{
			fmt.Sprintf("%s daily leap skip", locs[1]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 2, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 2, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 3, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 4, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 5, 2, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s daily leap unskip", locs[1]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 2, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 3, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 4, 2, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap skip", locs[1]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 15, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 3, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 4, 14, 5, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap unskip", locs[1]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 15, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 3, 14, 5, 0, 0, locs[1]),
				time.Date(1981, 4, 4, 14, 5, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap skip", locs[1]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 23, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 4, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 5, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 6, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap unskip", locs[1]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 23, 0, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 4, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 5, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 6, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s quarter-hourly leap skip", locs[1]),
			"0 */15 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 23, 15, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 3, 31, 23, 30, 0, 0, locs[1]),
				time.Date(1981, 3, 31, 23, 45, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 15, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s quarter-hourly leap unskip", locs[1]),
			"0 */15 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 23, 15, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 3, 31, 23, 30, 0, 0, locs[1]),
				time.Date(1981, 3, 31, 23, 45, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 15, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s daily third-hourly leap skip", locs[1]),
			"0 */20 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 2, 40, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 2, 2, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 2, 20, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 2, 40, 0, 0, locs[1]),
				time.Date(1981, 4, 3, 2, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s daily third-hourly leap unskip", locs[1]),
			"0 */20 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(1981, 3, 31, 2, 40, 0, 0, locs[1]),
			[]time.Time{
				time.Date(1981, 4, 1, 3, 0, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 20, 0, 0, locs[1]),
				time.Date(1981, 4, 1, 3, 40, 0, 0, locs[1]),
				time.Date(1981, 4, 2, 2, 0, 0, 0, locs[1]),
			},
		},
		{
			fmt.Sprintf("%s daily leap skip", locs[2]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 2, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 6, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 7, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 8, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 9, 2, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s daily leap unskip", locs[2]),
			"0 0 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 2, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 7, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 8, 2, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap skip", locs[2]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 15, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 7, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 8, 14, 5, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s time after daily leap unskip", locs[2]),
			"0 5 14 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 15, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 7, 14, 5, 0, 0, locs[2]),
				time.Date(2014, 10, 8, 14, 5, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap skip", locs[2]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 3, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 4, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 5, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 6, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s hourly leap unskip", locs[2]),
			"0 0 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 3, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 4, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 5, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s quarter-hourly leap skip", locs[2]),
			"0 */15 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 5, 1, 15, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 1, 45, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 2, 45, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s quarter-hourly leap unskip", locs[2]),
			"0 */15 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 5, 1, 15, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 1, 45, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 2, 45, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s third-hourly leap skip", locs[2]),
			"0 */20 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 2, 40, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 2, 40, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 20, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 40, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s third-hourly leap unskip", locs[2]),
			"0 */20 2 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireEarly},
			time.Date(2014, 10, 4, 2, 40, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2014, 10, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2014, 10, 5, 2, 40, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 0, 0, 0, locs[2]),
				time.Date(2014, 10, 6, 2, 20, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire early", locs[2]),
			"0 45 1 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]).Add(45 * time.Minute),
			[]time.Time{
				time.Date(2015, 4, 6, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 7, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 8, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 9, 1, 45, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire late", locs[2]),
			"0 45 1 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 1, 45, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 6, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 7, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 8, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 9, 1, 45, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s daily fall fire twice", locs[2]),
			"0 45 1 * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]).Add(45 * time.Minute),
				time.Date(2015, 4, 5, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 6, 1, 45, 0, 0, locs[2]),
				time.Date(2015, 4, 7, 1, 45, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire early", locs[2]),
			"0 30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2015, 4, 5, 0, 30, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 0, 30, 0, 0, locs[2]).Add(1 * time.Hour),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 3, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 4, 30, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire late", locs[2]),
			"0 30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 0, 30, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 3, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 4, 30, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s hourly fall fire twice", locs[2]),
			"0 30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 0, 30, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 0, 30, 0, 0, locs[2]).Add(1 * time.Hour),
				time.Date(2015, 4, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 3, 30, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s half-hourly fall fire early", locs[2]),
			"0 */30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly},
			time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]).Add(30 * time.Minute),
				time.Date(2015, 4, 5, 2, 0, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 3, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s half-hourly fall fire late", locs[2]),
			"0 */30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 0, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 3, 0, 0, 0, locs[2]),
			},
		},
		{
			fmt.Sprintf("%s half-hourly fall fire twice", locs[2]),
			"0 */30 * * * * *",
			cronexpr.Options{DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate},
			time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]),
			[]time.Time{
				time.Date(2015, 4, 5, 1, 0, 0, 0, locs[2]).Add(30 * time.Minute),
				time.Date(2015, 4, 5, 1, 30, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 0, 0, 0, locs[2]),
				time.Date(2015, 4, 5, 2, 30, 0, 0, locs[2]),
			},
		},
	}

	for _, tc := range cases {
		s, err := cronexpr.ParseWithOptions(tc.expr, tc.opts)
		if err != nil {
			t.Fatalf("parser error: %s", err)
		}

		runs := s.NextN(tc.from, 4)
		if len(runs) != 4 {
			t.Errorf("Case %s: Expected 4 runs, got %d", tc.name, len(runs))
		}

		for i := 0; i < len(runs); i++ {
			if !runs[i].Equal(tc.expected[i]) {
				t.Errorf("Case %s: Expected %v, got %v", tc.name, tc.expected[i], runs[i])
			}
		}
	}
}

// Issue: https://github.com/gorhill/cronexpr/issues/16
func TestInterval_Interval60Issue(t *testing.T) {
	_, err := cronexpr.Parse("*/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = cronexpr.Parse("*/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 61 should return err")
	}

	_, err = cronexpr.Parse("2/60 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}

	_, err = cronexpr.Parse("2-20/61 * * * * *")
	if err == nil {
		t.Errorf("parsing with interval 60 should return err")
	}
}

/******************************************************************************/

var benchmarkExpressions = []string{
	"* * * * *",
	"@hourly",
	"@weekly",
	"@yearly",
	"30 3 15W 3/3 *",
	"30 0 0 1-31/5 Oct-Dec * 2000,2006,2008,2013-2015",
	"0 0 0 * Feb-Nov/2 thu#3 2000-2050",
}
var benchmarkExpressionsLen = len(benchmarkExpressions)

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cronexpr.MustParse(benchmarkExpressions[i%benchmarkExpressionsLen])
	}
}

func BenchmarkNext(b *testing.B) {
	exprs := make([]*cronexpr.Expression, benchmarkExpressionsLen)
	for i := 0; i < benchmarkExpressionsLen; i++ {
		exprs[i] = cronexpr.MustParse(benchmarkExpressions[i])
	}
	from := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := exprs[i%benchmarkExpressionsLen]
		next := expr.Next(from)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
		next = expr.Next(next)
	}
}
