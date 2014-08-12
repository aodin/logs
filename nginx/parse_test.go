package nginx

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	records, err := OpenAndParseFile("./nginx_access_examples.log")
	if err != nil {
		t.Fatalf("Error during ParseFile(): %s", err)
	}
	if len(records) != 10 {
		t.Fatalf("Unexpected number of ParseFile() records: %d", len(records))
	}

	// Test an individual record's fields
	r0 := records[0]
	if r0.Status != 400 {
		t.Errorf("Unexpected status from ParseFile(): %d", r0.Status)
	}
	x := `41.227.38.172 - - [14/Nov/2013:06:59:03 +0000] "-" 400 0 "-" "-"`
	output := r0.String()
	if output != x {
		t.Errorf("Unexpected String() output of a Record: %s", output)
	}
	expectedTime := time.Date(2013, time.November, 14, 6, 59, 3, 0, time.UTC)
	if r0.Timestamp != expectedTime {
		t.Errorf(
			"Unexpected timestamp from ParseFile(): %s != %s",
			r0.Timestamp,
			expectedTime,
		)
	}
}
