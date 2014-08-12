package nginx

import "time"

// CountDays will count the number of records that include the given URL and
// separate them by day in the given timezone.
func CountDays(records []Record, url string, loc *time.Location) (int64, []DateCount) {
	var total int64
	byDay := DateCounter{}
	for _, record := range records {
		if url == record.URL {
			byDay[DateFromTime(record.Timestamp.In(loc))] += 1
			total += 1
		}
	}
	return total, byDay.Range()
}

// Count will count the number of records that include the given URL.
func Count(records []Record, url string) (total int64) {
	for _, record := range records {
		if url == record.URL {
			total += 1
		}
	}
	return
}
