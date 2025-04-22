package si

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
	"unicode/utf8"
)

// Tokenizer turns an input string into a sequence of tokens
type Tokenizer struct {
	input    string
	tokens   []Token
	position int
}

// NewTokenizer creates a new tokenizer for the input
func NewTokenizer(input string) *Tokenizer {
	tokens, _ := tokenizeFully(input)
	return &Tokenizer{
		input:    input,
		tokens:   tokens,
		position: 0,
	}
}

// Next returns the next token and advances
func (t *Tokenizer) Next() Token {
	if t.position >= len(t.tokens) {
		// Return EOF if we're past the end
		return Token{
			Kind: EOF,
			Pos:  scanner.Position{},
		}
	}

	token := t.tokens[t.position]
	t.position++
	return token
}

// Peek returns the next token without advancing
func (t *Tokenizer) Peek() Token {
	if t.position >= len(t.tokens) {
		// Return EOF if we're past the end
		return Token{
			Kind: EOF,
			Pos:  scanner.Position{},
		}
	}

	return t.tokens[t.position]
}

// readRuneAt safely reads a rune from a string at a given byte position
func readRuneAt(s string, pos int) (rune, int) {
	if pos >= len(s) {
		return 0, 0 // End of string
	}
	return utf8.DecodeRuneInString(s[pos:])
}

// tokenizeFully tokenizes the entire input at once
func tokenizeFully(input string) ([]Token, error) {
	// Normalize the input to make parsing easier
	input = normalizeInput(input)

	var tokens []Token
	var pos int

	for pos < len(input) {
		// Skip spaces
		r, width := readRuneAt(input, pos)
		if isSpace(r) {
			pos += width
			continue
		}

		// Single-character tokens
		if r == '(' {
			tokens = append(tokens, Token{Kind: LParen, Value: "(", Pos: scanner.Position{Offset: pos}})
			pos += width
			continue
		} else if r == ')' {
			tokens = append(tokens, Token{Kind: RParen, Value: ")", Pos: scanner.Position{Offset: pos}})
			pos += width
			continue
		} else if r == '*' || r == '·' {
			val := string(r)
			tokens = append(tokens, Token{Kind: Multiply, Value: val, Pos: scanner.Position{Offset: pos}})
			pos += width
			continue
		} else if r == '/' {
			tokens = append(tokens, Token{Kind: Divide, Value: "/", Pos: scanner.Position{Offset: pos}})
			pos += width
			continue
		} else if r == '^' {
			tokens = append(tokens, Token{Kind: Power, Value: "^", Pos: scanner.Position{Offset: pos}})
			pos += width
			continue
		}

		// Numbers
		if unicode.IsDigit(r) {
			start := pos
			for pos < len(input) {
				r, width := readRuneAt(input, pos)
				if !unicode.IsDigit(r) && r != '.' {
					break
				}
				pos += width
			}
			numStr := input[start:pos]
			_, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return tokens, fmt.Errorf("invalid number %q at position %d", numStr, start)
			}
			tokens = append(tokens, Token{Kind: Number, Value: numStr, Pos: scanner.Position{Offset: start}})
			continue
		}

		// Identifiers
		if unicode.IsLetter(r) || isSpecialIdentifierStart(r) {
			start := pos
			for pos < len(input) {
				r, width := readRuneAt(input, pos)
				if !isIdentifierChar(r) {
					break
				}
				pos += width
			}
			ident := input[start:pos]
			tokens = append(tokens, Token{Kind: Identifier, Value: ident, Pos: scanner.Position{Offset: start}})
			continue
		}

		// Invalid character
		return tokens, fmt.Errorf("invalid character %q at position %d", r, pos)
	}

	// Add EOF token
	tokens = append(tokens, Token{Kind: EOF, Pos: scanner.Position{Offset: pos}})

	return tokens, nil
}

// isSpecialIdentifierStart checks if a rune is a valid start of an identifier (special characters)
func isSpecialIdentifierStart(r rune) bool {
	return r == '%' || r == '°' || r == 'µ' || r == 'μ' || r == 'Ω'
}

// Helpers for tokenization
func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isIdentifierChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || isSpecialIdentifierStart(r)
}

// normalizeInput prepares input for tokenization
func normalizeInput(input string) string {
	// Add spaces around operators for clear tokenization
	input = strings.Replace(input, "(", " ( ", -1)
	input = strings.Replace(input, ")", " ) ", -1)
	input = strings.Replace(input, "/", " / ", -1)
	input = strings.Replace(input, "*", " * ", -1)
	input = strings.Replace(input, "·", " · ", -1)
	input = strings.Replace(input, "^", " ^ ", -1)

	// Normalize spaces
	for strings.Contains(input, "  ") {
		input = strings.Replace(input, "  ", " ", -1)
	}
	return strings.TrimSpace(input)
}

// tokenize is a legacy function for testing
func tokenize(input string) ([]Token, error) {
	return tokenizeFully(input)
}
