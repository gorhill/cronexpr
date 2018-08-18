/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr.go
 * Version: 1.0
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

// Package cronexpr parses cron time expressions.
package cronexpr

/******************************************************************************/

import (
	"fmt"
	"sort"
	"time"
)

/******************************************************************************/

// A Expression represents a specific cron time expression as defined at
// <https://github.com/gorhill/cronexpr#implementation>
type Expression struct {
	Expression             string
	SecondList             []int
	MinuteList             []int
	HourList               []int
	DaysOfMonth            map[int]bool
	WorkdaysOfMonth        map[int]bool
	LastDayOfMonth         bool
	LastWorkdayOfMonth     bool
	DaysOfMonthRestricted  bool
	ActualDaysOfMonthList  []int
	MonthList              []int
	DaysOfWeek             map[int]bool
	SpecificWeekDaysOfWeek map[int]bool
	LastWeekDaysOfWeek     map[int]bool
	DaysOfWeekRestricted   bool
	YearList               []int
}

/******************************************************************************/

// MustParse returns a new Expression pointer. It expects a well-formed cron
// expression. If a malformed cron expression is supplied, it will `panic`.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
func MustParse(cronLine string) *Expression {
	expr, err := Parse(cronLine)
	if err != nil {
		panic(err)
	}
	return expr
}

/******************************************************************************/

// Parse returns a new Expression pointer. An error is returned if a malformed
// cron expression is supplied.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
func Parse(cronLine string) (*Expression, error) {

	// Maybe one of the built-in aliases is being used
	cron := cronNormalizer.Replace(cronLine)

	indices := fieldFinder.FindAllStringIndex(cron, -1)
	fieldCount := len(indices)
	if fieldCount < 5 {
		return nil, fmt.Errorf("missing field(s)")
	}
	// ignore fields beyond 7th
	if fieldCount > 7 {
		fieldCount = 7
	}

	var expr = Expression{}
	var field = 0
	var err error

	// second field (optional)
	if fieldCount == 7 {
		err = expr.secondFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
		field += 1
	} else {
		expr.SecondList = []int{0}
	}

	// minute field
	err = expr.minuteFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// hour field
	err = expr.hourFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of month field
	err = expr.domFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// month field
	err = expr.monthFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of week field
	err = expr.dowFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// year field
	if field < fieldCount {
		err = expr.yearFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
	} else {
		expr.YearList = yearDescriptor.defaultList
	}

	return &expr, nil
}

/******************************************************************************/

// Next returns the closest time instant immediately following `fromTime` which
// matches the cron expression `expr`.
//
// The `time.Location` of the returned time instant is the same as that of
// `fromTime`.
//
// The zero value of time.Time is returned if no matching time instant exists
// or if a `fromTime` is itself a zero value.
func (expr *Expression) Next(fromTime time.Time) time.Time {
	// Special case
	if fromTime.IsZero() {
		return fromTime
	}

	// Since expr.nextSecond()-expr.nextMonth() expects that the
	// supplied time stamp is a perfect match to the underlying cron
	// expression, and since this function is an entry point where `fromTime`
	// does not necessarily matches the underlying cron expression,
	// we first need to ensure supplied time stamp matches
	// the cron expression. If not, this means the supplied time
	// stamp falls in between matching time stamps, thus we move
	// to closest future matching immediately upon encountering a mismatching
	// time stamp.

	// year
	v := fromTime.Year()
	i := sort.SearchInts(expr.YearList, v)
	if i == len(expr.YearList) {
		return time.Time{}
	}
	if v != expr.YearList[i] {
		return expr.nextYear(fromTime)
	}
	// month
	v = int(fromTime.Month())
	i = sort.SearchInts(expr.MonthList, v)
	if i == len(expr.MonthList) {
		return expr.nextYear(fromTime)
	}
	if v != expr.MonthList[i] {
		return expr.nextMonth(fromTime)
	}

	expr.ActualDaysOfMonthList = expr.calculateActualDaysOfMonth(fromTime.Year(), int(fromTime.Month()))
	if len(expr.ActualDaysOfMonthList) == 0 {
		return expr.nextMonth(fromTime)
	}

	// day of month
	v = fromTime.Day()
	i = sort.SearchInts(expr.ActualDaysOfMonthList, v)
	if i == len(expr.ActualDaysOfMonthList) {
		return expr.nextMonth(fromTime)
	}
	if v != expr.ActualDaysOfMonthList[i] {
		return expr.nextDayOfMonth(fromTime)
	}
	// hour
	v = fromTime.Hour()
	i = sort.SearchInts(expr.HourList, v)
	if i == len(expr.HourList) {
		return expr.nextDayOfMonth(fromTime)
	}
	if v != expr.HourList[i] {
		return expr.nextHour(fromTime)
	}
	// minute
	v = fromTime.Minute()
	i = sort.SearchInts(expr.MinuteList, v)
	if i == len(expr.MinuteList) {
		return expr.nextHour(fromTime)
	}
	if v != expr.MinuteList[i] {
		return expr.nextMinute(fromTime)
	}
	// second
	v = fromTime.Second()
	i = sort.SearchInts(expr.SecondList, v)
	if i == len(expr.SecondList) {
		return expr.nextMinute(fromTime)
	}

	// If we reach this point, there is nothing better to do
	// than to move to the next second

	return expr.nextSecond(fromTime)
}

/******************************************************************************/

// NextN returns a slice of `n` closest time instants immediately following
// `fromTime` which match the cron expression `expr`.
//
// The time instants in the returned slice are in chronological ascending order.
// The `time.Location` of the returned time instants is the same as that of
// `fromTime`.
//
// A slice with len between [0-`n`] is returned, that is, if not enough existing
// matching time instants exist, the number of returned entries will be less
// than `n`.
func (expr *Expression) NextN(fromTime time.Time, n uint) []time.Time {
	nextTimes := make([]time.Time, 0, n)
	if n > 0 {
		fromTime = expr.Next(fromTime)
		for {
			if fromTime.IsZero() {
				break
			}
			nextTimes = append(nextTimes, fromTime)
			n -= 1
			if n == 0 {
				break
			}
			fromTime = expr.nextSecond(fromTime)
		}
	}
	return nextTimes
}
