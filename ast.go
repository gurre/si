package si

import (
	"fmt"
)

// Node represents a node in the abstract syntax tree
type Node interface {
	// Eval evaluates the node to produce a Unit
	Eval(ctx Context) (Unit, error)
	// String returns a string representation of the node
	String() string
}

// IdentNode represents an identifier in the expression (unit symbol)
type IdentNode struct {
	Symbol string
}

// Eval resolves the identifier to a Unit
func (n *IdentNode) Eval(ctx Context) (Unit, error) {
	return ctx.Resolve(n.Symbol)
}

// String returns the symbol name
func (n *IdentNode) String() string {
	return n.Symbol
}

// NumberNode represents a numeric literal in the expression
type NumberNode struct {
	Value float64
}

// Eval returns a dimensionless Unit with the value
func (n *NumberNode) Eval(ctx Context) (Unit, error) {
	return Scalar(n.Value), nil
}

// String returns the string representation of the number
func (n *NumberNode) String() string {
	return fmt.Sprintf("%g", n.Value)
}

// BinaryNode represents a binary operation (multiplication or division)
type BinaryNode struct {
	Op    TokenKind
	Left  Node
	Right Node
}

// Eval evaluates both sides of the binary operation and applies the operation
func (n *BinaryNode) Eval(ctx Context) (Unit, error) {
	left, err := n.Left.Eval(ctx)
	if err != nil {
		return Unit{}, fmt.Errorf("error evaluating left side: %w", err)
	}

	right, err := n.Right.Eval(ctx)
	if err != nil {
		return Unit{}, fmt.Errorf("error evaluating right side: %w", err)
	}

	switch n.Op {
	case Multiply:
		return left.Mul(right), nil
	case Divide:
		return left.Div(right), nil
	default:
		return Unit{}, fmt.Errorf("unsupported binary operation: %v", n.Op)
	}
}

// String returns a string representation of the binary operation
func (n *BinaryNode) String() string {
	op := "*"
	if n.Op == Divide {
		op = "/"
	}
	return fmt.Sprintf("(%s %s %s)", n.Left, op, n.Right)
}

// PowerNode represents an expression raised to a power
type PowerNode struct {
	Base Node
	Exp  int
}

// Eval evaluates the base and raises it to the power
func (n *PowerNode) Eval(ctx Context) (Unit, error) {
	base, err := n.Base.Eval(ctx)
	if err != nil {
		return Unit{}, fmt.Errorf("error evaluating base: %w", err)
	}

	return base.Pow(n.Exp), nil
}

// String returns a string representation of the power operation
func (n *PowerNode) String() string {
	return fmt.Sprintf("%s^%d", n.Base, n.Exp)
}

// GroupNode represents a grouped expression (in parentheses)
type GroupNode struct {
	Inner Node
}

// Eval evaluates the inner expression
func (n *GroupNode) Eval(ctx Context) (Unit, error) {
	return n.Inner.Eval(ctx)
}

// String returns a string representation of the grouped expression
func (n *GroupNode) String() string {
	return fmt.Sprintf("(%s)", n.Inner)
}

// Context provides resolution of units and prefixes
type Context interface {
	// Resolve converts a symbol to a Unit
	Resolve(symbol string) (Unit, error)
}
