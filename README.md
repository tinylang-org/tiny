<p align="center">
<h2 align="center"> Tiny programming language </h2>
</p>

**Tiny** - **tiny** programming language project for **tiny** embbed projects with **tiny** easy-to-learn syntax, **tiny** compiler and **tiny** package manager.

# Philosophy
Philosophy of the Tiny programming language is as follows: *do not overcomplicate the things*.

This is a problem of programming languages like C++ and Rust. They have huge standart library and huge amount of different techniques on writing an efficient code. Because of this even proffesional programmers that write on this languages cannot own some concepts of these languages even after years of development.

That does **not** mean that the language must be as simple as possible. Great example in this case is Go. Go because of its simplicity loses the scalability of large projects written on it, which is also not right. It has strange OOP and bad error handling.

Tiny programming language is about to solve the problem of overcomplication in programming languages. So it is basically "C with classes".

Here is a code that demonstrates OOP in Tiny:

```
pub struct Account {
   pub readonly age: int;
   pub readonly name: string;
   password: string;

   pub init(age: int, name: string, password: string) {
    this.age = age;
    this.name = name;
    this.password = password;
   }

   pub destroy() {
    destroy this.name; // analog of free
    destroy this.password;
   }
}

pub fun main() {
    var me = new Account(14, "Adi", "jifsjdif");
    destroy me;
}
```

In other words: Tiny is as easy as Golang, but more scalable.

# Design
Tiny language specification is going to move [here](https://github.com/tinylang-org/tiny-spec). 

# Progress

- [x] Lexer
- [ ] Parser
- [ ] LLVM frontend
- [ ] Generics
- [ ] Package manager (tpm)
- [ ] Using external C APIs
- [ ] Dynamic linking