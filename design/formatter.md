Design Specification: AST-Based Unit Formatter
==============================================

Overview
--------
The Unit Formatter is responsible for rendering Unit values into readable,
standardized strings. It works by traversing the Abstract Syntax Tree (AST)
of unit expressions and generating a string based on formatting rules.

This formatter supports canonical SI forms, user-friendly display, and
custom formatting profiles (e.g., scientific vs engineering style).

Goals
-----
- Generate human-readable unit strings (e.g., "kg*m/s^2", "W/(m^2*K^4)")
- Support full SI dimensions, derived units, and symbolic names
- Allow configurable formatting style (symbols, multiplication, division)
- Integrate cleanly with the existing AST parser
- Round-trip: parsing and formatting should be reversible when possible

Architecture
------------
The formatter architecture includes the following components:

1. **Formatter Interface**
    Provides a pluggable contract for formatting nodes:

        type Formatter interface {
            Format(Node) (string, error)
        }

2. **DefaultFormatter**
    Implements `Formatter` with canonical SI output.
    Example: `(kg*m)/(s^2)` → "kg*m/s^2"

3. **Node Traversal**
    Uses a visitor pattern to walk the AST nodes and emit a string expression:

        - IdentNode: returns the base symbol
        - BinaryNode: recursively formats left/right with * or /
        - PowerNode: appends ^exp
        - GroupNode: wraps expression in ()

4. **Configurable Output Options**
    The formatter may accept an optional config:

        type FormatOptions struct {
            MultSymbol  string // default "*"
            DivSymbol   string // default "/"
            ExponentFmt string // default "^%d"
            UseParens   bool   // default true
            Simplify    bool   // default false (show m^1, kg^1, etc.)
        }

5. **Context-Aware Symbols**
    Formatter will try to collapse known dimensions to symbolic names (e.g., "N" for "kg·m/s^2") if configured:

        FormatterOptions{
            CollapseSymbols: true,
            KnownSymbols: map[Dimension]string{
                Newton.Dimension: "N",
            },
        }

Function Signatures
-------------------
Primary formatter API:

    func FormatUnit(u Unit) string

With AST-level control:

    func FormatAST(node Node, opts *FormatOptions) (string, error)

With pluggable formatters:

    type Formatter interface {
        Format(Node) (string, error)
    }

    var f Formatter = &DefaultFormatter{}
    str, err := f.Format(astNode)

Examples of Output Strings
--------------------------
    Unit{Value: 1, Dimension: [1 0 -2 0 0 0 0]}        → "m/s^2"
    Unit{Value: 1, Dimension: [1 1 -2 0 0 0 0]}        → "kg*m/s^2"
    AST: Binary(Div, Power(Ident(m), 2), Power(Ident(s), 2)) → "m^2/s^2"
    AST: Group(Binary(Mul, Ident(kg), Ident(m)))      → "(kg*m)"
    With known symbols:                               → "N"

Formatter Modes (configurable)
------------------------------
- SI Canonical: "kg*m/s^2"
- Fully Expanded: "m^1*s^-2*kg^1"
- Pretty Print: "N·m/s"
- JSON/Code:  "kg*m*s^-2" (no division)

Extensibility
-------------
- Support units in numerator/denominator split view
- Optionally render prefix and numeric value (e.g., "1.23 km/h")
- Allow localization of symbols (e.g., Unicode `·`, localized decimals)

Future Enhancements
-------------------
- Simplify dimensions into canonical SI derived unit symbols
- Format expressions to LaTeX (e.g., `\\frac{kg\\cdot m}{s^2}`)
- Round-trip testing: Parse(Format(Parse(x))) == Parse(x)
- Support alternate output formats (JSON, YAML, Markdown)