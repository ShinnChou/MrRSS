package feed

import (
	"log"
	"time"
)

// Progress tracks the progress of feed fetching operations
type Progress struct {
	Total     int  `json:"total"`
	Current   int  `json:"current"`
	IsRunning bool `json:"is_running"`
}

// GetProgress returns the current progress of the feed fetching operation
func (f *Fetcher) GetProgress() Progress {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.progress
}

// waitForProgressComplete waits for any running operation to complete with a timeout.
// Returns true if the wait was successful, false if timeout occurred.
func (f *Fetcher) waitForProgressComplete(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for f.GetProgress().IsRunning {
		if time.Now().After(deadline) {
			log.Println("Timeout waiting for previous operation to complete")
			return false
		}
		time.Sleep(100 * time.Millisecond)
	}
	return true
}

// setProgress sets the progress state (thread-safe)
func (f *Fetcher) setProgress(total, current int, isRunning bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.progress.Total = total
	f.progress.Current = current
	f.progress.IsRunning = isRunning
}

// incrementProgress increments the current progress counter (thread-safe)
func (f *Fetcher) incrementProgress() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.progress.Current++
}

// startProgress initializes progress tracking for a new operation
func (f *Fetcher) startProgress(total int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.progress.IsRunning {
		return false
	}

	f.progress.IsRunning = true
	f.progress.Total = total
	f.progress.Current = 0
	return true
}

// finishProgress marks the operation as complete
func (f *Fetcher) finishProgress() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.progress.IsRunning = false
}
