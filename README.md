# Mutex

Mutex is a dynamically-typed scripting language implemented in Go. It features a tree-walk interpreter with support for mutable and immutable variables, control flow constructs, and lexical scoping.

## Installation

### Prerequisites
- Go 1.25.3 or higher

### Installing
```bash
go install github.com/caelondev/mutex@latest
```

### Running
```bash
mutex              # Start interactive REPL
mutex <filepath>   # Execute a Mutex source file
```

## Language Reference

### Types

Mutex supports the following data types:

```mutex
// Number - IEEE 754 double-precision floating-point
42
3.14159
// even this
-69
-67.41

// String - UTF-8 encoded text
"Hello, World!"

// Boolean - logical values
true
false

// Nil - represents absence of value
nil

// Arrays - represents a list of data compressed into one
["Mutex", 0, 1, "Hello, World", true, 3.14159]
```

### Variable Declarations

Variables in Mutex can be declared as either mutable or immutable:

**Mutable variables** can be reassigned after declaration:
```mutex
var mut counter = 0;
counter = counter + 1;

// and even this
counter += 1
```

**Immutable variables** cannot be reassigned:
```mutex
var imm max_retries = 3;
// max_retries = 5;  // Error: cannot reassign immutable variable
```

### Arrays

**Mutex arrays** can be instantiated with the syntax
```mutex
[0, 1, 2, ...]
```

Arrays are objects, they're a reference in the memory. hence they cannot be directly compared...
```mutex
[] == [] // INVALID!
```

#### Indexing

If you want to compare the data inside an array, you need to index the values
```mutex
var mut fruits = ["apple", "orange", "banana"]

fruits[0] // apple
fruits[1] // orange
fruits[2] // banana
```

You can also do nested arrays
```mutex
var mut random = [[0, "Hello"], ["Apple", 50], "World"]

random[0][1] // "Hello"
random[2] // "World"
```

#### Re-assigning

Re-assigning array values is as easy as reassigning a variable value
```mutex
var mut placeholders = ["foo", "bar"]

placeholders[0] = "baz" // This replaces "foo" with "baz"
```

### Expressions

#### Arithmetic Operators
```mutex
10 + 5    // Addition: 15
10 - 5    // Subtraction: 5
10 * 5    // Multiplication: 50
10 / 5    // Division: 2
10 % 3    // Modulo: 1
```

#### Compound Assignment Operators
```mutex
x += 5    // Shorthand for `x = x + 5`
x -= 5    // Shorthand for `x = x - 5`
x *= 5    // Shorthand for `x = x * 5`
x /= 5    // Shorthand for `x = x / 5`
x %= 3    // Shorthand for `x = x % 5`
```

#### Increment/Decrement Operators
```mutex
x++ // Shorthand for `x = x + 1`
x-- // Shorthand for `x = x - 1`
```

#### Comparison Operators
```mutex
5 < 10     // Less than
5 <= 10    // Less than or equal
5 > 10     // Greater than
5 >= 10    // Greater than or equal
5 == 10    // Equal to
5 != 10    // Not equal to
```

#### Logical Operators
```mutex
true and true // true
false or true // true
not true // false

not (true and true) //false
```

### Control Flow

#### Conditional Statements

Execute code based on boolean conditions:

```mutex
if (temperature > 30) {
    // Hot weather
} else if (temperature > 20) {
    // Mild weather
} else {
    // Cold weather
}
```

Parentheses around conditions are optional:
```mutex
if temperature > 30 {
    // Also valid
}
```



#### While Loops

Execute code repeatedly while a condition holds true:

```mutex
var mut i = 0;
while (i < 5 and i >= 0) {
    i = i + 1;
}
// i is now 5
```

#### For Loops

Execute code with explicit initialization, condition, and increment:

```mutex
for (var mut i = 0; i < 10; i++) {
    // Loop body executes 10 times
}
// i is not accessible here (scoped to loop)
```

### Scoping Rules

Mutex uses lexical scoping:

- Variables declared outside blocks are accessible within nested blocks
- Variables declared within blocks are not accessible outside those blocks
- For loop initializers are scoped to the loop body
- While loop bodies create new scopes, but the loop itself does not

```mutex
var mut x = 10;

if (x > 5) {
    var mut y = 20;  // y is scoped to this block
    x += y;          // x is accessible from outer scope
}

// x is 30
// y is not accessible here
```

## Examples

### Sum of Natural Numbers
```mutex
var mut sum = 0;
for (var mut i = 1; i <= 100; i++) {
    sum += i;
}

// sum is now 5050
```

### Factorial Calculation
```mutex
var mut n = 5;
var mut factorial = 1;
var mut i = 1;

while (i <= n) {
    factorial *= i;
    i += 1;
}

// factorial is now 120
```

### FizzBuzz
```mutex
for (var mut i = 1; i <= 15; i += 1) {
    var mut output = "";
    
    if (i % 3 == 0) {
        output = "Fizz";
    }
    
    if (i % 5 == 0) {
        output = output + "Buzz";
    }
}
```

## Implementation

Mutex is implemented as a tree-walk interpreter consisting of:

- **Lexer** - Tokenizes source code into lexical tokens
- **Parser** - Constructs an abstract syntax tree (AST) using Pratt parsing
- **Evaluator** - Traverses and executes the AST with environment-based variable storage

The interpreter supports proper lexical scoping through environment chaining, where each scope maintains a reference to its parent scope.

## Author

Built by [caelondev](https://github.com/caelondev)

## License

MIT
