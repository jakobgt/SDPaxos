package main

import (
	"time"
	"github.com/olekukonko/tablewriter"
	"os"
	"fmt"
)

// _histoGranularityMs is the granularity of the historgram. This
// must be divisble by 1000.
const _histoGranularityMs = 5

// stats is a simple stats collection package
type stats struct {
	iterations   int
	successes    int
	errors       int
	errorMap     map[string]int
	latencyHisto []int
	limitLatency time.Duration
	minLatency   time.Duration
	maxLatency   time.Duration
	maxSec       int
	// The following is used to track correctness.
	entityMap map[string]int
	startBal  int
	endBal    int
}

func newStats(maxSec int) *stats {
	return &stats{
		errorMap:     map[string]int{},
		entityMap:    map[string]int{},
		latencyHisto: make([]int, maxSec*(1000/_histoGranularityMs)),
		limitLatency: time.Duration(maxSec) * time.Second,
		minLatency:   time.Duration(maxSec) * time.Second,
		maxSec:       maxSec,
	}
}

func (s *stats) update(err error, latency time.Duration) {
	s.iterations++
	if err != nil {
		s.addErrorMsg(err.Error(), 1)
		return
	}
	if latency >= s.limitLatency {
		s.addErrorMsg("elapsed too long", 1)
		return
	}
	s.successes++
	if s.minLatency > latency {
		s.minLatency = latency
	}
	if s.maxLatency < latency {
		s.maxLatency = latency
	}
	idx := int64(latency/time.Millisecond) / _histoGranularityMs
	s.latencyHisto[idx]++
}

func (s *stats) percentile(percentile int) time.Duration {
	total := 0
	for _, c := range s.latencyHisto {
		total += c
	}
	cutoff := (total * percentile) / 100
	count := 0
	for i := range s.latencyHisto {
		count += s.latencyHisto[i]
		if count > cutoff {
			return time.Duration((i+1)*_histoGranularityMs) * time.Millisecond
		}
	}
	return 0
}

func (s *stats) addErrorMsg(msg string, n int) {
	s.errors += n
	if c, ok := s.errorMap[msg]; ok {
		s.errorMap[msg] = c + n
	} else {
		s.errorMap[msg] = n
	}
}

func (s *stats) merge(o *stats) {
	s.iterations += o.iterations
	s.successes += o.successes
	for msg, count := range o.errorMap {
		s.addErrorMsg(msg, count)
	}
	if o.minLatency < s.minLatency {
		s.minLatency = o.minLatency
	}
	if o.maxLatency > s.maxLatency {
		s.maxLatency = o.maxLatency
	}
	for i, c := range o.latencyHisto {
		s.latencyHisto[i] += c
	}
	for user, bal := range o.entityMap {
		s.entityMap[user] += bal
	}
}

func (s *stats) asPercent(val int) float32 {
	return (float32(val) * 100.0) / float32(s.iterations)
}

func (s *stats) show(duration time.Duration) {
	dur := duration.Seconds()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Parameters", "Values"})
	table.AppendBulk([][]string{
		{"Concurrency", fmt.Sprintf("%d", *T)},
		{"Duration", fmt.Sprintf("%f.02s", dur)},
		{"Iterations", fmt.Sprintf("%d", s.iterations)},
		{"Successes", fmt.Sprintf("%d (%02.02f%%)", s.successes, s.asPercent(s.successes))},
		{"Errors", fmt.Sprintf("%d (%02.02f%%)", s.errors, s.asPercent(s.errors))},
		{"IPS", fmt.Sprintf("%.02f", float64(s.iterations)/float64(dur))},
		{"Latency (max)", fmt.Sprintf("%v", s.maxLatency)},
		{"Latency (p99)", fmt.Sprintf("%v", s.percentile(99))},
		{"Latency (p95)", fmt.Sprintf("%v", s.percentile(95))},
		{"Latency (p50)", fmt.Sprintf("%v", s.percentile(50))},
		{"Latency (min)", fmt.Sprintf("%v", s.minLatency)},
	})
	if s.endBal > 0 {
		table.AppendBulk([][]string{
			{"Balance (start)", fmt.Sprintf("%d", s.startBal)},
			{"Balance (end)", fmt.Sprintf("%d", s.endBal)},
			{"Balance (expected)", fmt.Sprintf("%d", s.startBal+(10*s.successes))},
		})
	}

	table.SetAlignment(tablewriter.ALIGN_LEFT)
	fmt.Println("Summary")
	table.Render()
	if s.errors > 0 {
		fmt.Printf("\nError Details\n")
		table = tablewriter.NewWriter(os.Stdout)
		for errmsg, count := range s.errorMap {
			table.SetHeader([]string{"Error", "Count", "%"})
			table.Append(
				[]string{errmsg, fmt.Sprintf("%d", count), fmt.Sprintf("%.2f", s.asPercent(count))},
			)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.Render()
	}
}