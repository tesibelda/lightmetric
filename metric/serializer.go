// lightmetric helps with easy manipulation of simple telegraf metrics
//  without telegraf libraries dependencies
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)

package metric

import (
	"fmt"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
)

// BytesFormat represents supported decode formats
type BytesFormat uint8

const (
	Go BytesFormat = iota
	InfluxLp
)

// String returns a []byte representation of the metric in the given format
//
//	(https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md)
//
// Currently it only supports influx line protocol
func (m *Metric) Bytes(format BytesFormat) []byte {
	switch format {
	case Go:
		return m.bytesGo()
	case InfluxLp:
		return m.bytesInfluxLp()
	}
	return nil
}

// String returns a string representation of the metric in the given format
func (m *Metric) String(format BytesFormat) string {
	return string(m.Bytes(format))
}

// bytesGo returns a representation of the metric as commonly used with go
func (m *Metric) bytesGo() []byte {
	return []byte(fmt.Sprintf("%s %v %v %d", m.name, m.Tags(), m.Fields(), m.tm.UnixNano()))
}

// bytesInfluxLp returns a representation of the metric in influx line protocol
func (m *Metric) bytesInfluxLp() []byte {
	var (
		enc lineprotocol.Encoder
		val lineprotocol.Value
		ok  bool
	)

	// use lax encoding to avoid tags 'out of order' error
	enc.SetLax(true)
	enc.StartLine(m.name)
	for k, v := range m.Tags() {
		if len(k) > 0 && len(v) > 0 {
			enc.AddTag(k, v)
		}
	}
	for k, v := range m.Fields() {
		if val, ok = lineprotocol.NewValue(v); ok {
			enc.AddField(k, val)
		}
	}
	enc.EndLine(m.tm)

	return enc.Bytes()
}
