package si

import (
	"fmt"
	"text/scanner"
)

// TokenKind represents the type of token during lexical analysis
type TokenKind int

const (
	// Token types
	Invalid TokenKind = iota
	EOF
	Identifier
	Number
	Multiply // * or ·
	Divide   // /
	Power    // ^
	LParen   // (
	RParen   // )
)

// Token represents a lexical token
type Token struct {
	Kind  TokenKind
	Value string
	Pos   scanner.Position
}

// String returns a string representation of the token
func (t Token) String() string {
	switch t.Kind {
	case EOF:
		return "EOF"
	case Invalid:
		return fmt.Sprintf("Invalid(%s)", t.Value)
	default:
		return fmt.Sprintf("%v(%s)", t.Kind, t.Value)
	}
}

// String returns the name of the token kind
func (k TokenKind) String() string {
	switch k {
	case Invalid:
		return "Invalid"
	case EOF:
		return "EOF"
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case Power:
		return "Power"
	case LParen:
		return "LParen"
	case RParen:
		return "RParen"
	default:
		return fmt.Sprintf("TokenKind(%d)", k)
	}
}
