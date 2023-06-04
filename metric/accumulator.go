// lightmetric helps with easy manipulation of telegraf like metrics
//  without telegraf libraries dependencies
//
// License: The MIT License (MIT)

package metric

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Accumulator struct {
	errfile   io.Writer
	metrics   chan<- Metric
	precision time.Duration
}

// NewAccumulator returns a new Accumulator instance
func NewAccumulator(
	errfile io.Writer,
	metrics chan<- Metric,
) *Accumulator {
	acc := Accumulator{
		errfile:   errfile,
		metrics:   metrics,
		precision: time.Nanosecond,
	}
	return &acc
}

func (ac *Accumulator) AddFields(
	measurement string,
	fields map[string]interface{},
	tags map[string]string,
	t time.Time,
) {
	m := New(measurement, tags, fields, t)
	m.SetTime(m.Time().Round(ac.precision))
	ac.metrics <- m
}

func (ac *Accumulator) AddMetric(m Metric) {
	m.SetTime(m.Time().Round(ac.precision))
	ac.metrics <- m
}

// AddError passes a runtime error to the Accumulator.
// The error will be tagged with the plugin name and written to the log.
func (ac *Accumulator) AddError(err error) {
	if err != nil {
		_, werr := fmt.Fprintf(ac.errfile, "Error in plugin: %s\n", err)
		if werr != nil {
			fmt.Fprintln(os.Stderr, "Error in plugin: "+err.Error())
			fmt.Fprintln(os.Stderr, "Error logging previous error: "+werr.Error())
		}
	}
}

func (ac *Accumulator) SetPrecision(precision time.Duration) {
	ac.precision = precision
}
