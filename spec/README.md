<p align="center">
<image width="70%" src="./tiny.png"> <br>
<h2 align="center"> Tiny programming language </h2>
</p>


**Tiny** - **tiny** programming language project for **tiny** embbed projects with **tiny** easy-to-learn syntax, **tiny** compiler and **tiny** package manager.

## Why Tiny?
This language adheres to the philosophy: make it fast and simple. 

Imagine you need to create an embedded application. You need to do this fairly quickly, but speed is critical in your project. In this case, you can choose the C++ programming language, but in this case, your project will be cumbersome, it will become difficult to include external dependencies, and learning C++ itself is rather problematic.

You can choose the Rust programming language, which has a borrow checker, a large community, and is pretty fast. However, as in the case of C++, learning this language is very problematic, because it opens up completely new concepts for an already mastered programmer. Also dealing with the borrow checker can be really painfull experience for some people. 

There is a question about such things as the memory management system.

Since the philosophy of the language is to run programs quickly, and garbage collection sometimes takes an unthinkable amount of resources or time (depends on the language), this way of managing memory is not suitable for the language. 

See: https://discord.com/blog/why-discord-is-switching-from-go-to-rust#:~:text=Discord%20is%20a%20product%20focused,and%20messages%20you%20have%20read.

Controlling memory by your own doesn't really fit the philosophy of the language either - it's not nice and it's not **simple**. Therefore, C/C++ would not be very suitable for a simple project.

As it was already said, borrow checker also isn't suitable for the goal of **simplicity**. That is why I decided to create my own memory management system.

## Tiny-way
```
pub fun main() {
    println("hello world");
}
```