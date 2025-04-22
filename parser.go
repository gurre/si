package si

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// sortedPrefixes ensures longest prefixes match first when parsing.
var sortedPrefixes = func() []string {
	keys := make([]string, 0, len(Prefixes))
	for k := range Prefixes {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}()

// formatDimension renders a dimension into a string representation.
func formatDimension(d Dimension) string {
	symbols := []string{"m", "kg", "s", "A", "K", "mol", "cd"}
	var numerator []string
	var denominator []string

	for i, exp := range d {
		if exp == 0 {
			continue
		}

		if exp > 0 {
			if exp == 1 {
				numerator = append(numerator, symbols[i])
			} else {
				numerator = append(numerator, fmt.Sprintf("%s^%d", symbols[i], exp))
			}
		} else {
			if exp == -1 {
				denominator = append(denominator, symbols[i])
			} else {
				denominator = append(denominator, fmt.Sprintf("%s^%d", symbols[i], -exp))
			}
		}
	}

	// Handle special cases
	if len(numerator) == 0 && len(denominator) == 0 {
		return "1"
	}

	if len(numerator) == 0 {
		numerator = append(numerator, "1")
	}

	// Build the numerator part
	numStr := strings.Join(numerator, "·")

	// If there are denominators, append them with "/"
	if len(denominator) > 0 {
		denomStr := strings.Join(denominator, "·")
		return fmt.Sprintf("%s/%s", numStr, denomStr)
	}

	return numStr
}

// ParseUnit parses a unit string like "km/h" into a Unit.
// It handles basic units, prefixed units, compound units, and symbolic units.
func ParseUnit(input string) (Unit, error) {
	// Handle special cases
	if input == "" || input == "1" {
		return One, nil
	}

	if unit, ok := SymbolicUnits[input]; ok {
		return unit, nil
	}

	// Split by division
	parts := strings.Split(input, "/")

	// Parse numerator
	unit, err := parseExpression(parts[0])
	if err != nil {
		return One, fmt.Errorf("invalid numerator: %w", err)
	}

	// Parse and apply denominators
	for i := 1; i < len(parts); i++ {
		denomPart := parts[i]
		// Check for exponents
		expParts := strings.Split(denomPart, "^")
		denom, err := parseExpression(expParts[0])
		if err != nil {
			return One, fmt.Errorf("invalid denominator: %w", err)
		}

		// Apply exponent if present
		if len(expParts) > 1 {
			exp, err := strconv.Atoi(expParts[1])
			if err != nil {
				return One, fmt.Errorf("invalid exponent: %w", err)
			}
			denom = denom.Pow(exp)
		}

		unit = unit.Div(denom)
	}

	return unit, nil
}

// parseExpression handles unit expressions that may contain products (· or *)
func parseExpression(expr string) (Unit, error) {
	var parts []string

	// Handle special cases for complex units with fractions
	complexUnits := map[string]Unit{
		"J/(kg·K)": {1, Dimension{2, -1, -2, 0, -1, 0, 0}}, // Specific heat capacity
		"W/(m·K)":  {1, Dimension{1, 1, -3, 0, -1, 0, 0}},  // Thermal conductivity
		"N·s/kg":   {1, Dimension{1, 0, -1, 0, 0, 0, 0}},   // Specific impulse
	}

	// Try matching complex units directly
	if unit, found := complexUnits[expr]; found {
		return unit, nil
	}

	// Check for other complex cases with fractions or parentheses
	if strings.Contains(expr, "/") && strings.Contains(expr, "(") {
		// For complex expressions like "W/(m·K)" that aren't in our map
		parts := strings.Split(expr, "/")
		if len(parts) == 2 {
			numerator, err := parseExpression(parts[0])
			if err != nil {
				return One, err
			}

			// Handle parentheses in the denominator
			denom := parts[1]
			if strings.HasPrefix(denom, "(") && strings.HasSuffix(denom, ")") {
				denom = denom[1 : len(denom)-1]
			}

			denominator, err := parseExpression(denom)
			if err != nil {
				return One, err
			}

			return numerator.Div(denominator), nil
		}
	}

	// Check for product notation and split accordingly
	switch {
	case strings.Contains(expr, "·"):
		parts = strings.Split(expr, "·")
	case strings.Contains(expr, "*"):
		parts = strings.Split(expr, "*")
	default:
		// Single unit, no product
		return parseBasicUnit(expr)
	}

	// Apply multiplication for all parts
	result := One
	for _, part := range parts {
		unit, err := parseBasicUnit(part)
		if err != nil {
			return One, err
		}
		result = result.Mul(unit)
	}
	return result, nil
}

// parseBasicUnit handles simple units with optional prefixes
func parseBasicUnit(unit string) (Unit, error) {
	if unit == "" {
		return One, fmt.Errorf("empty unit")
	}

	// Predefined units
	baseUnits := map[string]Unit{
		"m":   Meter,
		"kg":  Kilogram,
		"s":   Second,
		"A":   Ampere,
		"K":   Kelvin,
		"mol": Mole,
		"cd":  Candela,
	}

	derivedUnits := map[string]Unit{
		"Hz":  Hertz,
		"N":   Newton,
		"Pa":  Pascal,
		"J":   Joule,
		"W":   Watt,
		"C":   Coulomb,
		"V":   Volt,
		"h":   {3600, TimeDim},  // hour = 3600 seconds
		"min": {60, TimeDim},    // minute = 60 seconds
		"d":   {86400, TimeDim}, // day = 24 hours = 86400 seconds
		"B":   {1, Dimensionless},
		"iB":  {1, Dimensionless}, // For MiB, GiB, etc.
	}

	// Try direct match first
	if u, ok := baseUnits[unit]; ok {
		return u, nil
	}
	if u, ok := derivedUnits[unit]; ok {
		return u, nil
	}

	// Handle binary prefixes (KiB, MiB, etc.)
	if strings.HasSuffix(unit, "iB") {
		for _, prefix := range []string{"K", "M", "G", "T", "P", "E"} {
			if strings.HasPrefix(unit, prefix) && unit == prefix+"iB" {
				return Unit{Prefixes[prefix+"i"], Dimensionless}, nil
			}
		}
	}

	// Try with prefixes
	for _, prefix := range sortedPrefixes {
		if prefix == "" {
			continue
		}

		if !strings.HasPrefix(unit, prefix) {
			continue
		}

		suffix := unit[len(prefix):]

		// Try base units with this prefix
		if u, ok := baseUnits[suffix]; ok {
			return Unit{Prefixes[prefix], u.Dimension}, nil
		}

		// Try derived units with this prefix
		if u, ok := derivedUnits[suffix]; ok {
			return Unit{Prefixes[prefix] * u.Value, u.Dimension}, nil
		}
	}

	return One, fmt.Errorf("unrecognized unit: %s", unit)
}

// Parse splits and parses a full unit expression like "100 km/h".
// This extracts the numeric value and unit component from a string.
func Parse(input string) (Unit, error) {
	fields := strings.Fields(input)

	// Handle case with only a number (dimensionless unit)
	if len(fields) == 1 {
		val, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return One, fmt.Errorf("invalid numeric value: %w", err)
		}
		return Scalar(val), nil
	}

	if len(fields) < 2 {
		return One, fmt.Errorf("invalid unit expression: %s", input)
	}

	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return One, fmt.Errorf("invalid numeric value: %w", err)
	}

	u, err := ParseUnit(strings.Join(fields[1:], ""))
	if err != nil {
		return One, err
	}

	u.Value *= val
	return u, nil
}

// MustParse works like Parse but panics on error.
// Use this function only when you know the input is valid.
func MustParse(input string) Unit {
	u, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return u
}

// New creates a unit with a given value and unit string (e.g. "kg").
func New(value float64, symbol string) Unit {
	// Special case for grams
	if symbol == "g" {
		return Unit{value / 1000, Mass}
	}

	// First try to parse as a direct unit
	u, err := ParseUnit(symbol)
	if err != nil {
		// Return dimensionless unit as fallback
		return Scalar(value)
	}

	u.Value = value
	return u
}
