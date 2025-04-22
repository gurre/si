package si

import (
	"fmt"
	"math"
)

// PrefixedFormatter extends the DefaultFormatter with support for SI prefixes
type PrefixedFormatter struct {
	DefaultFormatter
}

// NewPrefixedFormatter creates a new formatter that supports SI prefixes
func NewPrefixedFormatter() *PrefixedFormatter {
	return &PrefixedFormatter{
		DefaultFormatter: *NewDefaultFormatter(),
	}
}

// Format formats a node, applying SI prefixes to base units when appropriate
func (f *PrefixedFormatter) Format(node Node) (string, error) {
	// First get the standard format
	str, err := f.DefaultFormatter.Format(node)
	if err != nil {
		return "", err
	}

	// For simple values, no prefixing is needed
	if value, ok := extractSimpleValue(node); ok {
		prefix, scaled := computePrefix(value)
		if prefix != "" {
			switch n := node.(type) {
			case *IdentNode:
				return fmt.Sprintf("%g %s%s", scaled, prefix, n.Symbol), nil
			}
		}
	}

	return str, nil
}

// FormatUnitWithPrefix formats a unit with appropriate SI prefixes
func FormatUnitWithPrefix(u Unit) string {
	// If dimensionless, just return the value
	if u.Dimension == Dimensionless {
		return fmt.Sprintf("%g", u.Value)
	}

	// Handle special case for scaled base units
	for i, exp := range u.Dimension {
		if isBaseSIUnit(u.Dimension, i, exp) {
			symbols := []string{"m", "kg", "s", "A", "K", "mol", "cd"}
			prefix, scaled := computePrefix(u.Value)
			return fmt.Sprintf("%g %s%s", scaled, prefix, symbols[i])
		}
	}

	// Handle special cases for derived units
	if u.Dimension == Newton.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sN", scaled, prefix)
	} else if u.Dimension == Pascal.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sPa", scaled, prefix)
	} else if u.Dimension == Joule.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sJ", scaled, prefix)
	} else if u.Dimension == Watt.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sW", scaled, prefix)
	} else if u.Dimension == Hertz.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sHz", scaled, prefix)
	} else if u.Dimension == Volt.Dimension {
		prefix, scaled := computePrefix(u.Value)
		return fmt.Sprintf("%g %sV", scaled, prefix)
	}

	// Fall back to standard formatting
	return FormatUnit(u)
}

// extractSimpleValue attempts to extract a simple numeric value from an AST node
func extractSimpleValue(node Node) (float64, bool) {
	switch n := node.(type) {
	case *NumberNode:
		return n.Value, true
	case *IdentNode:
		// We can't determine a value from just an identifier
		return 0, false
	default:
		return 0, false
	}
}

// computePrefix computes the appropriate SI prefix for a value
func computePrefix(value float64) (string, float64) {
	absValue := math.Abs(value)

	switch {
	case absValue == 0:
		return "", 0
	case absValue >= 1e9:
		return "G", value / 1e9
	case absValue >= 1e6:
		return "M", value / 1e6
	case absValue >= 1e3:
		return "k", value / 1e3
	case absValue >= 1:
		return "", value
	case absValue >= 1e-3:
		return "m", value * 1e3
	case absValue >= 1e-6:
		return "μ", value * 1e6
	case absValue >= 1e-9:
		return "n", value * 1e9
	default:
		return "p", value * 1e12
	}
}

// isBaseSIUnit checks if a dimension represents a simple base SI unit
func isBaseSIUnit(dim Dimension, index int, exponent int) bool {
	if exponent != 1 {
		return false
	}

	// Check if all other dimensions are zero
	for i, e := range dim {
		if i != index && e != 0 {
			return false
		}
	}

	return true
}
