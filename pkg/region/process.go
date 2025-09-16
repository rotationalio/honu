package region

import "sync"

var (
	processMu     sync.Mutex
	processRegion *Region
)

func SetProcessRegion(r Region) {
	processMu.Lock()
	processRegion = &r
	processMu.Unlock()
}

func ProcessRegion() Region {
	processMu.Lock()
	defer processMu.Unlock()
	if processRegion == nil {
		panic("process region not set")
	}
	return *processRegion
}
