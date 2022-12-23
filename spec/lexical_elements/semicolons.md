# Semicolons in 0.1a
Semicolons in Tiny 0.1a play role of delimiter of statements in source code. They are **not** autoinserted by the compiler. And they are required to be placed at the end of each statement (except for function and struct definition):
```tiny
struct Person {
    pub name: string; // <- semicolon is required
    pub age: int32; // <- semicolon is required
} // <- no semicolon required here
```
```tiny
import "sub_lib.tl"; // <- semicolon required

pub fun sum(a: int32, b: int32): int32 {
    return a + b; // <- semicolon is required
} // <- no semicolon required
```

# Semicolons now
Semicolons in Tiny 0.1a+ are not required and will be autoinserted by the compiler:
```tiny
import "sub_lib.tl"

pub fun sum(a: int32, b: int32): int32 {
    return a + b
}
```
