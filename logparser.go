package logparser

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

// LogParser will parse the logs
type LogParser struct {
	r                    io.Reader
	logData              []LogLine
	TotalUniqueHits      int    `json:"total_unique_hits"`
	TotalUniqueSpiders   int    `json:"total_unique_spiders"`
	NumberOfDays         int    `json:"num_days"`
	UniqueHitsPerDay     int    `json:"unique_hits_per_day"`
	BandwidthLast24Hours string `json:"bandwidth_last_24_hours"`
	bandwidthLast24Hours int
	Name                 string `json:"name,omitempty"`
}

type LogLine struct {
	IP     string
	Time   time.Time
	Status int
	Method string
	HTTP   string
	Size   int
	Route  string
}

// OptionName sets the name
func OptionName(name string) func(*LogParser) error {
	return func(lp *LogParser) error {
		lp.Name = name
		return nil
	}
}

// New returns a new LogParser
func New(r io.Reader, options ...func(*LogParser) error) (lp *LogParser, err error) {
	lp = new(LogParser)
	lp.r = r

	for _, option := range options {
		err = option(lp)
		if err != nil {
			return
		}
	}
	return
}

func (lp *LogParser) Parse() (err error) {
	// gather data
	err = lp.parseReader()
	if err != nil {
		return
	}

	// gather stats
	err = lp.getStats()
	if err != nil {
		return
	}
	return
}

func (lp *LogParser) getStats() (err error) {
	// reset data
	lp.bandwidthLast24Hours = 0

	// get total number of days
	lp.NumberOfDays = int(lp.logData[len(lp.logData)-1].Time.Sub(lp.logData[0].Time).Hours() / 24)
	lastDay := lp.logData[len(lp.logData)-1].Time

	uniqueIPs := make(map[string]int)
	spiderIPs := make(map[string]struct{})
	for _, data := range lp.logData {
		// get IP count
		if _, ok := uniqueIPs[data.IP]; !ok {
			uniqueIPs[data.IP] = 0
		}
		uniqueIPs[data.IP]++
		if lastDay.Sub(data.Time).Hours() < 24 {
			lp.bandwidthLast24Hours += data.Size
		}

		// check for spiders
		if strings.Contains(data.Route, "robots.txt") {
			spiderIPs[data.IP] = struct{}{}
		}
	}
	lp.TotalUniqueHits = len(uniqueIPs)
	lp.UniqueHitsPerDay = lp.TotalUniqueHits / lp.NumberOfDays
	lp.BandwidthLast24Hours = humanize.Bytes(uint64(lp.bandwidthLast24Hours))
	lp.TotalUniqueSpiders = len(spiderIPs)
	return
}

// Parse will parse the loaded data
func (lp *LogParser) parseReader() (err error) {
	// count lines
	var lines = 0
	var teeBuffer bytes.Buffer
	tee := io.TeeReader(lp.r, &teeBuffer)
	scannerLines := bufio.NewScanner(tee)
	for scannerLines.Scan() {
		lines++
	}
	err = scannerLines.Err()
	if err != nil {
		return
	}

	// parse lines
	lp.logData = make([]LogLine, lines)
	scanner := bufio.NewScanner(&teeBuffer)
	i := 0
	for scanner.Scan() {
		lp.logData[i], err = parseCommon(scanner.Text())
		if err != nil {
			log.Println("err:", err.Error())
			continue
		}
		i++
	}
	err = scanner.Err()
	return
}

const layout = "02/Jan/2006:15:04:05 -0700" // 01/02 03:04:05PM '06 -0700
// parseCommon parses log of the form:
// 66.1.2.3 - - [22/May/2018:15:52:57 -0600] "GET /book HTTP/1.1" 403 38
func parseCommon(s string) (ll LogLine, err error) {
	fields := strings.Fields(s)
	ll.IP = fields[0]
	dateString := fields[3] + " " + fields[4]
	ll.Time, err = time.Parse(layout, dateString[1:len(dateString)-1])
	if err != nil {
		ll.Time, err = time.Parse(strings.Fields(layout)[0], fields[3][1:])
		if err != nil {
			err = errors.Wrap(err, "could not get time")
			return
		}
	}
	ll.Method = fields[5][1:]
	ll.Route = fields[6]
	ll.HTTP = fields[7][:len(fields[7])-1]
	ll.Status, err = strconv.Atoi(fields[8])
	if err != nil {
		err = errors.Wrap(err, "could not get status")
		return
	}
	ll.Size, err = strconv.Atoi(fields[9])
	if err != nil {
		err = errors.Wrap(err, "could not get size")
		return
	}
	return
}
