package logs

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

type Log struct {
	Origin string
	Time   time.Time
	Method string
	Path   string
	Status int
	Size   int
}

var (
	// fina.ansremote.com - - [06/Aug/2022:15:29:31 +0000] "GET /icons/blank.xbm HTTP/1.1" 304 0
	// msp1-16.nas.mr.net - - [06/Aug/2022:02:38:43 +0000] "GET /images/nasa-logo.gif HTTP/1.1" 404 -
	logRe = regexp.MustCompile(`(.*) - - \[(.*)\] "([A-Z]+) ([^ ]+).*?" (\d+) (-|\d+)`)
)

// ParseLine returns parse log from line.
func ParseLine(line string) (Log, error) {
	matches := logRe.FindStringSubmatch(line)

	origin := matches[1]
	time, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return Log{}, fmt.Errorf("can't parse time - %w", err)
	}
	method := matches[3]
	path := matches[4]
	status, err := strconv.Atoi(matches[5])
	if err != nil {
		return Log{}, fmt.Errorf("bad status: %q - %w", matches[5], err)
	}
	size := -1
	if matches[6] != "-" {
		size, err = strconv.Atoi(matches[6])
		if err != nil {
			return Log{}, fmt.Errorf("bad size: %q - %w", matches[6], err)
		}
	}

	log := Log{
		Origin: origin,
		Time:   time.UTC(),
		Method: method,
		Path:   path,
		Status: status,
		Size:   size,
	}
	return log, nil
}

// Scanner scans logs from reader.
type Scanner struct {
	scan *bufio.Scanner
	lnum int
	log  Log
	err  error
}

// NewScanner return new scanner from r.
func NewScanner(r io.Reader) *Scanner {
	s := Scanner{
		scan: bufio.NewScanner(r),
	}
	return &s
}

// Log returns the current log.
func (s *Scanner) Log() Log {
	return s.log
}

// Line returns the current line.
func (s *Scanner) Line() int {
	return s.lnum
}

// Next moves to the next log.
func (s *Scanner) Next() bool {
	if !s.scan.Scan() {
		return false
	}

	s.lnum++
	line := s.scan.Text()
	if line == "" { // Skip empty lines
		return s.Next()
	}

	s.log, s.err = ParseLine(line)
	if s.err != nil {
		return false
	}
	return true
}

// Err returns the current error.
func (s *Scanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scan.Err()
}
