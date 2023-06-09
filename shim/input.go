// shim is a telegraf input plugin shim using lightmetric
//
// License: The MIT License (MIT)

package shim

import (
	"bufio"
	"context"
	"fmt"
	"sync"

	"github.com/tesibelda/lightmetric/metric"
)

type GatherFunc func(context.Context, metric.Accumulator) error

func (s *Shim) RunInput(g GatherFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.watchForShutdown(cancel)

	acc := metric.NewAccumulator(s.shimname, s.metricCh).WithErrorWriter(s.stderr)
	acc.SetPrecision(s.precision)

	s.gatherPromptCh = make(chan empty, 1)
	go func() {
		s.startGathering(ctx, g, *acc)
		// closing the metric channel gracefully stops writing to stdout
		close(s.metricCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.writeProcessedMetrics()
		if err != nil {
			fmt.Fprintln(s.stderr, err.Error())
		}
		wg.Done()
	}()

	go func() {
		scanner := bufio.NewScanner(s.stdin)
		for scanner.Scan() {
			// push a non-blocking message to trigger metric collection.
			s.pushCollectMetricsRequest()
		}

		cancel() // cancel gracefully stops gathering
	}()

	wg.Wait() // wait for writing to stdout to finish
	return nil
}

func (s *Shim) startGathering(
	ctx context.Context,
	g GatherFunc,
	acc metric.Accumulator,
) {
	var err error

	for {
		// give priority to stopping.
		if hasQuit(ctx) {
			return
		}
		// see what's up
		select {
		case <-ctx.Done():
			return
		case <-s.gatherPromptCh:
			if err = g(ctx, acc); err != nil {
				fmt.Fprintf(s.stderr, "failed to gather metrics: %s\n", err)
			}
		}
	}
}
