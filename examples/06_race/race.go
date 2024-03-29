package race

import "sync"

var (
	// Counter represents a globally shared resource guarded by a mutex.
	// It is not safe for concurrent access.
	// Use [Inc] to modify it.
	Counter = 1
	mux     = sync.Mutex{}
)

// Inc is a concurrent safe way of incrementing the [Counter]
func Inc() {
	mux.Lock()
	defer mux.Unlock()
	Counter++
}
