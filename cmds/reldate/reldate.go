/**
 * Generates a date in YYYY-MM-DD format based on a relative time
 * description (e.g. -1 week, +3 years)
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2014 all rights reserved.
 * Released under the Simplified BSD License
 * See: http://opensource.org/licenses/bsd-license.php
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const yyyymmdd = "2006-01-02"

var (
	help          bool
	endOfMonthFor bool
	relativeTo    string
	relativeT     time.Time
)

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	fmt.Fprintf(fh, `%s
 USAGE %s [TIME_INCREMENT TIME_UNIT|WEEKDAY_NAME]
    

 EXAMPLES
 
 Two days from today: %s 2 days
 Three weeks ago: %s -- -3 weeks
 Three weeks from 2014-01-01: %s --from=2014-01-01 3 weeks
 Three days before 2014-01-01: %s --from=2014-01-01 -- -3 days
 The Friday of this week: %s Friday
 The Monday in week containing 2015-02-06: %s --from=2015-02-06 Monday

 Time increments are a positive or negative integer. Time unit can be
 either day(s), week(s), month(s), or year(s). Weekday names are
 case insentive (e.g. Monday and monday). They can be abbreviated
 to the first three letters of the name, e.g. Sunday can be Sun, Monday
 can be Mon, Tuesday can be Tue, Wednesday can be Wed, Thursday can
 be Thu, Friday can be Fri or Saturday can be Sat.

 OPTIONS

`, msg, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t(defaults to %s) %s\n", f.Name, f.Value, f.Usage)
	})

	fmt.Fprintf(fh, `

 copyright (c) 2014 all rights reserved.
 Released under the Simplified BSD License
 See: http://opensource.org/licenses/bsd-license.php

`)
	os.Exit(exit_code)
}

func endOfMonth(t1 time.Time) string {
	location := t1.Location()
	year := t1.Year()
	month := t1.Month()
	if month == 12 {
		year++
	}
	month++
	t2 := time.Date(year, month, 1, 0, 0, 0, 0, location)
	return t2.Add(-time.Hour).Format(yyyymmdd)
}

func init() {
	const (
		relativeToUsage = "Date the relative time is calculated from."
		helpUsage       = "Display this help document."
		endOfMonthUsage = "Display the end of month day. E.g. 2012-02-29"
	)

	flag.StringVar(&relativeTo, "from", relativeTo, relativeToUsage)
	flag.StringVar(&relativeTo, "f", relativeTo, relativeToUsage)
	flag.BoolVar(&endOfMonthFor, "end-of-month", endOfMonthFor, endOfMonthUsage)
	flag.BoolVar(&help, "help", help, helpUsage)
	flag.BoolVar(&help, "h", help, helpUsage)
}

func assertOk(e error, failMsg string) {
	if e != nil {
		usage(1, fmt.Sprintf(" %s\n %s\n", failMsg, e))
	}
}

func weekdayOffset(weekday time.Weekday) int {
	switch {
	case weekday == time.Sunday:
		return 0
	case weekday == time.Monday:
		return 1
	case weekday == time.Tuesday:
		return 2
	case weekday == time.Wednesday:
		return 3
	case weekday == time.Thursday:
		return 4
	case weekday == time.Friday:
		return 5
	case weekday == time.Saturday:
		return 6
	}
	return 0
}

func relativeWeekday(t time.Time, weekday time.Weekday) (time.Time, error) {
	// Normalize to Sunday then add weekday constant
	switch {
	case t.Weekday() == time.Sunday:
		return t.AddDate(0, 0, weekdayOffset(weekday)), nil
	case t.Weekday() == time.Monday:
		return t.AddDate(0, 0, (-1 + weekdayOffset(weekday))), nil
	case t.Weekday() == time.Tuesday:
		return t.AddDate(0, 0, (-2 + weekdayOffset(weekday))), nil
	case t.Weekday() == time.Wednesday:
		return t.AddDate(0, 0, (-3 + weekdayOffset(weekday))), nil
	case t.Weekday() == time.Thursday:
		return t.AddDate(0, 0, (-4 + weekdayOffset(weekday))), nil
	case t.Weekday() == time.Friday:
		return t.AddDate(0, 0, (-5 + weekdayOffset(weekday))), nil
	case t.Weekday() == time.Saturday:
		return t.AddDate(0, 0, (-6 + weekdayOffset(weekday))), nil
	}
	return t, errors.New("Expecting Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, or Saturday.")
}

func relativeTime(t time.Time, i int, u string) (time.Time, error) {
	switch {
	case strings.HasPrefix(u, "sun"):
		return relativeWeekday(t, time.Sunday)
	case strings.HasPrefix(u, "mon"):
		return relativeWeekday(t, time.Monday)
	case strings.HasPrefix(u, "tue"):
		return relativeWeekday(t, time.Tuesday)
	case strings.HasPrefix(u, "wed"):
		return relativeWeekday(t, time.Wednesday)
	case strings.HasPrefix(u, "thu"):
		return relativeWeekday(t, time.Thursday)
	case strings.HasPrefix(u, "fri"):
		return relativeWeekday(t, time.Friday)
	case strings.HasPrefix(u, "sat"):
		return relativeWeekday(t, time.Saturday)
	case strings.HasPrefix(u, "day"):
		return t.AddDate(0, 0, i), nil
	case strings.HasPrefix(u, "week"):
		return t.AddDate(0, 0, 7*i), nil
	case strings.HasPrefix(u, "month"):
		return t.AddDate(0, i, 0), nil
	case strings.HasPrefix(u, "year"):
		return t.AddDate(i, 0, 0), nil
	}
	return t, errors.New("Time unit must be day(s), week(s), month(s) or year(s) or weekday name.")
}

func main() {
	var (
		err        error
		unitString string
	)

	flag.Parse()
	if help == true {
		usage(0, "")
	}

	argc := flag.NArg()
	argv := flag.Args()

	if argc < 1 && endOfMonthFor == false {
		usage(1, "Missing time increment and units (e.g. +2 days) or weekday name (e.g. Monday, Mon).\n")
	} else if argc > 2 {
		usage(1, "Too many command line arguments.\n")
	}

	relativeT = time.Now()
	if relativeTo != "" {
		relativeT, err = time.Parse(yyyymmdd, relativeTo)
		assertOk(err, "Cannot parse the from date.\n")
	}

	if endOfMonthFor == true {
		fmt.Println(endOfMonth(relativeT))
		os.Exit(0)
	}

	timeInc := 0
	if argc == 2 {
		unitString = strings.ToLower(argv[1])
		timeInc, err = strconv.Atoi(argv[0])
		assertOk(err, "Time increment should be a positive or negative integer.\n")
	} else {
		// We may have a weekday string
		unitString = strings.ToLower(argv[0])
	}
	t, err := relativeTime(relativeT, timeInc, unitString)
	assertOk(err, "Did not understand command.")
	fmt.Println(t.Format(yyyymmdd))
}
