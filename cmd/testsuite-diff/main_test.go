package main

import (
	"testing"
)

func TestDiff(t *testing.T) {

	testCases := map[string]struct {
		suiteA         map[string]bool
		suiteB         map[string]bool
		hasDiff        bool
		SuitesInTarget []string
	}{
		"Test in suiteA and suiteB": {
			suiteA:         map[string]bool{"hello": true},
			suiteB:         map[string]bool{"hello": true},
			hasDiff:        false,
			SuitesInTarget: []string{"hello"},
		},
		"Test in suiteA missing in suiteB": {
			suiteA:         map[string]bool{"hello": true},
			suiteB:         map[string]bool{"world": true},
			hasDiff:        true,
			SuitesInTarget: []string{},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			hasDiff, SuitesInTarget := diff(testCase.suiteA, "foo.yaml", testCase.suiteB, "bar.yaml", "+")
			if hasDiff != testCase.hasDiff {
				t.Errorf("Diff unexpected")
			}
			if len(SuitesInTarget) != len(testCase.SuitesInTarget) {
				t.Errorf("Diff unexpected")
			}
		})
	}

}
