package xlog

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	RegisterHelper("ago", func(t time.Time) string {
		if READONLY {
			return t.Format("Monday 2 January 2006")
		} else {
			return ago(time.Since(t))
		}
	})
	RegisterHelper("isFontAwesome", func(i string) bool {
		return len(i) > 3 && i[0:3] == "fa-"
	})
}

// RegisterHelper registers a new helper function. all helpers are used when compiling
// templates. so registering helpers function must happen before the server
// starts as compiling templates happened right before starting the http server.
func RegisterHelper(name string, f interface{}) {
	if _, ok := helpers[name]; ok {
		log.Fatalf("Helper: %s already registered", name)
	}

	helpers[name] = f
}

// A function that takes time.duration and return a string representation of the
// duration in human readable way such as "3 seconds ago". "5 hours 30 minutes
// ago". The precision of this function is 2. which means it returns the largest
// unit of time possible and the next one after it. for example days + hours, or
// Hours + minutes or Minutes + seconds...etc
func ago(t time.Duration) string {
	const day = time.Hour * 24
	const week = day * 7
	const month = day * 30
	const year = day * 365
	const maxPrecision = 2

	var o strings.Builder

	if t.Seconds() < 1 {
		o.WriteString("Less than a second ")
	}

	for precision := 0; t.Seconds() > 1 && precision < maxPrecision; precision++ {
		switch {
		case t >= year:
			years := t / year
			t -= years * year
			o.WriteString(fmt.Sprintf("%d years ", years))
		case t >= month:
			months := t / month
			t -= months * month
			o.WriteString(fmt.Sprintf("%d months ", months))
		case t >= week:
			weeks := t / week
			t -= weeks * week
			o.WriteString(fmt.Sprintf("%d weeks ", weeks))
		case t >= day:
			days := t / day
			t -= days * day
			o.WriteString(fmt.Sprintf("%d days ", days))
		case t >= time.Hour:
			hours := t / time.Hour
			t -= hours * time.Hour
			o.WriteString(fmt.Sprintf("%d hours ", hours))
		case t >= time.Minute:
			minutes := t / time.Minute
			t -= minutes * time.Minute
			o.WriteString(fmt.Sprintf("%d minutes ", minutes))
		case t >= time.Second:
			seconds := t / time.Second
			t -= seconds * time.Second
			o.WriteString(fmt.Sprintf("%d seconds ", seconds))
		}
	}

	o.WriteString("ago")

	return o.String()
}
