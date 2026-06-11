package service

import (
	"testing"
	"time"
)

// TestCalculateAge checks dynamic age calculation logic against edge cases.
func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      string
		ref      string
		expected int
	}{
		{
			name:     "Birthday has passed this year",
			dob:      "1990-05-10",
			ref:      "2025-06-15",
			expected: 35,
		},
		{
			name:     "Birthday is today",
			dob:      "1990-05-10",
			ref:      "2025-05-10",
			expected: 35,
		},
		{
			name:     "Birthday is tomorrow (has not passed yet)",
			dob:      "1990-05-10",
			ref:      "2025-05-09",
			expected: 34,
		},
		{
			name:     "Birthday month is later in the year",
			dob:      "1990-12-01",
			ref:      "2025-05-10",
			expected: 34,
		},
		{
			name:     "Leap year birthday - non-leap reference (Feb 28)",
			dob:      "2004-02-29",
			ref:      "2005-02-28",
			expected: 0,
		},
		{
			name:     "Leap year birthday - non-leap reference (March 1)",
			dob:      "2004-02-29",
			ref:      "2005-03-01",
			expected: 1,
		},
		{
			name:     "Leap year birthday - leap reference (Feb 29)",
			dob:      "2004-02-29",
			ref:      "2008-02-29",
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dobTime, err := time.Parse("2006-01-02", tt.dob)
			if err != nil {
				t.Fatalf("Failed to parse DOB test case: %v", err)
			}
			refTime, err := time.Parse("2006-01-02", tt.ref)
			if err != nil {
				t.Fatalf("Failed to parse reference time test case: %v", err)
			}

			got := CalculateAge(dobTime, refTime)
			if got != tt.expected {
				t.Errorf("CalculateAge(%s, %s) = %d; want %d", tt.dob, tt.ref, got, tt.expected)
			}
		})
	}
}
