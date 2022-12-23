# Comments
Comments serve as a program documentation. They can be:
1. *One line comments*, start with "//" and stop at [EOL](../terminology_used.md#eol) or [EOF](../terminology_used.md#eof). Example:
```tiny
// If user logged in
var loggedIn: bool = false
```
2. *Multi line comments*, start with "/*" and end at "*/". If there will be [EOF](../terminology_used.md#eof) appeared in the [scanning process](../terminology_used.md#lexer) when processing multi line comment, error will appear. Example of multi line comments usage:
```tiny
/*  Function that calculates the factorial

    It uses recursion, that is why it is not recommended to pass big values. Stack overflow
    is bad ._. */
pub fun factorial(n: int64): int64 {
    if n < 2 {
        return 1
    }
    
    return factorial(n - 1) * n
}
```