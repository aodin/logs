package nginx

import (
	"sort"
	"testing"
	"time"
)

func expectDate(t *testing.T, a Date, i int, yy int, mm time.Month, dd int) {
	if a.Year != yy {
		t.Errorf("Unexpected year for index %d: %d", i, a.Year)
	}
	if a.Month != mm {
		t.Errorf("Unexpected month for index %d: %s", i, a.Month)
	}
	if a.Day != dd {
		t.Errorf("Unexpected day for index %d: %d", i, a.Day)
	}
}

// Test the sorting of dates
func TestDates(t *testing.T) {
	examples := []Date{
		{2014, 3, 4},
		{2014, 2, 4},
		{2014, 3, 2},
		{2014, 4, 4},
		{2015, 3, 4},
		{2013, 3, 4},
	}
	sort.Sort(Dates(examples))

	expectDate(t, examples[0], 0, 2013, 3, 4)
	expectDate(t, examples[2], 0, 2014, 3, 2)
	expectDate(t, examples[3], 0, 2014, 3, 4)
	expectDate(t, examples[5], 0, 2015, 3, 4)
}

func expectDateCount(t *testing.T, a DateCount, d Date, c int64) {
	dd := a.Date
	cc := a.Count
	if dd.Year != d.Year {
		t.Errorf("Unexpected year %d", dd.Year)
	}
	if dd.Month != d.Month {
		t.Errorf("Unexpected month %s", dd.Month)
	}
	if dd.Day != d.Day {
		t.Errorf("Unexpected day %d", dd.Day)
	}
	if cc != c {
		t.Errorf("Unexpected count: %d", cc)
	}
}

func TestDateCounter(t *testing.T) {
	feb27 := Date{2014, 2, 27}
	feb28 := Date{2014, 2, 28}
	mar1 := Date{2014, 3, 1}
	dc := DateCounter{}
	dc[feb27] += 1
	dc[mar1] += 1

	counts := dc.Range()
	if len(counts) != 3 {
		t.Fatalf("Unexpected length of counts: %d", len(counts))
	}
	expectDateCount(t, counts[0], feb27, 1)
	expectDateCount(t, counts[1], feb28, 0)
	expectDateCount(t, counts[2], mar1, 1)
}
