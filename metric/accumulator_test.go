// lightmetric helps with easy manipulation of telegraf like metrics
//  without telegraf libraries dependencies
//
// License: The MIT License (MIT)

package metric //nolint: testpackage

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddFields(t *testing.T) {
	metrics := make(chan Metric, 10)
	defer close(metrics)
	a := NewAccumulator("TestPlugin", metrics)

	tags := map[string]string{"foo": "bar"}
	fields := map[string]interface{}{
		"usage": float64(99),
	}
	now := time.Now()
	a.AddFields("acctest", fields, tags, now)

	testm := <-metrics

	require.Equal(t, "acctest", testm.Name())
	actual, ok := testm.GetField("usage")

	require.True(t, ok)
	require.Equal(t, float64(99), actual)

	actual, ok = testm.GetTag("foo")
	require.True(t, ok)
	require.Equal(t, "bar", actual)

	tm := testm.Time()
	// okay if monotonic clock differs
	require.True(t, now.Equal(tm))
}

func TestAccAddError(t *testing.T) {
	errBuf := bytes.NewBuffer(nil)

	metrics := make(chan Metric, 10)
	defer close(metrics)
	a := NewAccumulator("TestPlugin", metrics).WithErrorWriter(errBuf)

	a.AddError(errors.New("foo"))
	a.AddError(errors.New("bar"))
	a.AddError(errors.New("baz"))

	errs := bytes.Split(errBuf.Bytes(), []byte{'\n'})
	require.Len(t, errs, 4) // 4 because of trailing newline
	assert.Contains(t, string(errs[0]), "TestPlugin")
	assert.Contains(t, string(errs[0]), "foo")
	assert.Contains(t, string(errs[1]), "TestPlugin")
	assert.Contains(t, string(errs[1]), "bar")
	assert.Contains(t, string(errs[2]), "TestPlugin")
	assert.Contains(t, string(errs[2]), "baz")
}

func TestSetPrecision(t *testing.T) {
	tests := []struct {
		name      string
		unset     bool
		precision time.Duration
		timestamp time.Time
		expected  time.Time
	}{
		{
			name:      "default precision is nanosecond",
			unset:     true,
			timestamp: time.Date(2006, time.February, 10, 12, 0, 0, 82912748, time.UTC),
			expected:  time.Date(2006, time.February, 10, 12, 0, 0, 82912748, time.UTC),
		},
		{
			name:      "second interval",
			precision: time.Second,
			timestamp: time.Date(2006, time.February, 10, 12, 0, 0, 82912748, time.UTC),
			expected:  time.Date(2006, time.February, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "microsecond interval",
			precision: time.Microsecond,
			timestamp: time.Date(2006, time.February, 10, 12, 0, 0, 82912748, time.UTC),
			expected:  time.Date(2006, time.February, 10, 12, 0, 0, 82913000, time.UTC),
		},
		{
			name:      "2 second precision",
			precision: 2 * time.Second,
			timestamp: time.Date(2006, time.February, 10, 12, 0, 2, 4, time.UTC),
			expected:  time.Date(2006, time.February, 10, 12, 0, 2, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make(chan Metric, 10)

			a := NewAccumulator("TestPlugin", metrics)
			if !tt.unset {
				a.SetPrecision(tt.precision)
			}

			a.AddFields("acctest",
				map[string]interface{}{"value": float64(101)},
				map[string]string{},
				tt.timestamp,
			)

			testm := <-metrics
			require.Equal(t, tt.expected, testm.Time())

			close(metrics)
		})
	}
}
