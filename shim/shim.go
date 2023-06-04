// shim is a telegraf input plugin shim using lightmetric
//
// License: The MIT License (MIT)

package shim

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tesibelda/lightmetric/metric"
)

type empty struct{}

type Shim struct {
	precision time.Duration
	outFormat metric.BytesFormat

	// streams
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	// outgoing metric channel
	metricCh       chan metric.Metric
	gatherPromptCh chan empty
}

// New creates a new shim interface
func New() *Shim {
	return &Shim{
		precision: time.Nanosecond,
		outFormat: metric.InfluxLp,
		stdin:     os.Stdin,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		metricCh:  make(chan metric.Metric, 1),
	}
}

// WithPrecision sets time precision to the shim's metric accumulator
func (s *Shim) WithPrecision(precision time.Duration) *Shim {
	s.precision = precision
	return s
}

// WithOutputFormat sets the format for the serialization of the output metrics
func (s *Shim) WithOutputFormat(format metric.BytesFormat) *Shim {
	s.outFormat = format
	return s
}

// pushCollectMetricsRequest pushes a non-blocking (nil) message to the
// gatherPromptCh channel to trigger metric collection.
// The channel is defined with a buffer of 1, so while it's full, subsequent
// requests are discarded.
func (s *Shim) pushCollectMetricsRequest() {
	// push a message out to each channel to collect metrics. don't block.
	select {
	case s.gatherPromptCh <- empty{}:
	default:
	}
}

func (s *Shim) writeProcessedMetrics() error {
	var mby []byte
	var err error

	for { //nolint:gosimple // for-select used on purpose
		select {
		case m, open := <-s.metricCh:
			if !open {
				return nil
			}
			// Serialize metric
			if mby = m.Bytes(s.outFormat); mby == nil {
				return fmt.Errorf("failed to serialize metric %s", m.Name())
			}
			// Write this to stdout
			if _, err = fmt.Fprint(s.stdout, string(mby)); err != nil {
				return fmt.Errorf("failed to write metric: %w", err)
			}
		}
	}
}

func (s *Shim) watchForShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit // user-triggered quit
		// cancel, but keep looping until the metric channel closes.
		cancel()
	}()
}

func hasQuit(ctx context.Context) bool {
	return ctx.Err() != nil
}
