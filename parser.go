package si

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser implements a recursive descent parser for unit expressions
type Parser struct {
	tokenizer *Tokenizer
	err       error
}

// NewParser creates a new parser for the input
func NewParser(input string) *Parser {
	// Normalize middle dot to asterisk for consistency
	input = strings.Replace(input, "·", "*", -1)

	// Preprocess for better tokenizing of operators
	// Add spaces around operators and parentheses for easier tokenization
	return &Parser{
		tokenizer: NewTokenizer(input),
	}
}

// Parse parses a unit expression and returns the AST
func (p *Parser) Parse() (Node, error) {
	node := p.parseTerm()
	if p.err != nil {
		return nil, p.err
	}

	// Check for unconsumed tokens
	token := p.tokenizer.Next()
	if token.Kind != EOF {
		return nil, fmt.Errorf("unexpected token at end of input: %s", token.Value)
	}

	return node, nil
}

// parseTerm parses a term (multiplication/division)
func (p *Parser) parseTerm() Node {
	if p.err != nil {
		return nil
	}

	left := p.parsePower()
	if p.err != nil {
		return nil
	}

	for {
		token := p.tokenizer.Peek()
		if token.Kind != Multiply && token.Kind != Divide {
			break
		}

		p.tokenizer.Next() // Consume the operator

		right := p.parsePower()
		if p.err != nil {
			return nil
		}

		left = &BinaryNode{
			Op:    token.Kind,
			Left:  left,
			Right: right,
		}
	}

	return left
}

// parsePower parses a power expression (factor^exponent)
func (p *Parser) parsePower() Node {
	if p.err != nil {
		return nil
	}

	base := p.parseFactor()
	if p.err != nil {
		return nil
	}

	// Check for power operator
	token := p.tokenizer.Peek()
	if token.Kind == Power {
		p.tokenizer.Next() // Consume ^

		// Parse exponent (must be a number)
		token = p.tokenizer.Next()
		if token.Kind != Number {
			p.err = fmt.Errorf("expected number for exponent, got %v", token)
			return nil
		}

		exp, err := strconv.Atoi(token.Value)
		if err != nil {
			p.err = fmt.Errorf("invalid exponent %q: %w", token.Value, err)
			return nil
		}

		return &PowerNode{
			Base: base,
			Exp:  exp,
		}
	}

	return base
}

// parseFactor parses an atomic expression or parenthesized expression
func (p *Parser) parseFactor() Node {
	if p.err != nil {
		return nil
	}

	token := p.tokenizer.Next()

	switch token.Kind {
	case Number:
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			p.err = fmt.Errorf("invalid number %q: %w", token.Value, err)
			return nil
		}
		return &NumberNode{Value: value}

	case Identifier:
		return &IdentNode{Symbol: token.Value}

	case LParen:
		expr := p.parseTerm()
		if p.err != nil {
			return nil
		}

		token = p.tokenizer.Next()
		if token.Kind != RParen {
			p.err = fmt.Errorf("expected closing parenthesis, got %v", token)
			return nil
		}

		return &GroupNode{Inner: expr}

	default:
		p.err = fmt.Errorf("unexpected token %v", token)
		return nil
	}
}

// ParseUnitAST parses a unit expression into an AST
func ParseUnitAST(input string) (Node, error) {
	// Replace middle dot with asterisk for consistency
	input = strings.Replace(input, "·", "*", -1)

	parser := NewParser(input)
	return parser.Parse()
}

// EvalAST evaluates an AST with the given context
func EvalAST(node Node, ctx Context) (Unit, error) {
	if node == nil {
		return Unit{}, fmt.Errorf("cannot evaluate nil node")
	}
	return node.Eval(ctx)
}

// ParseComplexUnit is the main entry point for parsing a unit expression
func ParseComplexUnit(input string, ctx Context) (Unit, error) {
	// Replace middle dot with asterisk for consistency
	input = strings.Replace(input, "·", "*", -1)

	ast, err := ParseUnitAST(input)
	if err != nil {
		return Unit{}, err
	}

	return EvalAST(ast, ctx)
}

// preprocessInput prepares a unit expression for tokenization
func preprocessInput(input string) string {
	// Add spaces around operators and parentheses for easier tokenization
	input = addSpacesAroundOperators(input)

	// Remove redundant spaces
	return normalizeSpaces(input)
}

// addSpacesAroundOperators adds spaces around operators and parentheses
func addSpacesAroundOperators(input string) string {
	input = addSpaceAround(input, "(")
	input = addSpaceAround(input, ")")
	input = addSpaceAround(input, "/")
	input = addSpaceAround(input, "*")
	input = addSpaceAround(input, "^")
	return input
}

// addSpaceAround adds spaces around a specific character
func addSpaceAround(input, char string) string {
	return strings.Replace(input, char, " "+char+" ", -1)
}

// normalizeSpaces removes redundant spaces
func normalizeSpaces(input string) string {
	for strings.Contains(input, "  ") {
		input = strings.Replace(input, "  ", " ", -1)
	}
	return strings.TrimSpace(input)
}
