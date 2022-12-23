# Top level statements

> ```ebnf
> program_unit = { import_statement } { tl_statement } |
>                namespace_declaration { import_statement } { tl_statement } .
> ```

Every tiny source file is called "program unit". Program unit can be represented using Abstract Syntax Trees.

Program unit AST consists of namespace declaration, top level statement (but not statement) and import statement AST nodes. Example:

```tiny
namespace main // ok (namespace declaration)

import "test.tl" // ok (import statement)

return 0 // syntax error (unexpected statement node)
```

Import statements go first:
```tiny
namespace main // ok (namespace declaration)

pub fun main() {
	// some code here
}

import "lib.tl" // syntax err (unexpected import statement)
``` 