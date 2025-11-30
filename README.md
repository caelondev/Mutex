# Mutex

A tree-walk interpreter built in Go, following [Crafting Interpreters](https://craftinginterpreters.com).

Mutex is a dynamically-typed scripting language built as a learning project to understand how programming languages work from the ground up.

# SYNTAX

### DATATYPES
```Mutex
// Strings ("")
// Numbers (10, 10.5)
// Booleans (true, false)
// Nil (nil)
```

### VARIABLE DECLARATION
Mutable variables (Variables that can be reassigned)
```Mutex
var mut foo = 10;
```

Immutable variables (Variables that cannot be reassiyned)
```Mutex
var imm bar = "Hello, World!";
```

### VARIABLE RE-ASSIGN
```mutex
var mut foo = "bar"; // MUST BE `mut` TO MUTATE/REASSIGN!

foo = "buzz";
```

## Credits

- Built by [caelondev](https://github.com/caelondev)
- Based on [Crafting Interpreters](https://craftinginterpreters.com) by Robert Nystrom
