package region_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/honu/pkg/region"
)

func TestProcess(t *testing.T) {
	t.Run("Unset", func(t *testing.T) {
		require.Panics(t, func() { region.ProcessRegion() }, "accessing unset process region should panic")
	})

	t.Run("Set", func(t *testing.T) {
		region.SetProcessRegion(region.DEVELOPMENT)
		r := region.ProcessRegion()
		require.Equal(t, region.DEVELOPMENT, r, "process region did not match set value")
	})

	t.Run("Concurrency", func(t *testing.T) {
		region.SetProcessRegion(region.GCP_US_EAST_1B)

		var wg sync.WaitGroup
		wg.Go(func() {
			r := region.ProcessRegion()
			require.Equal(t, region.GCP_US_EAST_1B, r, "process region did not match set value in goroutine")
		})
		wg.Wait()

	})
}
