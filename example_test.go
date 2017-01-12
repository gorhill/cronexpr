/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/example_test.go
 * File: example_test.go
 * Version: 1.0
 * License: GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *
 */

package cronexpr_test

/******************************************************************************/

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

/******************************************************************************/

// ExampleMustParse
func ExampleMustParse() {
	t := time.Date(2013, time.August, 31, 0, 0, 0, 0, time.UTC)
	nextTimes := cronexpr.MustParse("0 0 29 2 *").NextN(t, 5)
	for i := range nextTimes {
		fmt.Println(nextTimes[i].Format(time.RFC1123))
		// Output:
		// Mon, 29 Feb 2016 00:00:00 UTC
		// Sat, 29 Feb 2020 00:00:00 UTC
		// Thu, 29 Feb 2024 00:00:00 UTC
		// Tue, 29 Feb 2028 00:00:00 UTC
		// Sun, 29 Feb 2032 00:00:00 UTC
	}
}

// Configure the parser to skip times in DST leaps and
// repeat times in DST falls
func ExampleParseWithOptions_dst1() {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Date(2014, 3, 8, 1, 0, 0, 0, loc)
	expr, _ := cronexpr.ParseWithOptions("0 0 2 * * * *", cronexpr.Options{
		DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate,
	})

	fmt.Println("DST leap times:")
	nextTimes := expr.NextN(t, 4)
	for i := range nextTimes {
		fmt.Println(nextTimes[i].Format(time.RFC1123))
	}

	t = time.Date(2014, 10, 31, 1, 0, 0, 0, loc)
	expr, _ = cronexpr.ParseWithOptions("0 0 1 * * * *", cronexpr.Options{
		DSTFlags: cronexpr.DSTFallFireEarly | cronexpr.DSTFallFireLate,
	})

	fmt.Println("DST fall times:")
	nextTimes = expr.NextN(t, 4)
	for i := range nextTimes {
		fmt.Println(nextTimes[i].Format(time.RFC1123))
	}

	// Output:
	// DST leap times:
	// Sat, 08 Mar 2014 02:00:00 PST
	// Mon, 10 Mar 2014 02:00:00 PDT
	// Tue, 11 Mar 2014 02:00:00 PDT
	// Wed, 12 Mar 2014 02:00:00 PDT
	// DST fall times:
	// Sat, 01 Nov 2014 01:00:00 PDT
	// Sun, 02 Nov 2014 01:00:00 PDT
	// Sun, 02 Nov 2014 01:00:00 PST
	// Mon, 03 Nov 2014 01:00:00 PST
}

// Configure the parser to unskip times in DST leaps and
// fire late in DST falls
func ExampleParseWithOptions_dst2() {
	loc, _ := time.LoadLocation("America/Los_Angeles")
	t := time.Date(2014, 3, 8, 1, 0, 0, 0, loc)
	expr, _ := cronexpr.ParseWithOptions("0 0 2 * * * *", cronexpr.Options{
		DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireLate,
	})

	fmt.Println("DST leap times:")
	nextTimes := expr.NextN(t, 4)
	for i := range nextTimes {
		fmt.Println(nextTimes[i].Format(time.RFC1123))
	}

	t = time.Date(2014, 10, 31, 1, 0, 0, 0, loc)
	expr, _ = cronexpr.ParseWithOptions("0 0 1 * * * *", cronexpr.Options{
		DSTFlags: cronexpr.DSTLeapUnskip | cronexpr.DSTFallFireLate,
	})

	fmt.Println("DST fall times:")
	nextTimes = expr.NextN(t, 4)
	for i := range nextTimes {
		fmt.Println(nextTimes[i].Format(time.RFC1123))
	}

	// Output:
	// DST leap times:
	// Sat, 08 Mar 2014 02:00:00 PST
	// Sun, 09 Mar 2014 03:00:00 PDT
	// Mon, 10 Mar 2014 02:00:00 PDT
	// Tue, 11 Mar 2014 02:00:00 PDT
	// DST fall times:
	// Sat, 01 Nov 2014 01:00:00 PDT
	// Sun, 02 Nov 2014 01:00:00 PST
	// Mon, 03 Nov 2014 01:00:00 PST
	// Tue, 04 Nov 2014 01:00:00 PST
}
