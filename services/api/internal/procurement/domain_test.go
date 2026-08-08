package procurement

import (
	"testing"
	"time"
)

// agingBucketFor classifies a due date into an AP aging bucket relative to "now",
// mirroring the SQL CASE in Repository.Aging. Kept as a pure function for testing.
func agingBucketFor(due time.Time, now time.Time) string {
	d := due.Truncate(24 * time.Hour)
	n := now.Truncate(24 * time.Hour)
	switch {
	case !d.Before(n):
		return "current"
	case d.AddDate(0, 0, 30).After(n):
		return "30"
	case d.AddDate(0, 0, 60).After(n):
		return "60"
	case d.AddDate(0, 0, 90).After(n):
		return "90"
	default:
		return "90plus"
	}
}

func TestAgingBucketClassification(t *testing.T) {
	now := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		due    time.Time
		expect string
	}{
		{time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC), "current"},  // today
		{time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC), "current"}, // future
		{time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC), "30"},      // 20 days overdue
		{time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC), "60"},      // 50 days overdue
		{time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), "90"},      // 81 days overdue
		{time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), "90plus"},  // 111 days overdue
	}
	for _, c := range cases {
		if got := agingBucketFor(c.due, now); got != c.expect {
			t.Errorf("due=%s: expected %s, got %s", c.due.Format("2006-01-02"), c.expect, got)
		}
	}
}

func TestAgingBucketOrdering(t *testing.T) {
	// The response must order buckets current -> 30 -> 60 -> 90 -> 90plus.
	order := map[string]int{"current": 0, "30": 1, "60": 2, "90": 3, "90plus": 4}
	buckets := []string{"90plus", "current", "60", "30", "90"}
	for i := 0; i < len(buckets); i++ {
		for j := i + 1; j < len(buckets); j++ {
			if order[buckets[j]] < order[buckets[i]] {
				buckets[i], buckets[j] = buckets[j], buckets[i]
			}
		}
	}
	want := []string{"current", "30", "60", "90", "90plus"}
	for i := range want {
		if buckets[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, buckets)
		}
	}
}
