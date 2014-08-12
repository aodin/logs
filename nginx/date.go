package nginx

import (
    "fmt"
	"sort"
	"time"
)

// Date is a simple year, month, and day representation. Go's "time" package
// has only a datetime representation (which could be truncated to serve as a
// date) and I desired a simple implementation.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d Date) String() string {
    return fmt.Sprintf("%s %d, %d", d.Month, d.Day, d.Year)
}

// ToTime returns a UTC zoned time for this date
func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// DateFromTime constructs a date from a time
func DateFromTime(t time.Time) Date {
	return Date{t.Year(), t.Month(), t.Day()}
}

// Dates implements the sort.Sort interface. By default it will sort dates
// in ascending order.
type Dates []Date

func (d Dates) Len() int {
	return len(d)
}

func (d Dates) Less(i, j int) bool {
	if d[i].Year == d[j].Year {
		if d[i].Month == d[j].Month {
			return d[i].Day < d[j].Day
		}
		return d[i].Month < d[j].Month
	}
	return d[i].Year < d[j].Year
}

func (d Dates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// DateCount merely contains a date and count.
type DateCount struct {
	Date  Date
	Count int64
}

// DateCounter is used to count views for a given date.
type DateCounter map[Date]int64

// Range returns the ordered dates in the date counter, including any missing
// dates.
func (dc DateCounter) Range() (counts []DateCount) {
	// Skip sorting if there is at most one date
	if len(dc) < 2 {
		counts = make([]DateCount, len(dc))
		var i int
		for dd, count := range dc {
			counts[i] = DateCount{Date: dd, Count: count}
		}
		return
	}
	// Determine the min and max dates
	dates := make([]Date, len(dc))
	var i int
	for dd := range dc {
		dates[i] = dd
		i += 1
	}
	sort.Sort(Dates(dates))
	min := dates[0]
	max := dates[len(dates)-1]

	// Construct the range of dates
	minT := min.ToTime()
	maxT := max.ToTime()

	// This should be an even number!
	days := int(maxT.Sub(minT).Hours() / 24.0)

	// Fill in any missing dates
	counts = make([]DateCount, days+1)
	for i := 0; i <= days; i++ {
		dd := DateFromTime(minT.AddDate(0, 0, i))
		counts[i] = DateCount{dd, dc[dd]}
	}
	return
}
