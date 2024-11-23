// shim is a telegraf input plugin shim using lightmetric
//
// License: The MIT License (MIT)

package shim //nolint: testpackage

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/tesibelda/lightmetric/metric"
)

func TestShimSetsUpLogger(t *testing.T) {
	stderrReader, stderrWriter := io.Pipe()
	stdinReader, stdinWriter := io.Pipe()

	runErroringInputPlugin(t, stdinReader, nil, stderrWriter)

	_, err := stdinWriter.Write([]byte("\n"))
	require.NoError(t, err)

	// <-metricProcessed

	r := bufio.NewReader(stderrReader)
	out, err := r.ReadString('\n')
	require.NoError(t, err)
	require.Contains(t, out, "Error in plugin ShimTest: intentional")

	err = stdinWriter.Close()
	require.NoError(t, err)
}

func runErroringInputPlugin(
	t *testing.T,
	stdin io.Reader,
	stdout, stderr io.Writer,
) (chan bool, chan bool) {
	metricProcessed := make(chan bool, 1)
	exited := make(chan bool, 1)
	inp := &erroringInput{}

	shim := New("ShimTest").WithPrecision(time.Millisecond)
	if stdin != nil {
		shim.stdin = stdin
	}
	if stdout != nil {
		shim.stdout = stdout
	}
	if stderr != nil {
		shim.stderr = stderr
		log.SetOutput(stderr)
	}
	go func() {
		err := shim.RunInput(inp.Gather)
		require.NoError(t, err)
		exited <- true
	}()
	return metricProcessed, exited
}

type erroringInput struct {
}

func (i *erroringInput) Gather(_ context.Context, acc *metric.Accumulator) error {
	acc.AddError(errors.New("intentional"))
	return nil
}
