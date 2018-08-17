package si

import (
	"testing"
)

func TestUnit(t *testing.T) {
	unit := Length
	result := unit.String()
	expected := "m"
	if result != expected {
		t.Fatalf("Expected %s to equal %s", result, expected)
	}
}
