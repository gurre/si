package si

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// assertUnitAlmostEqual checks that two units are almost equal (value within epsilon)
func assertUnitAlmostEqual(t *testing.T, got, want Unit, name string) {
	t.Helper()

	// Check dimension first
	if got.Dimension != want.Dimension {
		t.Errorf("%s dimension = %v, want %v", name, got.Dimension, want.Dimension)
		return
	}

	// Check value with epsilon
	epsilon := 1e-10
	absVal := math.Abs(got.Value - want.Value)
	if absVal > epsilon {
		t.Errorf("%s value = %v, want %v (± %v)", name, got.Value, want.Value, epsilon)
	}
}

// TestTokenizer verifies tokenization
func TestTokenizer(t *testing.T) {
	tests := []struct {
		input string
		want  []TokenKind
	}{
		{"m", []TokenKind{Identifier, EOF}},
		{"m/s", []TokenKind{Identifier, Divide, Identifier, EOF}},
		{"kg*m/s^2", []TokenKind{Identifier, Multiply, Identifier, Divide, Identifier, Power, Number, EOF}},
		{"(kg*m)/(s^2)", []TokenKind{LParen, Identifier, Multiply, Identifier, RParen, Divide, LParen, Identifier, Power, Number, RParen, EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Replace middle dot with asterisk to match the implementation
			normalizedInput := strings.Replace(tt.input, "·", "*", -1)

			tokens, err := tokenize(normalizedInput)
			if err != nil {
				t.Fatalf("tokenize(%q) error: %v", normalizedInput, err)
			}

			if len(tokens) != len(tt.want) {
				t.Fatalf("tokenize(%q) got %d tokens, want %d", normalizedInput, len(tokens), len(tt.want))
			}

			for i, token := range tokens {
				if token.Kind != tt.want[i] {
					t.Errorf("tokenize(%q) token[%d] = %v, want %v", normalizedInput, i, token.Kind, tt.want[i])
				}
			}
		})
	}
}

// TestParseUnitAST verifies AST construction
func TestParseUnitAST(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"m", true},
		{"kg", true},
		{"s", true},
		{"m/s", true},
		{"kg*m/s^2", true},
		{"(kg·m)/(s^2)", true},
		{"((kg*m)/s)^2", true},
		{"W/(m^2*K^4)", true},
		{"km/h", true},
		{"g/(cm^2*s)", true},
		{"(N*s)/m", true},
		{"(kg*m)/(s^2*K)", true},
		{"m/s/", false},  // Missing denominator
		{"kg*/s", false}, // Invalid operator usage
		{"m^2^3", false}, // Double exponentiation
		{"(kg", false},   // Unclosed parenthesis
		{"kg)", false},   // Unmatched closing parenthesis
		{"kg^x", false},  // Non-numeric exponent
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ast, err := ParseUnitAST(tt.input)
			if tt.valid {
				if err != nil {
					t.Errorf("ParseUnitAST(%q) error: %v", tt.input, err)
				} else if ast == nil {
					t.Errorf("ParseUnitAST(%q) returned nil AST without error", tt.input)
				}
			} else {
				if err == nil {
					t.Errorf("ParseUnitAST(%q) expected error, got nil", tt.input)
				}
			}
		})
	}
}

