// sample of a lightmetric telegraf exec input plugin
//
// License: The MIT License (MIT)

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tesibelda/lightmetric/metric"
)

func main() {
	var ctags = make(map[string]string)
	var cfields = make(map[string]interface{})
	var mytag = "Shaun"

	flag.Parse()
	values := flag.Args()
	if len(values) > 0 {
		mytag = values[0]
	}

	rand.Seed(time.Now().UnixNano())

	ctags["sheep"] = mytag
	cfields["counter"] = rand.Intn(100) //nolint:gosec
	t := metric.TimeWithPrecision(time.Now(), time.Millisecond)
	m := metric.New("randplugin", ctags, cfields, t)

	fmt.Fprintln(os.Stdout, m.String(metric.InfluxLp))
}
