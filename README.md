<p align="center">
<h2 align="center"> Tiny programming language </h2>
</p>

**Tiny** - **tiny** programming language project for **tiny** embbed projects with **tiny** easy-to-learn syntax, **tiny** compiler and **tiny** package manager.

```mermaid
flowchart LR;
    A[Tiny source file] -->|Parser| B(Abstract Syntax Tree);
    B -->|TinyC C translator| D[C code]
    D -->|Clang frontend| E[LLVM IR]
    G[External C library, static linking] -->|Clang frontend| E[LLVM IR]
    E[LLVM IR] -->|LLVM backend| I[Object]
    L[External C library objects] -->|lld| P[Native binary]
    I -->|lld| P[native binary]
```
