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