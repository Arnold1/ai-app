package build

import (
	"fmt"
	t "time"
)

var number string
var time string
var sha string
var startupTime t.Time

//Build represents the build number (if supplied) of the project
type Build struct {
	Number    string `json:"number"`
	Time      string `json:"time"`
	SHA       string `json:"sha"`
	StartedAt string `json:"startedAt"`
	Uptime    string `json:"uptime"`
}

func (build *Build) String() string {
	if build.Number != "" {
		return fmt.Sprintf("Build %s built on: %s (%s)\n ", build.Number, build.Time, build.SHA)
	}

	return "No Build Info"
}

func Init() {
	startupTime = t.Now()
}

func formatSince(start t.Time) string {
	const (
		Decisecond = 100 * t.Millisecond
		Day        = 24 * t.Hour
	)
	ts := t.Since(start)
	sign := t.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / t.Hour
	ts = ts % t.Hour
	m := ts / t.Minute
	ts = ts % t.Minute
	s := ts / t.Second
	ts = ts % t.Second
	f := ts / Decisecond
	return fmt.Sprintf("%dd%dh%dm%d.%ds", d, h, m, s, f)
}

//Info generates a build struct based on the ld flag based variables
func Info() *Build {
	return &Build{Number: number, Time: time, SHA: sha, StartedAt: startupTime.String(), Uptime: formatSince(startupTime)}
}