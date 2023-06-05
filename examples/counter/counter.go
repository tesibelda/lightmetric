// sample of a lightmetric telegraf input plugin shim
//
// License: The MIT License (MIT)

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tesibelda/lightmetric/metric"
	"github.com/tesibelda/lightmetric/shim"
)

func main() {
	var mytag = "Shaun"

	flag.Parse()
	values := flag.Args()
	if len(values) > 0 {
		mytag = values[0]
	}

	p := NewCounter(mytag)
	_ = p.Start()
	execd := shim.New().WithPrecision(time.Second)
	err := execd.RunInput(p.Gather)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running lightmetric telegraf input shim: %s\n", err)
		os.Exit(1)
	}
	_ = p.Stop()
}

type Counter struct {
	sheep   string
	counter int
}

func NewCounter(name string) *Counter {
	return &Counter{
		sheep: name,
	}
}

// Start starts whatever logic needed before for gathering metrics
func (c *Counter) Start() error {
	return nil
}

// Stop stops whatever logic needed before plugin exit
func (c *Counter) Stop() error {
	return nil
}

// Gather get the metrics values and adds them to the given acummulator
func (c *Counter) Gather(ctx context.Context, acc metric.Accumulator) error {
	var ctags = make(map[string]string)
	var cfields = make(map[string]interface{})

	ctags["sheep"] = c.sheep
	c.counter++
	cfields["counter"] = c.counter
	acc.AddFields("counterplugin", cfields, ctags, time.Now())
	return nil
}
