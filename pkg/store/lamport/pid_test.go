package lamport_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/store/lamport"
)

func TestPID(t *testing.T) {
	pid := lamport.PID(42)
	require.Equal(t, lamport.Scalar{42, 1}, pid.Next(nil), "nil pid didn't return expected next")
	require.Equal(t, lamport.Scalar{42, 1}, pid.Next(&lamport.Scalar{}), "zero pid didn't return expected next")
	require.Equal(t, lamport.Scalar{42, 19}, pid.Next(&lamport.Scalar{42, 18}), "same pid didn't return expected next")
}

func TestProcess(t *testing.T) {
	t.Run("PanicUnset", func(t *testing.T) {
		require.Panics(t, func() { lamport.ProcessID() }, "ProcessID didn't panic when unset")
		require.Panics(t, func() { lamport.Next(nil) }, "ProcessID didn't panic when unset")
	})

	t.Run("SetAndGet", func(t *testing.T) {
		lamport.SetProcessID(99)
		require.Equal(t, lamport.PID(99), lamport.ProcessID(), "ProcessID didn't return expected value after set")
	})

	t.Run("Next", func(t *testing.T) {
		lamport.SetProcessID(7)
		require.Equal(t, lamport.Scalar{7, 1}, lamport.Next(nil), "Next with nil didn't return expected value")
		require.Equal(t, lamport.Scalar{7, 19}, lamport.Next(&lamport.Scalar{7, 18}), "Next with same pid didn't return expected value")
	})

}
