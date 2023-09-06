package logs

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var lineCases = []struct {
	name string
	line string
	log  Log
}{
	{
		"regular",
		`edams.ksc.nasa.gov - - [01/Aug/2022:12:35:16 +0000] "GET /ksc.html HTTP/1.1" 200 7280`,
		Log{
			Origin: "edams.ksc.nasa.gov",
			Time:   time.Date(2022, time.August, 1, 12, 35, 16, 0, time.UTC),
			Method: "GET",
			Path:   "/ksc.html",
			Status: 200,
			Size:   7280,
		},
	},
	{
		"404",
		`198.150.176.207 - - [03/Aug/2022:19:09:03 +0000] "GET /history/apollo/apollo-13.html HTTP/1.1" 404 -`,
		Log{
			Origin: "198.150.176.207",
			Time:   time.Date(2022, time.August, 3, 19, 9, 3, 0, time.UTC),
			Method: "GET",
			Path:   "/history/apollo/apollo-13.html",
			Status: 404,
			Size:   -1,
		},
	},
	{
		"bob",
		`bos1b.delphi.com - - [03/Aug/2022:16:51:06 +0000] "GET /history/skylab/skylab-operations.txt" 200 13586`,
		Log{
			Origin: "bos1b.delphi.com",
			Time:   time.Date(2022, time.August, 3, 16, 51, 6, 0, time.UTC),
			Method: "GET",
			Path:   "/history/skylab/skylab-operations.txt",
			Status: 200,
			Size:   13586,
		},
	},
}

func TestParseLine(t *testing.T) {
	for _, tc := range lineCases {
		t.Run(tc.name, func(t *testing.T) {
			log, err := ParseLine(tc.line)
			require.NoError(t, err)
			require.Equal(t, tc.log, log)
		})
	}
}

var testData = `vagrant.vf.mmc.com - - [04/Aug/2022:11:18:50 +0000] "GET /shuttle/technology/sts-newsref/sts-gnnc.html HTTP/1.1" 200 119561
slip2.comserv.ipn.mx - - [03/Aug/2022:21:41:55 +0000] "GET / HTTP/1.1" 200 7034
nslemmer.sri.com - - [04/Aug/2022:19:51:52 +0000] "GET /icons/image.xbm HTTP/1.1" 200 509
168.200.39.38 - - [03/Aug/2022:14:52:31 +0000] "GET /history/apollo/flight-summary.txt HTTP/1.1" 200 5086
www-c3.proxy.aol.com - - [06/Aug/2022:14:41:04 +0000] "GET /images/WORLD-logosmall.gif HTTP/1.1" 200 669
screamer.raxco.com - - [04/Aug/2022:12:30:13 +0000] "GET /shuttle/technology/sts-newsref/sts-rcs.html HTTP/1.1" 200 109356
caboose-sl12.ironhorse.com - - [04/Aug/2022:22:30:06 +0000] "GET /shuttle/missions/sts-69/count69.gif HTTP/1.1" 200 46053
132.208.99.95 - - [05/Aug/2022:21:18:34 +0000] "GET /icons/sound.xbm HTTP/1.1" 200 530
130.110.74.81 - - [01/Aug/2022:08:13:39 +0000] "GET /shuttle/missions/sts-59/mission-sts-59.html HTTP/1.1" 200 63060
141.102.80.167 - - [04/Aug/2022:12:54:10 +0000] "GET /shuttle/missions/sts-69/news HTTP/1.1" 302 -

www-c3.proxy.aol.com - - [01/Aug/2022:01:03:14 +0000] "GET /shuttle/countdown/images/countclock.gif HTTP/1.1" 200 13994
dial-05.escape.ca - - [03/Aug/2022:12:55:08 +0000] "GET /shuttle/countdown/lps/fr.gif HTTP/1.1" 200 30232
www-c1.proxy.aol.com - - [04/Aug/2022:11:15:29 +0000] "GET /images/ HTTP/1.1" 200 17688
198.150.176.207 - - [03/Aug/2022:19:09:03 +0000] "GET /history/apollo/apollo-13.html HTTP/1.1" 404 -
`

func TestScanner(t *testing.T) {
	r := strings.NewReader(testData)
	s := NewScanner(r)
	count := 0
	for s.Next() {
		count++
		require.NotEqual(t, s.Log().Path, "")
	}
	require.NoError(t, s.Err())
	require.Equal(t, 14, count, "record count")
	require.Equal(t, 15, s.Line(), "line count")
}

func TestScannerFile(t *testing.T) {
	file, err := os.Open("testdata/httpd.log")
	require.NoError(t, err, "open")
	defer file.Close()

	s := NewScanner(file)
	count := 0
	for s.Next() {
		count++
		require.NotEqual(t, s.Log().Path, "")
	}
	require.NoError(t, s.Err())
	require.Equal(t, 1984, count, "record count")
}