// TestEvalSimpleUnits verifies evaluation of simple units
func TestEvalSimpleUnits(t *testing.T) {
	ctx := NewStandardContext()

	tests := []struct {
		input string
		want  Unit
	}{
		{"m", Unit{1, Dimension{1, 0, 0, 0, 0, 0, 0}}},
		{"kg", Unit{1, Dimension{0, 1, 0, 0, 0, 0, 0}}},
		{"s", Unit{1, Dimension{0, 0, 1, 0, 0, 0, 0}}},
		{"g", Unit{0.001, Dimension{0, 1, 0, 0, 0, 0, 0}}},
		{"km", Unit{1000, Dimension{1, 0, 0, 0, 0, 0, 0}}},
		{"ms", Unit{0.001, Dimension{0, 0, 1, 0, 0, 0, 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ast, err := ParseUnitAST(tt.input)
			if err != nil {
				t.Fatalf("ParseUnitAST(%q) error: %v", tt.input, err)
			}

			got, err := EvalAST(ast, ctx)
			if err != nil {
				t.Fatalf("EvalAST(%q) error: %v", tt.input, err)
			}

			assertUnitAlmostEqual(t, got, tt.want, fmt.Sprintf("EvalAST(%q)", tt.input))
		})
	}
}

// TestEvalCompoundUnits verifies evaluation of compound units
func TestEvalCompoundUnits(t *testing.T) {
	ctx := NewStandardContext()

	tests := []struct {
		input string
		want  Unit
	}{
		// m/s - velocity
		{"m/s", Unit{1, Dimension{1, 0, -1, 0, 0, 0, 0}}},

		// km/h - velocity
		{"km/h", Unit{1000.0 / 3600.0, Dimension{1, 0, -1, 0, 0, 0, 0}}},

		// kg*m/s^2 - force (newton)
		{"kg*m/s^2", Unit{1, Dimension{1, 1, -2, 0, 0, 0, 0}}},

		// N - force
		{"N", Unit{1, Dimension{1, 1, -2, 0, 0, 0, 0}}},

		// J - energy
		{"J", Unit{1, Dimension{2, 1, -2, 0, 0, 0, 0}}},

		// W - power
		{"W", Unit{1, Dimension{2, 1, -3, 0, 0, 0, 0}}},

		// Pa - pressure
		{"Pa", Unit{1, Dimension{-1, 1, -2, 0, 0, 0, 0}}},

		// Complex: (kg*m)/(s^2)
		{"(kg*m)/(s^2)", Unit{1, Dimension{1, 1, -2, 0, 0, 0, 0}}},

		// Complex nested: ((kg*m)/s)^2
		{"((kg*m)/s)^2", Unit{1, Dimension{2, 2, -2, 0, 0, 0, 0}}},

		// Complex: W/(m^2*K^4) - Stefan-Boltzmann constant units
		{"W/(m^2*K^4)", Unit{1, Dimension{0, 1, -3, 0, -4, 0, 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ast, err := ParseUnitAST(tt.input)
			if err != nil {
				t.Fatalf("ParseUnitAST(%q) error: %v", tt.input, err)
			}

			got, err := EvalAST(ast, ctx)
			if err != nil {
				t.Fatalf("EvalAST(%q) error: %v", tt.input, err)
			}

			assertUnitAlmostEqual(t, got, tt.want, fmt.Sprintf("EvalAST(%q)", tt.input))
		})
	}
}

// TestParseComplexUnit verifies the main entry point function
func TestParseComplexUnit(t *testing.T) {
	ctx := NewStandardContext()

	tests := []struct {
		input string
		want  Unit
	}{
		{"m", Unit{1, Dimension{1, 0, 0, 0, 0, 0, 0}}},
		{"kg*m/s^2", Unit{1, Dimension{1, 1, -2, 0, 0, 0, 0}}},
		{"km/h", Unit{1000.0 / 3600.0, Dimension{1, 0, -1, 0, 0, 0, 0}}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseComplexUnit(tt.input, ctx)
			if err != nil {
				t.Fatalf("ParseComplexUnit(%q) error: %v", tt.input, err)
			}

			assertUnitAlmostEqual(t, got, tt.want, fmt.Sprintf("ParseComplexUnit(%q)", tt.input))
		})
	}
}

// TestParseInvalidUnits verifies error handling for invalid units
func TestParseInvalidUnits(t *testing.T) {
	ctx := NewStandardContext()

	tests := []string{
		"m/",    // Missing denominator
		"*m",    // Invalid operator usage
		"(kg",   // Unclosed parenthesis
		"kg^x",  // Non-numeric exponent
		"m^^2",  // Double power operator
		"xyzzy", // Unknown unit
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseComplexUnit(input, ctx)
			if err == nil {
				t.Errorf("ParseComplexUnit(%q) expected error, got nil", input)
			}
		})
	}
}
