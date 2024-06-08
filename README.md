# Flux 

### Compiler and Interpreter for general-purpose programming language with switch - matching patterns instruction.

Implementation of all elements in the project written in **Golang**,
including **Reader**, **Lexer**, **Parser** and **Interpreter**.

1. [Description of functionality](#description-of-functionality)
2. [Allowed data types](#allowed-data-types)
3. [Language-assumptions](#language-assumptions)
4. [Characteristic operations](#characteristic-operations)
5. [Built-in functions](#built-in-functions)
6. [EBNF specification and syntax](#ebnf-specification-and-syntax)
7. [Examples of allowable constructions and semantics](#examples-of-allowable-constructs-and-semantics)
8. [Error-handling and examples](#error-handling-and-examples)
10. [Install](#install)
11. [Type conversion](#type-conversion-and-type-combination-acceptable-for-multi-argument-operators-and-built-in-functions)
12. [Rules for passing variables to functions](#rules-for-passing-variables-to-functions)
13. [Module implementation](#module-implementation)
14. [Tests](#tests)

## Description of functionality

The “Flux” language design allows for:

- initialization and assignment of variables,
- performing arithmetic operations,
- support for conditional statements and loops,
- defining functions with or without arguments,
- type conversion using the `as` operator,
- calling functions,
- recursive functions,
- support for relational patterns in the switch statement.

---

## Allowed data types

- integer(int)
- float
- string
- boolean (bool)

---

## Assumptions of the language

- statically typed language,
- arguments passed by value,
- error handling at the lexical, syntactic and semantic levels,
- the value of the variable can be changed, but its type must match,
- every program written in the "flux" language must have a `main()` function, which contains the main body of the program and its operation begins with this function,
- an external variable can be covered by a variable with the same name located in a function block, loop, conditional statement or switch statement.
- defined functions cannot be overridden by functions with other arguments,
- built-in functions cannot be overridden

---

## Characteristic operations

The characteristic operation is relational pattern matching

- The `switch` instruction works on comparing variables, it can be invoked by declaring variables local to the expression, or using variables defined in higher `scopes`,
- There can be more than one variable declaration

Example call:

- with the declaration of a local variable for the expression:

```golang
switch int a := 2 + 2 {
    a == 4 => print("four"),
    default => print("definitely not four")
}
```

- using previously defined variables

```golang
int a := 3
switch {
    a == 4 => print("four"),
    default => print("definitely not four")
}
```

Behavior of the **switch** statement:

- the switch statement because the right side of the `=>` operator may have an expression or block opened with `{` and closed with `}`,
- the instruction loops through the current cases and calls the instructions for the first positively evaluated case,
- if there is a 'block` on the right, the switch statement does not return a value, but evaluates the defined block (a return can be defined in it),
- if we want the instruction to return a value, to the right of `=>` we should place an expression with the type of this value, e.g.:

- instruction returning the `int` value:

```golang
switch {
    default => 2
}
```

- instruction returning the `bool` value:

```golang
switch {
    default => true
}
```

- instruction not returning a value:

```golang
switch {
    default => { print("flux") }
}
```

---

## Built-in functions

The following functions are built into the language:

- `print(...)` - function that prints the passed values ​​to Stdout, works for any number of arguers,
- `println(...)` - a function similar to the 'print' function, which additionally separates each of the passed arguments with a newline
- `sqrt(var1, var2 float) -> float` - function returning the square root, accepts `float` type arguments, returns `float` type argument,
- `power(var1, var2 float) -> float` - a function that returns the number given as the first argument of the `float` type, raised to the power given as the second argument of the `float` type, returns a `float` value,

---

## EBNF specification and syntax

- terminal symbols marked with `*`

```go
program = { function_definition } ;

function_definition = identifier , "(", [ parameters ], ")", [ type_annotation ] , block ;

parameters = parameter_group , { "," , parameter_group } ;
parameter_group = identifier , { ",", identifier }, type_annotation ;

type_annotation = "int" | "float
" | "bool" | "str" ​​;

block = "{" , { statement } , "}" ;

statement = variable_declaration
 | assignment_or_call
 | conditional_statement
 | loop_statement
 | switch_statement
 | return_statement
 ;

variable_declaration = type_annotation, identifier, ":=", expression ;

assignment_or_call = identifier, ( "(", [ arguments ], ")" ) | ( "=", expression ) ;

conditional_statement = "if" , expression , block , [ "else" , block ] ;

loop_statement = "while" , expression, block ;

switch_statement = "switch", [( variable_declaration, { ",", variable_declaraion } ) ], "{", switch_case, { ",", switch_case "}" ;

switch_case = ( expression | "default" ), "=>", ( expression | block ) };

return_statement = "return" , [ expression ] ;



expression = conjunction_term, { "or", conjunction_term };

conjunction_term = relation_term, { "and", relation_term } ;

relation_term = additive_term, [ relation_operator, additive_term ] ;

relation_operator = ">="
 | ">"
 | "<="
 | "<"
 | "=="
 | "!="
 ;

additive_term = multiplicative_term, { ("+" | "-"), multiplicative_term } ;

multiplicative_term = casted_term, { ("*" | "/"), casted_term } ;

casted_term = unary_operator, [ "as", type_annotation ] ;

unary_operator = [ ("-" | "!") ], term ;

term = integer
 | float
 | bool
 | string
 | identifier_or_call
 | "(" , expression , ")"
 ;

identifier_or_call = identifier, [ "(", [ argumets ], ")" ] ;

arguments = expression , { "," , expression };

identifier = letter , { letter | digit | "_" } ;

float = integer , "." , digit , { digit } ;

*integer = "0" | positive_digit , { digit } ;

*string = '"', { literal }, '"' ;

*literal = letter
 | digit
 | symbols
 ;

*bool = "true" | "false" ;

*letter = "a" | "..." | "with" | "A" | "..." | "WITH" ;

*positive_digit = "1" | "2" | "3" | "4"| "5" | "6"| "7" | "8" | "9" ;

*digit = "0" | "1" | "2" | "3" | "4"| "5" | "6"| "7" | "8" | "9" ;

*symbols = "`" | "~" | "!" | "@" | "#" | "$" | "%" | "^" | "&" | "*" | "(" | ")" | "_" | "-" | "+" | "=" | "{" | "}" | "[" | "]" | ";" | ":" | "'" | "," | "." | "?" | "/" | "|" | "\" ;
```

---

## Examples of allowable constructions and semantics

Initialization and value assignment

```go
int a := 5
int b := 2

a = 8
```

---

Arithmetic operations

```go
int a := 3
a = a + 3 * (2 - 1)
```

---

Comments

```go
# This is a comment
```

---

Conditional statement

```go
if y > 5 {
    print("Y is greater than 5")
 } else {
    print("Y is less than or equal to 5")
 }
}

string name := "Ala has a dog"
if name == "Ala has a cat" {
    print("The cat belongs to Ania")
} else {
    print("It's neither Ala nor the cat")
}
```

---

While loop statement

```go
int num := 10
while num > 0 {
 print(num)
 num = num - i
}
```

---

Function with argument

```go
circleArea(r int) float    {
 return 3.14 * (r * r)
}

main(){
 int r := 2
 int a := circleArea(r)
 print(a)
}
# output: 12.56636
```

---

Recursive function

```go
fibonacci(n) int {
 if n <= 1 {
    return n
 } else {
    return fibonacci(n - 1) + fibonacci(n - 2)
 }
}

main(){
    print(fibonacci(3))
}

# output: 2
```

---

Type conversion

```go
int a := 5
string c := a as string
print(c) # "5"

int b := 0
bool d = b as bool # "false"
```

---

Built-in features

```go
print(sqrt(9 as float))
# output: 3

```

```go
print(power(3 as float, 2.0))
# output: 9
```

---

Relational patterns - switch instruction

```go
sumUp(a,b int) int {
 return a + b
}

whatWillGetMe(a,b int) string {
 switch int c := sumUp(a, b) {
 c>2 and c<=4 => "A pint",
 c==5         => "Decent beverage",
 c>5 and c<15 => "A NICE bevrage",
 c>15         => "Whole bottle",
 default      => "Nothing today!"
 }
}
main(){
 print(whatWillGetMe(2,3))
}
# output: Decent beverage
```

```go
giveMeWord() string {
 return "word"
}

nameNumber() int {
 string c := giveMeWord()
 switch {
 c == "Sammy" => 0,
 c == "World" => 1,
 c == "word" => 2,
 default => 3
 }
}

main(){
 print(nameNumber())
}

# ourput: 2
```

```go
getUserRole(userId int) string {
 return "admin"
}

checkPermission(role, permission string) bool {
 return role == "admin" and permission == "edit"
}

main() {
 int userId := 123

 switch string userRole := getUserRole(userId) {
    userRole == "admin" => {
        if checkPermission(userRole, "edit") {
            print("The user has edit permissions")
        } else {
            print("User does not have edit permissions")
        }
    },
    userRole == "user" => {
        print("User has limited permissions")
    },
    default => {
        print("Unknown user role")
    }
 }
}
# output: The user has edit permissions
```

## Error handling and examples

Error handling takes place at all levels, i.e.:

- lexer,
- parser,
- interpreter

Due to the use of the `panic()` method in golang, program processing is interrupted when the first error is encountered. The 'errorHandlers' function defined for the lexer and parser are responsible for 'catching' the error; the panic triggered in the interpreter is caught in the main.go function, which forwards the error content to Stdout.
Initially, the implementation was to be carried out with error propagation, after which the concept was changed in consultation with the host. A change is planned in the future, moving from the `panic()` function to passing errors via **error** values.
Each module has defined constants with error messages that contain the content and the place of occurrence.

---

Error format:

```go
`error [<line> : <column>]: <message>`
```

---

Unclosed string:

```go
string a := "this is a string
```

```go
error [1, 27] String not closed, perhaps you forgot "
```

---

Going beyond the int value limit

```go
main(){
 int a := 99999999999999...
}
```

```go
error [2, 14]: Int value limit Exceeded
```

---

Assignment error:

```go
main(){
 int a := 5
 a = "Ala has a cat"
}
```

```go
error [3, 3]: type mismatch: expected int, got string
```

---

Different type return error:

```go
sumUp(a, b int) float {
 return a + b
}

main(){
 print(sumUp(20, 11))
}
```

```go
error [2, 12]: invalid return type: int, expected: float
```

---

Error in switch design:

```go
kelvinToCelcius(temp int) int {
 return temp - 273
}

howCold(kelvin int) string {
 switch int c := kelvinToCelcius(kelvin) {
    c < -20 => "Freezing",
    c>0 and c<10 => "Chilling",
    c>=10 and c<20 => "Warm"
 }
}

main(){
 print(howCold(300))
}
```

```go
error [6, 1]: missing return, function should return type: string
```

---

Uninitialized variable:

```go
main(){
 print(a + 10)
}
```

```go
error [2, 9]: undefine: a
```

---

Undeclared function error:

```go
main() {
 print(unknownFunction())
}
```

```go
error [2, 9]: undefined function: unknownFunction
```

---

Error invalid number of function arguments:

```go
add(a, b int) int {
 return a + b
}
main() {
 print(add(5))
}
```

```go
error [5, 9]: function add expects 2 arguemnts but got: 1
```

---

Type mismatch error in conditional statement:

```go
main() {
 int a := 5
 if a == "test" {
 print("Equal")
 }
}
```

```go
error [3, 8]: cannot evaluate '==' operation with instances, mismatched types of int and string
```

---

Incorrect operator usage error:

```go
main() {
 int a := 20
 string b := "5"
 int c := a / b
}
```

```go
error [4, 14]: cannot evaluate '/' operation with instances of int and string
```

---

Incorrect use of relational operator error:

```go
main() {
 int a := 5
 if a < "test" {
 print("Less than")
 }
}
```

```go
error [3, 8]: cannot evaluate '<=' operation with instances, mismatched types of int and string
```

---

Division by zero error:

```go
main() {
 int a := 10
 int b := 0
 result := a / b
}
```

```go
error [4, 19]: Division by zero
```

## Install

To run a program written in **flux** you should:

- you need to have [Golang] compiler(https://go.dev/dl/)
- clone this repository and **cd** into it
- build the project `$ go build -o flux .`
- move the binary `$ sudo mv flux /usr/local/bin/` or run the program via `./flux`

---

## Input - streams/files and interpreter startup

A program written in Flux can be run from both a file and an input data stream.
Files should have the extension `.fl`.

The standard way to run a written program is to invoke the compiler with an argument specifying the path to the file:

```shell
flux example.fl
```

If the program accepts initial arguments, they should be given after the specified file:

```shell
flux example.fl 0 1
```

The program code can be passed from standard input using the `|` operator and typing the first argument as `-`:

```shell
echo 'main(){ print("hellooo") }' | flux -
```

or

```shell
flux - < example.fl
```

Calling the program from standard input with arguments:

```shell
echo 'main(a int){ print(a) }' | flux - 2
```

or

```shell
flux - < example.fl 0 2
```

The Flux language does not require any special configuration data to function properly.

The program interpreter gains access to standard output and input, which allows it to capture program results, show errors, and provide input data to the program.

---

## Type conversion and type combination acceptable for multi-argument operators and built-in functions

Because the language is strongly and statically typed, any type conversion is explicit and the `as` operator is available to perform it.

Type conversion for static typing:

| With    | To Integer | To Float | To String | To Boolean |
| ------- | ---------- | -------- | --------- | ---------- |
| Integer | -          | Explicit | Explicit  | Explicit   |
| Float   | Explicit   | -        | Explicit  | Explicit   |
| String  | Explicit   | Explicit | -         | Explicit   |
| Boolean | Explicit   | Explicit | Explicit  | -          |

For int to boolean:

- int 0 means `false`
- other than 0 means `true`

For string to boolean:

- empty string: "" means `false`
- a non-empty string means `true`

For float to boolean:

- float 0.0 means `false`
- other than 0 means `true`

**Operations `*`, `/`, `+`, `-`:**
Multiplication (`*`)

- **int \* int**: Returns the result as an integer value (`int`).
- **float \* float**: Returns the result as a floating point number (`float`).
- **int \* float** and **float \* int**: Returns the result as a floating point number (`float`).

Division (`/`)

- **int / int**: Returns the result as an integer (`int`).
- **float / float**: Returns the result as a floating point number (`float`).
- **int / float** and **float / int**: Returns the result as a floating point number (`float`).

Adding (`+`)

- **int + int**: Returns the result as an integer value (`int`).
- **float + float**: Returns the result as a floating point number (`float`).
- **int + float** and **float + int**: Returns the result as a floating point number (`float`).
- **string + string**: Concatenate strings.
- **int + string**, **float + string**, **string + int** and **string + float**: Concatenate a number or floating-point value with a string, returns (`string`).

Subtraction (`-`)

- **int - int**: Returns the result as an integer value (`int`).
- **float - float**: Returns the result as a floating point number (`float`).
- **int - float** and **float - int**: Returns the result as a floating point number (`float`).

## Rules for passing variables to functions

Variables are passed to the function by value. As there are no structures, passing a variable by reference does not seem to be necessary.

---

## Function overloading

Function overloading is not allowed, there cannot be two functions with the same name.

Built-in functions also cannot be overridden.

---

## Implementation of modules

1. **Lexical analyzer** (lexer):

- Processes the source code character by character, and according to the grammar produces tokens to identify and group lexemes such as identifiers, numbers, operators and keywords.
- Tokens store information about their location in the source code in the form `(line no., column no.)`.
- If a string of characters is encountered that is impossible to decode, the analyzer scans the string until it finds white character and returns the `UNDEFIND` token

2. **Syntactic parser** (parser):

- The parser takes as input a stream of tokens produced by the lexical parser.
- The task of the parser is to produce a parsing tree of the program in the form of `nodes'.
- Strictly waits for expected token when parsing expression.
- Syntax error handling implemented via `panic()`, containing information about the location of the incorrect expression in the program code.

3. **Interpreter**:

- Operates on a syntactic parsing tree.
- Written using the "Visitor" design pattern.
- The interpreter visits the elements of the syntax tree, evaluating their contents. Assigns values ​​to variables, checks type compatibility, compliance of arguments supplied to calls, runs called functions (including built-in functions).
- Makes sure that recursive calls do not exceed the defined limit (implementation using CallStack).
- Performs arithmetic operations, supports conditional statements, loops, function calls and other language constructs.

---

