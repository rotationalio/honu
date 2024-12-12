package lamport_test

import (
	"testing"

	"github.com/rotationalio/honu/pkg/store/lamport"
	"github.com/stretchr/testify/require"
)

func TestPID(t *testing.T) {
	pid := lamport.PID(42)
	require.Equal(t, &lamport.Scalar{42, 1}, pid.Next(nil), "nil pid didn't return expected next")
	require.Equal(t, &lamport.Scalar{42, 1}, pid.Next(&lamport.Scalar{}), "zero pid didn't return expected next")
	require.Equal(t, &lamport.Scalar{42, 19}, pid.Next(&lamport.Scalar{42, 18}), "same pid didn't return expected next")
}