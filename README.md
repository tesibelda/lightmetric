[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/tesibelda/lightmetric/raw/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tesibelda/lightmetric)](https://goreportcard.com/report/github.com/tesibelda/lightmetric)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tesibelda/lightmetric?display_name=release)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tesibelda/lightmetric)](https://pkg.go.dev/github.com/tesibelda/lightmetric)

# lightmetric

lightmetric contains libraries that help creating plugins for [telegraf](https://github.com/influxdata/telegraf) monitoring agent. metric library includes Metric and metric Accumlator; Shim helps creating execd input plugins like telegraf's [shim](https://github.com/influxdata/telegraf/tree/master/plugins/common/shim). Most of the code has been borrowed from telegraf's repository and most of it has not been modified, but it has been reorganized so that it provides:
* much less library dependencies
* smaller final binary size
* ability to set precision to the execd input plugin shim
* ability to use the shim's [context](https://pkg.go.dev/context) in the plugin's Gather function

Hopefully as telegraf evolves this library will not be helpful in the future.


# Examples

## Example of exec input plugin using lightmetric

See the complete example at [samples/rand](https://github.com/tesibelda/lightmetric/tree/main/samples/rand) folder:

```go
	ctags["sheep"] = mytag
	cfields["counter"] = rand.Intn(100)
	t := metric.TimeWithPrecision(time.Now(), time.Millisecond)
	m := metric.New("randplugin", ctags, cfields, t)
	fmt.Fprintln(os.Stdout, m.String(metric.InfluxLp))
```

## Example of execd input plugin using Shim

See the complete example at [samples/counter](https://github.com/tesibelda/lightmetric/tree/main/samples/counter) folder:

```go
	p := NewCounter(mytag)
	p.Start()
	execd := shim.New("SheepCounter").WithPrecision(time.Second)
	err := execd.RunInput(p.Gather)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running lightmetric telegraf input shim: %w\n", err)
		os.Exit(1)
	}
	p.Stop()
```


# Telegraf configuration

## For an exec input plugin

Use influx data format.

```
[[inputs.exec]]
  commands = ["plugins/inputs/execd/yourplugin"]
  data_format = "influx"
```

Reference: [exec input](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/exec) 

## For execd input plugin shim

Execd input shim requires telegraf to be configured with: signal = "STDIN"

```toml
[[inputs.execd]]
  command = ["plugins/inputs/execd/yourplugin --config plugins/inputs/execd/yourplugin.conf"]
  signal = "STDIN"
```

Reference: [execd input](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/execd)


# License

[The MIT License (MIT)](https://github.com/tesibelda/vcstat/blob/master/LICENSE)
