package nginx

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// nginxTimeLayout is the default nginx time layout.
const nginxTimeLayout = `[02/Jan/2006:15:04:05 -0700]`

// The default configuration for nginx access logs:
// $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
// Skip index 1, it's just a dash
type Record struct {
	IP        string
	User      string // index 2
	Timestamp time.Time
	Request   string
	Method    string
	URL       string
	Version   string
	Status    int
	Bytes     int
	Referer   string
	Agent     string // index 8
}

// String rebuilds the Record into the default nginx log format.
func (r Record) String() string {
	request := r.Request
	if request == "" {
		request = "-"
	}
	referer := r.Referer
	if referer == "" {
		referer = "-"
	}
	agent := r.Agent
	if agent == "" {
		agent = "-"
	}
	return fmt.Sprintf(`%s - %s %s "%s" %d %d "%s" "%s"`,
		r.IP,
		r.User,
		r.Timestamp.Format(nginxTimeLayout),
		request,
		r.Status,
		r.Bytes,
		referer,
		agent,
	)
}

func ParseRow(row []string) (r Record, err error) {
	// TODO Split the request into method / url / protocol?
	r.IP = row[0]
	r.User = row[2]
	// Resplice together the timestamp from columns 3 and 4
	ts := row[3] + " " + row[4]
	// Go's parse rules are quite dumb - it will fix the offset unless the
	// timezone is explicitly named
	if r.Timestamp, err = time.ParseInLocation(nginxTimeLayout, ts, time.UTC); err != nil {
		return
	}
	r.Request = row[5]
	// Attempt to split the request
	parts := strings.Split(r.Request, " ")

	// If there were three parts, separate the request
	if len(parts) == 3 {
		r.Method = parts[0]
		r.URL = parts[1]
		r.Version = parts[2]
	}

	status, err := strconv.ParseInt(row[6], 10, 32)
	if err != nil {
		return
	}
	r.Status = int(status)
	bytes, err := strconv.ParseInt(row[7], 10, 32)
	if err != nil {
		return
	}
	r.Bytes = int(bytes)
	r.Referer = row[8]
	r.Agent = row[9]
	return
}

func ParseFile(f io.ReadCloser) (records []Record, err error) {
	defer f.Close()
	// Nginx access logs are space separated values
	r := csv.NewReader(f)
	r.Comma = ' '

	// Convert to records
	var rows [][]string
	rows, err = r.ReadAll()
	if err != nil {
		return
	}

	records = make([]Record, len(rows))
	for i, row := range rows {
		record, parseErr := ParseRow(row)
		if parseErr != nil {
			// Wrap the parse error to give a line number
			err = fmt.Errorf("nginx: error parsing line %d: %s", i, parseErr)
			return
		}
		records[i] = record
	}
	return
}

func OpenAndParseFile(path string) (records []Record, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	return ParseFile(f)
}

func OpenAndParseZipFile(path string) (records []Record, err error) {
	gz, err := os.Open(path)
	if err != nil {
		return
	}
	// TODO Does the original file need to be closed?

	// Read the compressed file
	f, err := gzip.NewReader(gz)
	if err != nil {
		return
	}
	return ParseFile(f)
}

type Logs struct {
	Files         int
	Plain, Zipped []string
}

func ReadDirectory(dir string) (logs Logs, err error) {
	// TODO Recursive or skip sub-directories?
	// if file.IsDir() {
	logs.Plain = make([]string, 0)
	logs.Zipped = make([]string, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	logs.Files = len(files)

	for _, file := range files {
		name := file.Name()
		// Read only access logs
		if !strings.HasPrefix(name, "access") {
			continue
		}

		// TODO Use a regex for parsing?
		// The full filepath that will be used to open the file
		filepath := filepath.Join(dir, name)
		ext := strings.Split(strings.ToLower(name), ".")
		extLen := len(ext)
		if extLen < 2 {
			continue
		}
		if ext[extLen-1] == "log" || ext[extLen-2] == "log" {
			logs.Plain = append(logs.Plain, filepath)
		} else if ext[extLen-1] == "gz" {
			logs.Zipped = append(logs.Zipped, filepath)
		}
	}
	return
}
