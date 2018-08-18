/*!
 * Copyright 2013 Raymond Hill
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr_next.go
 * Version: 1.0
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

package cronexpr

/******************************************************************************/

import (
	"sort"
	"time"
)

/******************************************************************************/

var dowNormalizedOffsets = [][]int{
	{1, 8, 15, 22, 29},
	{2, 9, 16, 23, 30},
	{3, 10, 17, 24, 31},
	{4, 11, 18, 25},
	{5, 12, 19, 26},
	{6, 13, 20, 27},
	{7, 14, 21, 28},
}

/******************************************************************************/

func (expr *Expression) nextYear(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate year
	i := sort.SearchInts(expr.YearList, t.Year()+1)
	if i == len(expr.YearList) {
		return time.Time{}
	}
	// Year changed, need to recalculate actual days of month
	expr.ActualDaysOfMonthList = expr.calculateActualDaysOfMonth(expr.YearList[i], expr.MonthList[0])
	if len(expr.ActualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			expr.YearList[i],
			time.Month(expr.MonthList[0]),
			1,
			expr.HourList[0],
			expr.MinuteList[0],
			expr.SecondList[0],
			0,
			t.Location()))
	}
	return time.Date(
		expr.YearList[i],
		time.Month(expr.MonthList[0]),
		expr.ActualDaysOfMonthList[0],
		expr.HourList[0],
		expr.MinuteList[0],
		expr.SecondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextMonth(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate month
	i := sort.SearchInts(expr.MonthList, int(t.Month())+1)
	if i == len(expr.MonthList) {
		return expr.nextYear(t)
	}
	// Month changed, need to recalculate actual days of month
	expr.ActualDaysOfMonthList = expr.calculateActualDaysOfMonth(t.Year(), expr.MonthList[i])
	if len(expr.ActualDaysOfMonthList) == 0 {
		return expr.nextMonth(time.Date(
			t.Year(),
			time.Month(expr.MonthList[i]),
			1,
			expr.HourList[0],
			expr.MinuteList[0],
			expr.SecondList[0],
			0,
			t.Location()))
	}

	return time.Date(
		t.Year(),
		time.Month(expr.MonthList[i]),
		expr.ActualDaysOfMonthList[0],
		expr.HourList[0],
		expr.MinuteList[0],
		expr.SecondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextDayOfMonth(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate day of month
	i := sort.SearchInts(expr.ActualDaysOfMonthList, t.Day()+1)
	if i == len(expr.ActualDaysOfMonthList) {
		return expr.nextMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		expr.ActualDaysOfMonthList[i],
		expr.HourList[0],
		expr.MinuteList[0],
		expr.SecondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextHour(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate hour
	i := sort.SearchInts(expr.HourList, t.Hour()+1)
	if i == len(expr.HourList) {
		return expr.nextDayOfMonth(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		expr.HourList[i],
		expr.MinuteList[0],
		expr.SecondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextMinute(t time.Time) time.Time {
	// Find index at which item in list is greater or equal to
	// candidate minute
	i := sort.SearchInts(expr.MinuteList, t.Minute()+1)
	if i == len(expr.MinuteList) {
		return expr.nextHour(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		expr.MinuteList[i],
		expr.SecondList[0],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) nextSecond(t time.Time) time.Time {
	// nextSecond() assumes all other fields are exactly matched
	// to the cron expression

	// Find index at which item in list is greater or equal to
	// candidate second
	i := sort.SearchInts(expr.SecondList, t.Second()+1)
	if i == len(expr.SecondList) {
		return expr.nextMinute(t)
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		expr.SecondList[i],
		0,
		t.Location())
}

/******************************************************************************/

func (expr *Expression) calculateActualDaysOfMonth(year, month int) []int {
	actualDaysOfMonthMap := make(map[int]bool)
	firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	// As per crontab man page (http://linux.die.net/man/5/crontab#):
	//  "The day of a command's execution can be specified by two
	//  "fields - day of month, and day of week. If both fields are
	//  "restricted (ie, aren't *), the command will be run when
	//  "either field matches the current time"

	// If both fields are not restricted, all days of the month are a hit
	if expr.DaysOfMonthRestricted == false && expr.DaysOfWeekRestricted == false {
		return genericDefaultList[1 : lastDayOfMonth.Day()+1]
	}

	// day-of-month != `*`
	if expr.DaysOfMonthRestricted {
		// Last day of month
		if expr.LastDayOfMonth {
			actualDaysOfMonthMap[lastDayOfMonth.Day()] = true
		}
		// Last work day of month
		if expr.LastWorkdayOfMonth {
			actualDaysOfMonthMap[workdayOfMonth(lastDayOfMonth, lastDayOfMonth)] = true
		}
		// Days of month
		for v := range expr.DaysOfMonth {
			// Ignore days beyond end of month
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
		// Work days of month
		// As per Wikipedia: month boundaries are not crossed.
		for v := range expr.WorkdaysOfMonth {
			// Ignore days beyond end of month
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[workdayOfMonth(firstDayOfMonth.AddDate(0, 0, v-1), lastDayOfMonth)] = true
			}
		}
	}

	// day-of-week != `*`
	if expr.DaysOfWeekRestricted {
		// How far first sunday is from first day of month
		offset := 7 - int(firstDayOfMonth.Weekday())
		// days of week
		//  offset : (7 - day_of_week_of_1st_day_of_month)
		//  target : 1 + (7 * week_of_month) + (offset + day_of_week) % 7
		for v := range expr.DaysOfWeek {
			w := dowNormalizedOffsets[(offset+v)%7]
			actualDaysOfMonthMap[w[0]] = true
			actualDaysOfMonthMap[w[1]] = true
			actualDaysOfMonthMap[w[2]] = true
			actualDaysOfMonthMap[w[3]] = true
			if len(w) > 4 && w[4] <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[w[4]] = true
			}
		}
		// days of week of specific week in the month
		//  offset : (7 - day_of_week_of_1st_day_of_month)
		//  target : 1 + (7 * week_of_month) + (offset + day_of_week) % 7
		for v := range expr.SpecificWeekDaysOfWeek {
			v = 1 + 7*(v/7) + (offset+v)%7
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
		// Last days of week of the month
		lastWeekOrigin := firstDayOfMonth.AddDate(0, 1, -7)
		offset = 7 - int(lastWeekOrigin.Weekday())
		for v := range expr.LastWeekDaysOfWeek {
			v = lastWeekOrigin.Day() + (offset+v)%7
			if v <= lastDayOfMonth.Day() {
				actualDaysOfMonthMap[v] = true
			}
		}
	}

	return toList(actualDaysOfMonthMap)
}

func workdayOfMonth(targetDom, lastDom time.Time) int {
	// If saturday, then friday
	// If sunday, then monday
	dom := targetDom.Day()
	dow := targetDom.Weekday()
	if dow == time.Saturday {
		if dom > 1 {
			dom -= 1
		} else {
			dom += 2
		}
	} else if dow == time.Sunday {
		if dom < lastDom.Day() {
			dom += 1
		} else {
			dom -= 2
		}
	}
	return dom
}
