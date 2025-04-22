package si

import (
	"fmt"
	"strings"
)

// Formatter provides a contract for formatting AST nodes into strings
type Formatter interface {
	Format(Node) (string, error)
}

// FormatOptions configures how units are formatted
type FormatOptions struct {
	// MultSymbol defines the multiplication symbol (default "*")
	MultSymbol string
	// DivSymbol defines the division symbol (default "/")
	DivSymbol string
	// ExponentFmt defines the exponent format (default "^%d")
	ExponentFmt string
	// UseParens determines if parentheses should be used (default true)
	UseParens bool
	// Simplify controls whether to simplify units (default false)
	Simplify bool
	// CollapseSymbols controls whether to use known symbolic names (default true)
	CollapseSymbols bool
	// KnownSymbols maps dimensions to their symbolic names
	KnownSymbols map[Dimension]string
}

// DefaultFormatOptions returns the default formatting options
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		MultSymbol:      "*",
		DivSymbol:       "/",
		ExponentFmt:     "^%d",
		UseParens:       true,
		Simplify:        false,
		CollapseSymbols: true,
		KnownSymbols:    defaultKnownSymbols(),
	}
}

// defaultKnownSymbols returns a map of common dimensions to their symbolic names
func defaultKnownSymbols() map[Dimension]string {
	return map[Dimension]string{
		Newton.Dimension: "N",
		Joule.Dimension:  "J",
		Watt.Dimension:   "W",
		Pascal.Dimension: "Pa",
		Hertz.Dimension:  "Hz",
		Volt.Dimension:   "V",
	}
}

// DefaultFormatter implements the Formatter interface with default settings
type DefaultFormatter struct {
	Options FormatOptions
}

// NewDefaultFormatter creates a new DefaultFormatter with default options
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		Options: DefaultFormatOptions(),
	}
}

// Format formats a node into a string using the default options
func (f *DefaultFormatter) Format(node Node) (string, error) {
	if node == nil {
		return "", fmt.Errorf("cannot format nil node")
	}

	switch n := node.(type) {
	case *IdentNode:
		return n.Symbol, nil

	case *NumberNode:
		return fmt.Sprintf("%g", n.Value), nil

	case *BinaryNode:
		left, err := f.Format(n.Left)
		if err != nil {
			return "", err
		}

		right, err := f.Format(n.Right)
		if err != nil {
			return "", err
		}

		op := f.Options.MultSymbol
		if n.Op == Divide {
			op = f.Options.DivSymbol
		}

		// Apply parentheses if needed
		if f.Options.UseParens {
			// Only add parentheses around the left side if it's a binary operation
			// and we're doing division. This prevents output like "(kg*m)/s^2".
			if isBinaryNode(n.Left) && n.Op == Divide && !isBinaryDivision(n.Left.(*BinaryNode)) {
				left = "(" + left + ")"
			}

			// Always add parentheses around the right side in division if it's a binary operation
			if isBinaryNode(n.Right) && n.Op == Divide {
				right = "(" + right + ")"
			}
		}

		return left + op + right, nil

	case *PowerNode:
		base, err := f.Format(n.Base)
		if err != nil {
			return "", err
		}

		// Don't show exponent 1 if simplify is enabled
		if n.Exp == 1 && f.Options.Simplify {
			return base, nil
		}

		// Apply parentheses if base is a binary operation
		if f.Options.UseParens && isBinaryNode(n.Base) {
			base = "(" + base + ")"
		}

		return base + fmt.Sprintf(f.Options.ExponentFmt, n.Exp), nil

	case *GroupNode:
		inner, err := f.Format(n.Inner)
		if err != nil {
			return "", err
		}

		if f.Options.UseParens {
			return "(" + inner + ")", nil
		}
		return inner, nil

	default:
		return "", fmt.Errorf("unknown node type: %T", node)
	}
}

// isBinaryNode checks if a node is a binary operation
func isBinaryNode(node Node) bool {
	_, ok := node.(*BinaryNode)
	return ok
}

