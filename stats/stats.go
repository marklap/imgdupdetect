package stats

import (
	"fmt"
	"sync"
	"time"
)

// ScanStats holds statistics of the images processed
type ScanStats struct {
	sync.RWMutex
	Start            time.Time
	End              time.Time
	ImagesFound      int
	FingerPrintCount int
	DuplicatesFound  int
}

// NewScanStats creates a new Statistics object
func NewScanStats() *ScanStats {
	return &ScanStats{Start: time.Now()}
}

// Complete sets the end time of the scan
func (s *ScanStats) Complete() {
	s.End = time.Now()
}

// Duration returns the time it took to run the scan
func (s *ScanStats) Duration() time.Duration {
	if s.End.IsZero() {
		return time.Duration(0)
	}
	return s.End.Sub(s.Start)
}

// Rate returns the average time it takes to find, fingerprint and companre an image
func (s *ScanStats) Rate() time.Duration {
	d := s.Duration()
	if d == 0 {
		return d
	}
	return time.Duration(uint64(d) / uint64(s.FingerPrintCount))
}

// String returns a printable string of stats
func (s *ScanStats) String() string {
	return fmt.Sprintf("scanning took %s (avg %s/image); found %d images; fingerprinted %d images; %d duplicates found",
		s.Duration(), s.Rate(), s.ImagesFound, s.FingerPrintCount, s.DuplicatesFound)
}

// ImagesFoundIncr increments the image count counter
func (s *ScanStats) ImagesFoundIncr() {
	s.Lock()
	defer s.Unlock()
	s.ImagesFound++
}

// FingerPrintCountIncr increments the image count counter
func (s *ScanStats) FingerPrintCountIncr() {
	s.Lock()
	defer s.Unlock()
	s.FingerPrintCount++
}
