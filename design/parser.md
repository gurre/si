Design Specification: AST-Based Unit Expression Parser
=======================================================

Overview
--------
This design outlines a flexible and extensible architecture for parsing unit
expressions (e.g., "kg·m/s^2", "W/(m^2·K^4)"). The design is based on a
tokenizer + abstract syntax tree (AST) model to decouple lexing, parsing, and
evaluation. This allows better error handling, testability, and future
extensions (e.g., symbolic manipulation, simplification, serialization).

Goals
-----
- Provide a clean and testable architecture for parsing units
- Support all SI base and derived units with prefixes
- Allow nested expressions, grouping, and exponentiation
- Return detailed errors, not panics
- Enable semantic validation at the AST level
- Allow reuse across custom unit domains (e.g., astronomy, data science)

Components
----------
1. **Tokenizer (Lexer)**
    - Based on `text/scanner.Scanner`
    - Emits a sequence of tokens (identifiers, numbers, symbols)
    - Skips whitespace and validates characters

2. **AST Nodes**
    - Each expression is parsed into a tree of typed nodes:
      ```
      type Node interface {
          Eval(ctx Context) (Unit, error)
      }

      type IdentNode struct {
          Symbol string
      }

      type BinaryNode struct {
          Op    TokenKind // '*', '/', etc.
          Left  Node
          Right Node
      }

      type PowerNode struct {
          Base Node
          Exp  int
      }

      type GroupNode struct {
          Inner Node
      }
      ```

3. **Parser**
    - Recursive descent implementation that builds AST from tokens
    - Handles operator precedence and parentheses
    - Validates structure and attaches context-independent error messages

4. **Context (for evaluation)**
    - Provides mapping of known units, prefixes, and symbolic constants
    - Resolves IdentNode during evaluation
      ```
      type Context interface {
          Resolve(symbol string) (Unit, error)
      }
      ```

5. **Evaluator**
    - Walks the AST and returns a final `Unit`
    - Applies correct dimensional math and value scaling
    - Returns precise error messages (e.g., division by zero, unknown symbol)

Function Signatures
-------------------
Preferred entry point:
    func ParseComplexUnit(input string) (Unit, error)

AST-specific (for testing or meta usage):
    func ParseUnitAST(input string) (Node, error)
    func EvalAST(node Node, ctx Context) (Unit, error)

Examples of Input Strings
-------------------------
These should be supported and tested:
    "m"
    "s^2"
    "kg*m/s^2"
    "((kg*m)/s)^2"
    "W/(m^2*K^4)"
    "km/h"
    "g/(cm^2*s)"
    "(N*s)/m"
    "(kg*m)/(s^2*K)"
    "((W/m^2)^2)^3"

Extensibility
-------------
- Additional operators can be added (e.g., `%`, `+` if needed)
- Symbolic simplification routines can be layered on top of AST
- AST nodes can support formatting, stringification, metadata

Future Considerations
---------------------
- Implement AST serialization to JSON/YAML for debugging
- Support symbolic units (e.g., Planck constant, lightyear)
- Enable constant folding in AST
- Integrate with code generation tools (e.g., automatic unit test generators)


Implementation
--------------
1. The unit system is built around the Unit struct, which contains:
    - Value: A float64 representing the scalar magnitude of the physical quantity
    - Dimension: A [7]int array representing the exponents of the 7 SI base dimensions (Length, Mass, Time, Current, Temperature, Substance, Luminosity)
2. Units are registered in the StandardContext struct, which implements the Context interface. The StandardContext contains:
    - baseUnits: A map of base unit symbols (e.g., "m", "kg") to their Unit definitions
    - derivedUnits: A map of derived unit symbols (e.g., "N", "J") to their Unit definitions
    - prefixes: A map of unit prefixes (e.g., "k", "M") to their scaling factors
    - sortedPrefixes: A slice of prefixes sorted by length for proper matching
3. Registration happens in three main methods:
    - registerBaseUnits(): Registers the 7 SI base units (meter, kilogram, second, ampere, kelvin, mole, candela)
    - registerDerivedUnits(): Registers derived units (newton, joule, watt, etc.) by combining base units
    - registerPrefixes(): Registers SI prefixes (kilo, mega, etc.) and binary prefixes (kibi, mebi, etc.)
4. The entire process is initiated in the NewStandardContext() function, which creates a new context and calls these registration methods.
5. Unit symbols are resolved via the Resolve() method, which:
    - Checks for special cases (dimensionless units, gram)
    - Looks up exact matches in base and derived unit maps
    - Handles prefixed units by checking if a symbol starts with a known prefix
    - Returns a Unit with the appropriate value and dimension
6. Complex unit expressions like "kgm/s^2" are parsed using an AST-based parser:
    - ParseComplexUnit() tokenizes the input and builds an abstract syntax tree
    - The parser handles identifiers, numbers, parentheses, and operations (multiplication, division, powers)
    - The AST nodes are then evaluated with the context to produce a final Unit
