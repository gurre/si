Design Specification: Package Structure for SI Unit Library
==========================================================

Overview
--------
This library models SI units, parses and formats unit expressions, and handles
unit arithmetic and conversion. The package structure enforces separation of
concerns and avoids cyclic dependencies.

Root Package:    si
Supporting:      si/parser, si/format, si/dim, si/symbols, si/convert

-------------------------------------------
1. Root Package: si
-------------------------------------------
Purpose:
    - Defines core Unit and Dimension types
    - Exposes construction, parsing, formatting, and arithmetic
    - Provides SI constants and helpers

Key Types:
    type Unit struct {
        Value     float64
        Dimension [7]int
    }

    // Constructors and helpers
    func Parse(input string) (Unit, error)
    func Kilograms(f float64) Unit
    func Meters(f float64) Unit
    // ... other SI unit constructors

Key Rules:
    - Does not import any supporting package.
    - All public usage is via this package (other packages are utility).

-------------------------------------------
2. Parser Package: si/parser
-------------------------------------------
Purpose:
    - Provides AST-based parsing from unit string to AST and to Unit.
    - Isolated from root (except for importing Unit definition).

Interface:
    type Node interface {
        Eval(ctx Context) (si.Unit, error)
        String() string
    }

    type Parser interface {
        Parse(input string) (Node, error)
    }

    type Context interface {
        Resolve(symbol string) (si.Unit, error)
    }

    // Main entry
    func ParseComplexUnit(input string) (si.Unit, error)

Relationships:
    - Imports "github.com/gurre/si"
    - Used by si.Parse and by any advanced users
    - No imports from format/dim/symbols/convert

-------------------------------------------
3. Format Package: si/format
-------------------------------------------
Purpose:
    - Formats Units or AST Nodes to strings in various notations.

Interface:
    type Formatter interface {
        FormatNode(node parser.Node, opts *FormatOptions) (string, error)
        FormatUnit(u si.Unit, opts *FormatOptions) (string, error)
    }

    type FormatOptions struct {
        MultSymbol  string // e.g. "*", "·"
        DivSymbol   string // e.g. "/"
        ExponentFmt string // e.g. "^%d"
        Parens      bool
        CollapseSymbols bool // Prefer N over kg*m/s^2
    }

Relationships:
    - Imports "github.com/gurre/si"
    - Imports "github.com/gurre/si/parser"
    - May import "github.com/gurre/si/symbols" for symbol lookup

-------------------------------------------
4. Dimension Package: si/dim
-------------------------------------------
Purpose:
    - Provides math helpers for [7]int dimensions.

Interface:
    type Dimension [7]int

    func Add(a, b Dimension) Dimension
    func Sub(a, b Dimension) Dimension
    func Mul(a Dimension, k int) Dimension
    func Eq(a, b Dimension) bool
    func Zero() Dimension

Relationships:
    - Imports nothing or just "github.com/gurre/si" for type compatibility
    - Used by si.Unit methods and optionally parser

-------------------------------------------
5. Symbols Package: si/symbols
-------------------------------------------
Purpose:
    - Handles mapping between SI dimensions and canonical unit symbols.
    - Manages symbol registration for custom or derived units.

Interface:
    func Lookup(dim dim.Dimension) (symbol string, ok bool)
    func Register(dim dim.Dimension, symbol string)

Relationships:
    - Imports "github.com/gurre/si/dim"
    - Optionally used by si/format

-------------------------------------------
6. Convert Package: si/convert
-------------------------------------------
Purpose:
    - Handles unit conversions, SI/metric/binary prefix math.

Interface:
    func PrefixFactor(prefix string) (float64, bool)
    func ConvertTo(u si.Unit, target si.Unit) (si.Unit, error)

Relationships:
    - Imports "github.com/gurre/si"
    - May be called by si.Unit methods

-------------------------------------------
Relationships Summary Table
-------------------------------------------
| Package      | Depends on     | Used by      |
|--------------|---------------|--------------|
| si           | —             | All users    |
| si/parser    | si            | si, format   |
| si/format    | si, parser    | si, users    |
| si/dim       | (si, none)    | si, parser   |
| si/symbols   | dim           | format, si   |
| si/convert   | si            | si           |

-------------------------------------------
Encapsulation and Idiomatic API
-------------------------------------------
- Only `si.Unit` and main constructors/helpers are public API.
- Supporting packages are for internal extension, advanced use, and testing.
- Supporting packages never import each other—only the root si or lowest level needed.

-------------------------------------------
Example: How a user consumes the API
-------------------------------------------
import \"github.com/gurre/si\"

u, err := si.Parse(\"kg*m/s^2\")
if err != nil { ... }
fmt.Println(u.Format(nil)) // uses si/format under the hood

-------------------------------------------
Example: Extending with custom formatter
-------------------------------------------
import (
    \"github.com/gurre/si\"
    \"github.com/gurre/si/parser\"
    \"github.com/gurre/si/format\"
)

node, err := parser.Parse(\"(N*s)/m\")
f := &format.DefaultFormatter{}
str, _ := f.FormatNode(node, &format.FormatOptions{...})
fmt.Println(str)

-------------------------------------------
Summary
-------------------------------------------
- Single entrypoint for end-users: si.Unit, si.Parse, helpers
- AST, parser, formatter, dimension helpers, and symbol registry are modular and testable
- No cyclic dependencies, easy to extend or swap formatter/parser
- Consistent interfaces and clean API for advanced or internal use