// isBinaryDivision checks if a binary node is a division operation
func isBinaryDivision(node *BinaryNode) bool {
	return node.Op == Divide
}

// FormatUnit formats a Unit into a readable string
func FormatUnit(u Unit) string {
	// If dimensionless, just return the value
	if u.Dimension == Dimensionless {
		return fmt.Sprintf("%g", u.Value)
	}

	// Try to use formatter to generate a unit string
	fmtStr, err := formatUnitDimension(u)
	if err != nil {
		// Fallback to a simple dimension representation
		fmtStr = formatDimensionFallback(u.Dimension)
	}

	// Format with value if needed
	if u.Value != 1.0 {
		// Special formatting for energy density (pressure)
		if u.Dimension == Pascal.Dimension {
			return fmt.Sprintf("%g Pa", u.Value)
		}
		return fmt.Sprintf("%g %s", u.Value, fmtStr)
	}

	return fmtStr
}

// formatUnitDimension attempts to format a unit using the formatter
func formatUnitDimension(u Unit) (string, error) {
	// Check for known symbolic units first
	formatter := NewDefaultFormatter()

	if formatter.Options.CollapseSymbols {
		if symbol, ok := formatter.Options.KnownSymbols[u.Dimension]; ok {
			return symbol, nil
		}
	}

	// Special case handling for complex units
	if Watt.Div(Meter.Mul(Kelvin)).Dimension == u.Dimension {
		return "W/(m*K)", nil // Thermal conductivity
	} else if Joule.Div(Kilogram.Mul(Kelvin)).Dimension == u.Dimension {
		return "J/(kg*K)", nil // Specific heat capacity
	} else if u.Dimension == Joule.Div(Meter.Pow(3)).Dimension {
		return "Pa", nil // Energy density is equivalent to pressure
	}

	// Generate an AST for this dimension
	node, err := dimensionToAST(u.Dimension)
	if err != nil {
		return "", err
	}

	// Format the AST
	return formatter.Format(node)
}

// dimensionToAST converts a Dimension to an AST node
func dimensionToAST(dim Dimension) (Node, error) {
	symbols := []string{"m", "kg", "s", "A", "K", "mol", "cd"}
	var numerator []Node
	var denominator []Node

	// Process dimensions in a specific order to ensure consistent output
	// First add mass, then length, then other positive dimensions
	// This ensures kg*m instead of m*kg

	// First pass: add kg if present (mass)
	if dim[1] > 0 {
		identNode := &IdentNode{Symbol: symbols[1]}
		if dim[1] != 1 {
			numerator = append(numerator, &PowerNode{
				Base: identNode,
				Exp:  dim[1],
			})
		} else {
			numerator = append(numerator, identNode)
		}
	}

	// Second pass: add m if present (length)
	if dim[0] > 0 {
		identNode := &IdentNode{Symbol: symbols[0]}
		if dim[0] != 1 {
			numerator = append(numerator, &PowerNode{
				Base: identNode,
				Exp:  dim[0],
			})
		} else {
			numerator = append(numerator, identNode)
		}
	}

	// Third pass: add other dimensions in order
	for i, exp := range dim {
		if i == 0 || i == 1 || exp == 0 {
			continue // Skip length, mass (already handled) and zero exponents
		}

		identNode := &IdentNode{Symbol: symbols[i]}

		if exp > 0 {
			// Add exponent if not 1
			if exp != 1 {
				numerator = append(numerator, &PowerNode{
					Base: identNode,
					Exp:  exp,
				})
			} else {
				numerator = append(numerator, identNode)
			}
		} else {
			// Add exponent if not -1
			if exp != -1 {
				denominator = append(denominator, &PowerNode{
					Base: identNode,
					Exp:  -exp,
				})
			} else {
				denominator = append(denominator, identNode)
			}
		}
	}

	// Now handle negative exponents for length and mass
	if dim[0] < 0 {
		identNode := &IdentNode{Symbol: symbols[0]}
		if dim[0] != -1 {
			denominator = append(denominator, &PowerNode{
				Base: identNode,
				Exp:  -dim[0],
			})
		} else {
			denominator = append(denominator, identNode)
		}
	}

	if dim[1] < 0 {
		identNode := &IdentNode{Symbol: symbols[1]}
		if dim[1] != -1 {
			denominator = append(denominator, &PowerNode{
				Base: identNode,
				Exp:  -dim[1],
			})
		} else {
			denominator = append(denominator, identNode)
		}
	}

	// Handle special cases
	if len(numerator) == 0 && len(denominator) == 0 {
		return &NumberNode{Value: 1}, nil
	}

	if len(numerator) == 0 {
		numerator = append(numerator, &NumberNode{Value: 1})
	}

	// Build the numerator part
	var numNode Node
	if len(numerator) == 1 {
		numNode = numerator[0]
	} else {
		// Chain binary multiplication nodes
		numNode = numerator[0]
		for i := 1; i < len(numerator); i++ {
			numNode = &BinaryNode{
				Op:    Multiply,
				Left:  numNode,
				Right: numerator[i],
			}
		}
	}

	// If there are no denominators, return just the numerator
	if len(denominator) == 0 {
		return numNode, nil
	}

	// Build the denominator part
	var denomNode Node
	if len(denominator) == 1 {
		denomNode = denominator[0]
	} else {
		// Chain binary multiplication nodes
		denomNode = denominator[0]
		for i := 1; i < len(denominator); i++ {
			denomNode = &BinaryNode{
				Op:    Multiply,
				Left:  denomNode,
				Right: denominator[i],
			}
		}
	}

	// Combine numerator and denominator with division
	return &BinaryNode{
		Op:    Divide,
		Left:  numNode,
		Right: denomNode,
	}, nil
}

