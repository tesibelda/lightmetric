[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/tesibelda/lightmetric/raw/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tesibelda/lightmetric)](https://goreportcard.com/report/github.com/tesibelda/lightmetric)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tesibelda/lightmetric?display_name=release)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tesibelda/lightmetric)](https://pkg.go.dev/github.com/tesibelda/lightmetric)

# lightmetric

lightmetric contains libraries that help creating plugins for [telegraf](https://github.com/influxdata/telegraf) monitoring agent. metric library includes Metric and metric Accumulator that you may use for exec input plugins; Shim helps creating execd input plugins like telegraf's [shim](https://github.com/influxdata/telegraf/tree/master/plugins/common/shim). Most of the code has been borrowed from telegraf's repository and most of it has not been modified, but it has been reorganized so that it provides:
* much less library dependencies
* smaller final binary size
* ability to set precision to the execd input plugin shim
* ability to use the shim's [context](https://pkg.go.dev/context) in the plugin's Gather function

Hopefully as telegraf evolves this library will not be helpful in the future.


# Examples

## Example of an exec input plugin using Metric

See the complete example at [examples/rand](https://github.com/tesibelda/lightmetric/tree/main/examples/rand) folder:

```go
	ctags["sheep"] = mytag
	cfields["counter"] = rand.Intn(100)
	t := metric.TimeWithPrecision(time.Now(), time.Millisecond)
	m := metric.New("randplugin", ctags, cfields, t)
	fmt.Fprint(os.Stdout, m.String(metric.InfluxLp))
```

### Example output
```plain
randplugin,sheep=Shaun counter=66i 1686482784943000000
```

## Example of an execd input plugin using Shim

See the complete example at [examples/counter](https://github.com/tesibelda/lightmetric/tree/main/examples/counter) folder:

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

### Example output
```plain
counterplugin,sheep=Shaun counter=1i 1686483004000000000
counterplugin,sheep=Shaun counter=2i 1686483014000000000
counterplugin,sheep=Shaun counter=3i 1686483024000000000
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
