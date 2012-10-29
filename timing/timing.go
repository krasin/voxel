package timing

import (
	"fmt"
	"io"
	"sort"
	"time"
)

type Timing struct {
	Name  string
	Start time.Time
	End   time.Time
}

type TimingSlice []Timing

func (s TimingSlice) Len() int           { return len(s) }
func (s TimingSlice) Less(i, j int) bool { return s[i].End.Before(s[j].End) }
func (s TimingSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var timings = make(map[string]Timing)

func StartTiming(name string) {
	timings[name] = Timing{Name: name, Start: time.Now()}
}

func StopTiming(name string) {
	if t, ok := timings[name]; !ok {
		panic(fmt.Sprintf("StopTiming: unknown timing name: %s", name))
	} else {
		t.End = time.Now()
		timings[name] = t
	}
}

func (t Timing) Duration() time.Duration {
	return t.End.Sub(t.Start)
}

func GetTimings() (res []Timing) {
	for _, t := range timings {
		res = append(res, t)
	}
	sort.Sort(TimingSlice(res))
	return
}

func PrintTimings(w io.Writer) {
	fmt.Fprintf(w, "Timings:\n")
	for _, t := range GetTimings() {
		fmt.Fprintf(w, "\t%s:\t%d ms\n", t.Name, t.Duration().Nanoseconds()/1000000)
	}
}