// formatDimensionFallback provides a fallback dimension formatter
func formatDimensionFallback(d Dimension) string {
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
	numStr := strings.Join(numerator, "*")

	// If there are denominators, append them with "/"
	if len(denominator) > 0 {
		denomStr := strings.Join(denominator, "*")
		return fmt.Sprintf("%s/%s", numStr, denomStr)
	}

	return numStr
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// FormatAST formats an AST node using the provided options
func FormatAST(node Node, opts *FormatOptions) (string, error) {
	if opts == nil {
		defOpts := DefaultFormatOptions()
		opts = &defOpts
	}

	formatter := &DefaultFormatter{
		Options: *opts,
	}

	return formatter.Format(node)
}

// FormatUnitWithOptions formats a Unit using custom options
func FormatUnitWithOptions(u Unit, opts *FormatOptions) string {
	if opts == nil {
		defOpts := DefaultFormatOptions()
		opts = &defOpts
	}

	// If dimensionless, just return the value
	if u.Dimension == Dimensionless {
		return fmt.Sprintf("%g", u.Value)
	}

	// Try to use formatter with custom options
	formatter := &DefaultFormatter{
		Options: *opts,
	}

	// Check for known symbolic units first
	if formatter.Options.CollapseSymbols {
		if symbol, ok := formatter.Options.KnownSymbols[u.Dimension]; ok {
			if u.Value != 1.0 {
				return fmt.Sprintf("%g %s", u.Value, symbol)
			}
			return symbol
		}
	}

	// Generate an AST for this dimension
	node, err := dimensionToAST(u.Dimension)
	if err != nil {
		// Fallback to simple format
		unitStr := formatDimensionFallback(u.Dimension)
		if u.Value != 1.0 {
			return fmt.Sprintf("%g %s", u.Value, unitStr)
		}
		return unitStr
	}

	// Format the AST
	unitStr, err := formatter.Format(node)
	if err != nil {
		unitStr = formatDimensionFallback(u.Dimension)
	}

	if u.Value != 1.0 {
		return fmt.Sprintf("%g %s", u.Value, unitStr)
	}

	return unitStr
}
