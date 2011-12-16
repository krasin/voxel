package main

import (
	"fmt"
	"io"
	"time"
)

var starts map[string]time.Time
var stops map[string]time.Time

func StartTiming(name string) {
	starts[name] = time.Now()
}

func StopTiming(name string) {
	if _, ok := starts[name]; !ok {
		panic(fmt.Sprintf("StopTiming: unknown timing name: %s", name))
	}
	stops[name] = time.Now()
}

func GetTimings() (res map[string]time.Duration) {
	for name, start := range starts {
		if _, ok := stops[name]; !ok {
			continue
		}
		res[name] = stops[name].Sub(start)
	}
	return
}

func PrintTimings(w io.Writer) {
	fmt.Fprintf(w, "%v\n", GetTimings())
}
