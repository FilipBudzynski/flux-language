# TKOM - Dokumentacja Końcowa

### Flux - własny język programowania ogólnego przeznaczenia z instrukcją relacyjnych wzorców (switch - matching patterns).

Filip Budzyński, numer albumu: 319021

1. [Opis funkcjonalności](#opis-funkcjonalności)
2. [Dopuszczane typy danych](#dopuszczane-typy-danych)
3. [Przyjęte założenia języka](#przyjęte-założenia-języka)
4. [Charakterystyczne operacje](#charakterystyczne-operacje)
5. [Funkcje wbudowane](#funkcje-wbudowane)
6. [Specyfikacja i składnia EBNF](#specyfikacja-i-składnia-ebnf)
7. [Przykłady dopuszczalnych konstrukcji i semantyka](#przykłady-dopuszczalnych-konstrukcji-i-semantyka)
8. [Obsługa błędów i przykłady](#obsługa-błędów-i-przykłady)
9. [Rozróżniane Tokeny](#rozróżniane-tokeny)
10. [Uruchomienie](#uruchomienie)
11. [Konwersja typów](#konwersja-typów-i-kombinacja-typów-akceptowalna-dla-operatorów-wieloargumentowych-i-funkcji-wbudowanych)
12. [Zasady przekazywania zmiennych do funkcji](#zasady-przekazywania-zmiennych-do-funkcji)
13. [Realizacja modułów](#realizacja-modułów)
14. [Testy](#testy)

## Opis funkcjonalności

Implementacja wszystkich elementów projektu napisana w języku **Golang**,
w tym Reader, Lexer, Parser i Interpreter.

Projekt języka “Flux” umożliwia:

- inicjalizację i przypisywanie zmiennych,
- wykonywanie operacji arytmetycznych,
- obsługę instrukcji warunkowych i pętli,
- definiowanie funkcji z lub bez argumentów,
- konwersję typów przy użyciu operatora `as`,
- wywoływanie funkcji,
- funkcje rekurencyjne,
- obsługę relacyjnych wzorców w insturkcji switch.

---

## Dopuszczane typy danych

- integer (int)
- float (float)
- string (string)
- boolean (bool)

---

## Przyjęte założenia języka

- język statycznie typowany,
- argumenty przekazywane przez wartość,
- obsługa błędów na poziomie leksykalnym, syntaktycznym i semantycznym,
- wartość zmiennej może być zmieniona, ale musi zgadzać się jej typ,
- każdy program napisany w języku “flux” musi posiadać funkcję `main()`, która zawiera główne ciało programu i od tej funkcji rozpoczyna się jego działanie,
- zmienna zewnętrzna może zostać przykryta poprzez zmienną o tej samej nazwie znajdującej się w bloku funkcji, pętli, instrukcji warunkowej lub instrukcji switch.
- zdefiniowanych funkcji nie można przesłaniać funkcjami z innymi argumentami,
- funkcji wbudowanych nie można przesłaniać

---

## Charakterystyczne operacje

Operacją charakterystyczną jest dopasowanie wzorców relacyjnych

- Instrukcja `switch` działa na porównywaniu zmiennych, można ją wywołać deklarując zmienne lokalne dla wyrażenia, lub używając zmiennych zdefiniowanych w wyżyszch `scopa'ach`,
- Deklaracji zmiennych może być więcej niż jedna

Przykładowe wywołanie:

- z deklaracją zmiennej lokalnej dla wyrażenia:
  ```golang
  switch int a := 2 + 2 {
    a == 4  => print("cztery"),
    default => print("na pewno nie cztery")
  }
  ```
- z wykorzystaniem wcześniej zdefiowanych zmiennych
  ```golang
  int a := 3
  switch {
    a == 4  => print("cztery"),
    default => print("na pewno nie cztery")
  }
  ```

Zachowanie instrukcji **switch**:

- instukcja switch bo prawej stronie operatora `=>` może posiadać wyrażenie lub blok otwierany przy użyciu `{` i zamknięty przez `}`,
- instrukcja przechodzi przez napiasne przypadki i wywołuje instrukcje dla pierwszego pozytywnie zewaluowanego przypadku,
- jeżeli po prawej stronie znajduje się `block`, instrukcja switch nie zwraca wartości, lecz ewaluuje zdefiniowany blok (w nim możena zdefiniować return),
- jeżeli chcemy aby instrukcja zwróciła jakąś wartość, po prawej stronie `=>` należy umieścieć wyrażenie z typem tej wartości, np.:

  - instrukcja zwracająca wartość `int`:

  ```golang
  switch {
      default => 2
  }
  ```

  - instrukcja zwracająca wartość `bool`:

  ```golang
  switch {
      default => true
  }
  ```

  - instrukcja nie zwracająca wartość:

  ```golang
  switch {
      default => { print("flux") }
  }
  ```

---

## Funkcje wbudowane

W języku wbudowane są niżej wymienione funkcje:

- `print(...)` - funkcja drukująca przekazane wartości na Stdout, działa dla dowolnej liczby arguentów,
- `println(...)` - funkcja bliźniacza do funkcji 'print', która dodatkowo każdy z przekazanych argumentów odziela znakiem nowej lini
- `sqrt(var1, var2 float) -> float` - funkcja zwracająca pierwiastek kwadratowy, przyjmuje argumenty typu `float`, zwraca argument typu `float`,
- `power(var1, var2 float) -> float` - funkcja zwracająca liczbę podaną jako pierwszy argument typu `float`, podniesioną do potęgi podaną jako drugi argument typu `float`, zwraca wartość typu `float`,

---

## Specyfikacja i składnia EBNF

- symbole terminalne wyróżnione znakiem `*`

```go
program               = { function_definition } ;

function_definition   = identifier , "(", [ parameters ], ")", [ type_annotation ] , block ;

parameters            = parameter_group , { "," , parameter_group } ;
parameter_group       = identifier , { ",", identifier }, type_annotation ;

type_annotation       = "int" | "float" | "bool" | "str" ;

block                 = "{" , { statement } , "}" ;

statement             = variable_declaration
                      | assignment_or_call
                      | conditional_statement
                      | loop_statement
                      | switch_statement
                      | return_statement
                      ;

variable_declaration  = type_annotation, identifier, ":=", expression ;

assignment_or_call    = identifier,  ( "(", [ arguments ], ")" ) | ( "=", expression ) ;

conditional_statement = "if" , expression , block , [ "else" , block ] ;

loop_statement        = "while" , expression, block ;

switch_statement      = "switch", [( variable_declaration, { ",", variable_declaraion } ) ], "{", switch_case, { ",", switch_case "}" ;

switch_case           = ( expression | "default" ), "=>", ( expression | block ) } ;

return_statement      = "return" , [ expression ] ;



expression            = conjunction_term, { "or", conjunction_term } ;

conjunction_term      = relation_term, { "and", relation_term } ;

relation_term         = additive_term, [ relation_operator, additive_term ] ;

relation_operator     = ">="
                      | ">"
                      | "<="
                      | "<"
                      | "=="
                      | "!="
                      ;

additive_term         = multiplicative_term, { ("+" | "-"), multiplicative_term } ;

multiplicative_term   = casted_term, { ("*" | "/"), casted_term } ;

casted_term           = unary_operator, [ "as", type_annotation ] ;

unary_operator        = [ ("-" | "!") ], term ;

term                  = integer
                      | float
                      | bool
                      | string
                      | identifier_or_call
                      | "(" , expression , ")"
                      ;

identifier_or_call    = identifier, [ "(", [ argumets ], ")" ] ;

arguments             = expression , { "," , expression } ;

identifier            = letter , { letter | digit | "_" } ;

float                 = integer , "." , digit , { digit } ;

*integer              = "0" | positive_digit , { digit } ;

*string               = '"', { literal }, '"' ;

*literal              = letter
                      | digit
                      | symbols
                      ;

*bool                 = "true" | "false" ;

*letter               = "a" | "..." | "z" | "A" | "..." | "Z" ;

*positive_digit       = "1" | "2" | "3" | "4"| "5" | "6"| "7" | "8" | "9" ;

*digit                = "0" | "1" | "2" | "3" | "4"| "5" | "6"| "7" | "8" | "9" ;

*symbols              = "`"
                      | "~"
                      | "!"
                      | "@"
                      | "#"
                      | "$"
                      | "%"
                      | "^"
                      | "&"
                      | "*"
                      | "("
                      | ")"
                      | "_"
                      | "-"
                      | "+"
                      | "="
                      | "{"
                      | "}"
                      | "["
                      | "]"
                      | ";"
                      | ":"
                      | "'"
                      | ","
                      | "."
                      | "?"
                      | "/"
                      | "|"
                      | "\"
                      ;
```

---

## Przykłady dopuszczalnych konstrukcji i semantyka

Inicjalizacja i przypisanie wartości

```go
int a := 5
int b := 2

a = 8
```

---

Operacje arytmetyczne

```go
int a := 3
a = a + 3 * (2 - 1)
```

---

Komentarze

```go
# To jest komentrz
```

---

Instrukcja warunkowa

```go
if y > 5 {
    print("Y jest większe od 5")
    } else {
        print("Y jest mniejsze lub równe 5")
    }
}

string nazwa := "Ala ma psa"
if nazwa == "Ala ma kota" {
    print("Kot należy do Ani")
} else {
    print("Ani to Ala ani kot")
}
```

---

Instrukcja pętli while

```go
int num := 10
while num > 0 {
    print(num)
    num = num - i
}
```

---

Funkcja z argumentem

```go
circleArea(r int) float    {
    return 3.14 * (r * r)
}

main(){
    int r := 2
    float a := circleArea(r)
    print(a)
}
# output: 12.56636
```

---

Funkcja rekurencyjna

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

Konwersja typów

```go
int a := 5
string c := a as string
print(c)        # "5"

int b := 0
bool d := b as bool   # "false"
```

---

Funkcje wbudowane

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
        c == "Sammy"  => 0,
        c == "World"  => 1,
        c == "word"   => 2,
        default       => 3
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
                print("Użytkownik ma uprawnienia do edycji")
            } else {
                print("Użytkownik nie ma uprawnień do edycji")
            }
        },
        userRole == "user" => {
            print("Użytkownik ma ograniczone uprawnienia")
        },
        default => {
            print("Nieznana rola użytkownika")
        }
    }
}
# output: Użytkownik ma uprawnienia do edycji
```

## Obsługa błędów i przykłady

Obsługa błędów odbywa się na wszystkich poziomach, tj.:

- lexer,
- parser,
- interpreter

Ze względu na użycie metody `panic()` w golang'u, przetwarzanie programu jest przerywane po napotkaniu pierwszego błędu. Za 'łapanie' błędu są odpowiedzialne funkcję 'errorHandlers' zdefiniowane dla lexer'a i parser'a, panic wywoływany w interpreterze łapany jest w funkcji main.go, która przekazuje treść błędu na Stdout.
Początkowo realizacja miała odbyć się z propagacją błędu, po czym koncepcja została zmieniona przy uzgodnieniu z prowadzącym. W przyszłości planowana jest zmiana, przejście z funkcji `panic()` na przekazywanie błędów przez wartości **error**.
Każdy z modułów ma zdefiniowane stałe z wiadomościami o błędach, które zawierają treść, oraz miejsce wystąpienia.

---

Format błędów:

```go
`error [<line> : <column>]: <message>`
```

---

Niezamknięty string:

```go
string a := "to jest napis
```

```go
error [1, 27] String not closed, perhaps you forgot "
```

---

Wyjście poza limit wartości int

```go
main(){
    int a := 99999999999999...
}
```

```go
error [2, 14]: Int value limit Exceeded
```

---

Błąd przypisania:

```go
main(){
  int a := 5
  a = "Ala ma kota"
}
```

```go
error [3, 3]: type mismatch: expected int, got string
```

---

Błąd zwracania innego typu:

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

Błąd w konstrukcji switch, brak przypadku dla zakresu:

```go
kelvinToCelcius(temp int) int {
  return temp - 273
}

howCold(kelvin int) string {
  switch int c := kelvinToCelcius(kelvin) {
    c < -20       => "Freezing",
    c>0 and c<10   => "Chilling",
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

Niezainicjowana zmienna:

```go
main(){
  print(a + 10)
}
```

```go
error [2, 9]: undefind: a
```

---

Błąd niezadeklarowanej funkcji:

```go
main() {
  print(unknownFunction())
}
```

```go
error [2, 9]: undefined function: unknownFunction
```

---

Błąd nieprawidłowa liczba argumentów funkcji:

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

Błąd niezgodności typów w instrukcji warunkowej:

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

Błąd niepoprawnego użycia operatora:

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

Błąd niepoprawnego użycia operatora relacyjnego:

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

Błąd dzielenia przez zero:

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

## Rozróżniane Tokeny

---

Struktura tokenu

- TokeType
- Position
- Value

1. Zmienne:
   - **`IDENTIFIER`**
2. Operatory arytmetyczne:
   - **`PLUS`**
   - **`MINUS`**
   - **`MULTIPLY`**
   - **`DIVIDE`**
3. Operatory relacyjne:
   - **`EQUALS`**
   - **`NOT_EQUALS`**
   - **`GREATER_THAN`**
   - **`LESS_THAN`**
   - **`GREATER_THAN_OR_EQUAL`**
   - **`LESS_THAN_OR_EQUAL`**
4. Operatory logiczne:
   - **`AND`**
   - **`OR`**
   - **`NEGATE`** ( “-” lub “!”)
5. Słowa kluczowe:
   - **`IF`**
   - **`ELSE`**
   - **`WHILE`**
   - **`SWITCH`**
   - **`DEFAULT`**
   - **`AS`**
   - **`RETURN`**
6. Symbole specjalne:
   - **`ASSIGN`** (**`:=`**)
   - **`CASE_ARROW`** (**`=>`**)
   - **`LEFT_BRACE`** (**`{`**)
   - **`RIGHT_BRACE`** (**`}`**)
   - **`LEFT_PARENTHESIS`** (**`(`**)
   - **`RIGHT_PARENTHESIS`** (**`)`**)
   - **`COMMA`** (**`,`**)
7. Typy danych:
   - **`INT`**
   - **`FLOAT`**
   - **`STRING`**
   - **`BOOL`**
8. Wartości stałe:
   - **`CONST_INT`**
   - **`CONST_FLOAT`**
   - **`CONST_STRING`**
   - **`CONST_BOOL`**
   - **`CONST_TRUE`**
   - **`CONST_FALSE`**
9. Inne:
   - **`ETX`** (end of text)
   - **`UNDEFINED`**
   - **`COMMENT`**

---

## Uruchomienie

Aby uruchomić program napisany w języku **flux** należy:

- posiadać kompilator [Golang](https://go.dev/dl/)
- sklonować to repozytorium i przejść do niego
- zbudować projekt `$ go build -o flux .`
- przenieść binary `$ sudo mv flux /usr/local/bin/` lub uruchamiać program poprzez `./flux`

---

## Dane wejściowe - strumienie/pliki i uruchomienie interpretera

Program napisany w języku Flux może być uruchamiany zarówno z pliku, jak i ze strumienia danych wejściowych.
Pliki powinny mieć rozszerzenie `.fl`.

Standardowym sposobem uruchomienia napisanego programu jest wywołanie kompilatora z argumentem, podającym ścieżkę do pliku:

```shell
flux example.fl
```

Jeżeli program przyjmuje argumenty początkowe, należy je podać po wskazanym pliku:

```shell
flux example.fl 0 1
```

Kod programu może zostać przekazany ze standardowego wejścia używając operatora `|` oraz wpisując pierwszy argument jako `-`:

```shell
echo 'main(){ print("halooo") }' | flux -
```

lub

```shell
flux - < example.fl
```

Wywołanie programu ze standardowego wejścia z argumentami:

```shell
echo 'main(a int){ print(a) }' | flux - 2
```

lub

```shell
flux - < example.fl 0 2
```

Język Flux nie wymaga żadnych specjalnych danych konfiguracyjnych do poprawnego działania.

Interpreter programu dostaje dostęp do standardowego wyjścia i wejścia, co pozwala na przechwytywanie wyników działania programu i pokazywanie błędów oraz na podanie danych wejściowych do programu.

---

## Konwersja typów i kombinacja typów akceptowalna dla operatorów wieloargumentowych i funkcji wbudowanych

Ponieważ język jest silnie i statycznie typowany, każda konwersja typu jest jawna a do jej dokonania dostępny jest operator `as`.

Konwersja typów dla typowania statycznego:

| Z       | Do Integer | Do Float | Do String | Do Boolean |
| ------- | ---------- | -------- | --------- | ---------- |
| Integer | -          | Explicit | Explicit  | Explicit   |
| Float   | Explicit   | -        | Explicit  | Explicit   |
| String  | Explicit   | Explicit | -         | Explicit   |
| Boolean | Explicit   | Explicit | Explicit  | -          |

W przypadku int na boolean:

- int 0 oznacza `false`
- inne poza 0 oznaczają `true`

W przypadku string na boolean:

- pusty string: “” oznacza `false`
- niepusty string oznacza `true`

W przypadku float na boolean:

- float 0.0 oznacza `false`
- inne poza 0 oznacza `true`

**Operacje `*`, `/`, `+`, `-`:**
Mnożenie (`*`)

- **int \* int**: Zwraca wynik jako wartość całkowitą (`int`).
- **float \* float**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).
- **int \* float** oraz **float \* int**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).

Dzielenie (`/`)

- **int / int**: Zwraca wynik jako liczbę całkowitą (`int`).
- **float / float**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).
- **int / float** oraz **float / int**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).

Dodawanie (`+`)

- **int + int**: Zwraca wynik jako wartość całkowitą (`int`).
- **float + float**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).
- **int + float** oraz **float + int**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).
- **string + string**: Konkatenacja ciągów znaków.
- **int + string**, **float + string**, **string + int** oraz **string + float**: Konkatenacja liczby lub wartości zmiennoprzecinkowej z ciągiem znaków, zwraca (`string`).

Odejmowanie (`-`)

- **int - int**: Zwraca wynik jako wartość całkowitą (`int`).
- **float - float**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).
- **int - float** oraz **float - int**: Zwraca wynik jako liczbę zmiennoprzecinkową (`float`).

## Zasady przekazywania zmiennych do funkcji

Zmienne są przekazywane do funkcji przez wartość. Jako, że nie ma struktur to przekazywanie zmiennej przez referencje nie wydaje się być konieczne.

---

## Przeciążanie funkcji

Przeciążanie funkcji nie jest dozwolone, nie mogą istnieć dwie funkcje o takiej samej nazwie.

Funkcje wbudowane również nie mogą być przesłaniane.

---

## Realizacja modułów

1. **Analizator leksykalny** (lexer):
   - Przetwarza kod źródłowy, znak po znaku, i zgodnie z gramatyką produkuje tokeny do identyfikacji i grupowania leksemów, takich jak identyfikatory, liczby, operatory i słowa kluczowe.
   - Tokeny przechowują informacje o swoim położeniu w kodzie źródłowym w postaci `(nr linii, nr kolumny)`.
   - W przypadku natrafienia na niemożliwy do zdekodowania ciąg znaków, analizator skanuje ciąg aż do natrafienia na biały znak i zwraca token `UNDEFIND`
2. **Analizator składniowy** (parser):
   - Analizator składniowy jako wejście przyjmuje strumien tokenów wyprodukowany przez analizator leksykalny.
   - Zadaniem parsera jest wyprodukowanie drzewa rozbioru składniowego programu w formie `node'ów`.
   - Ściśle oczekuje na spodziewany token przy analizie wyrażenia.
   - Obsługa błędów składniowych realizowana poprzez `panic()`, zawierające informacje o położeniu błędnego wyrażenia w kodzie programu.
3. **Interpreter**:
   - Operuje na drzewie rozbioru składniowego.
   - Napisany przy użyciu wzorca projektowego "Visitor (wizytator)".
   - Interpreter odwiedza elementy drzewa skłądniowego, ewaluując ich zawartość. Nadaje wartości zmiennym, sprawdza zgodność typów, zgodność podawanych do wywołań argumentów, uruchamia wywoływane funkcje (również funkcje wbudowane).
   - Dba o to aby wywołania rekurencyjnie nie przekroczyly zdefiowanego limitu (implementacja za pomocą CallStack).
   - Wykonuje operacje arytmetyczne, obsługuje instrukcje warunkowe, pętle, wywołania funkcji oraz inne konstrukcje językowe.

---

## Testy

### Scanner

- TestScanner: Zbiór testów sprawdzający poprawność działania skanera dla różnych przykładowych wejść. Testy sprawdzają, czy wyniki zwrócone przez skaner są zgodne z oczekiwanymi znakami i ich pozycjami. Testuje same ciągi znaków jak i ciągi znaków ze znakami nowej lini i poprawnośc zachowania w przypadku `\r\n` oraz `\r\t`.

- TestScannerMultipleEOF: Testuje, czy skaner poprawnie zwraca EOF i jego pozycję w tekście. Definiuje pusty `input` i oczekuje, że skaner zwróci kilka EOF z tą samą pozycją.

### Lexer

Testy jednostkowe (unit tests):

- TestSingleTokens: Zbiór testów sprawdzający poprawność tworzenia każdego z tokenów.

Testy ciągu tokenów:

- TestLexerCodeExample: Zbiór testów sprawdzający poprawność analizy sekwencji tokenów w przykładowym kodzie źródłowym.

Testy obsługi błędów:

- TestStringNotClosed: Sprawdza, czy analizator leksykalny poprawnie zgłasza błąd w przypadku niezamkniętych ciągów znaków.
- TestIntValueLimitExceeded: Testuje, czy analizator leksykalny zgłasza błąd, gdy wartość liczby całkowitej przekracza limit.
- TestLexerStringTokenEscaping: Sprawdza, czy analizator leksykalny prawidłowo obsługuje znaki specjalne w ciągach znaków.
- TestLexerInvalidStringTokenEscaping: Testuje, czy analizator leksykalny zgłasza błąd dla niepoprawnych znaków specjalnych w ciągach znaków.
- TestStringValueLimitExceeded: Sprawdza, czy analizator leksykalny zgłasza błąd, gdy długość ciągu znaków przekracza limit.
- TestLexerErrorHandling: Testuje, czy analizator leksykalny poprawnie obsługuje błędy i zatrzymuje się po napotkaniu błędu, a także czy zgłasza oczekiwane błędy.

### Parser

Testy sprawdzające poprawnośc parsowania ciągu tokenów wejsciowych na elementy drzewa AST, od terminalnych do bardziej złożonych (zagłębionych).
Testy odpowiedzialne są również za sprawdzanie czy w określonych przypadkach wyrzucane są błędy z poprawnymi treściami.

- TestParseParameterGroup
- TestParseParameters
- TestParseIdentifier
- TestParseFunctionDefinitions
- TestParseExpressionIdentifierAsTerm
- TestParseExpressions: Zbiór testów sprawdzających poprawność każdego z expression zdefiowanego w EBNF
- TestParseVariableDeclaration: Zbiór testów sprawdzających poprawnośc przypisywania wartości do zmiennej
- TestParseNegateExpression
- TestParseVariableDeclarationsWithTerm
- TestParseFloatExpression
- TestParseStringExpression
- TestParseParenthesisExpression
- TestParseOrAndExpression
- TestParseIfStatement
- TestParseIfStatementWithElse
- TestParseWhileStatement
- TestSwitchStatement
- TestSwitchStatementWithDefault
- TestSwitchStatementWithBlock
- TestParseSwitchError
- TestParseProgram
- TestParseProgramFunctionWithType

### Interpreter

Testy wyrażeń liczbowych:

- TestVisitIntExpression: Sprawdza, czy odwiedzając wyrażenie liczbowe, otrzymuje się prawidłowy wynik.
- TestVisitFloatExpression: Testuje poprawność przetwarzania wyrażenia zmiennoprzecinkowego.
- TestVisitStringExpression: Sprawdza poprawność przetwarzania wyrażenia tekstowego.
- TestVisitBoolExpression: Testuje poprawność przetwarzania wyrażenia logicznego.

Testy wyrażeń logicznych:

- TestVisitAndExpression: Sprawdza, czy poprawnie przetwarzane są wyrażenia AND.
- TestVisitOrExpression: Testuje poprawność przetwarzania wyrażeń OR.

Testy wyrażeń arytmetycznych:

- TestVisitSumExpression: Sprawdza, czy wyrażenia sumy są przetwarzane poprawnie dla różnych typów danych.
- TestVisitSumExpressionString: Testuje przetwarzanie sumy wyrażeń tekstowych.
- TestVisitSubstractExpressionInt: Testuje poprawność przetwarzania wyrażenia odejmowania dla wartości integer.
- TestVisitSubstractExpressionFloat: Testuje poprawność przetwarzania wyrażenia odejmowania dla wartości float.
- TestVisitSubstractExpressionFloatMinusInt: Testuje poprawność przetwarzania wyrażenia odejmowania dla wartości float i integer.
- TestVisitSubstrackExpressionIntMinusFloat: Testuje poprawność przetwarzania wyrażenia odejmowania dla wartości integer i float.

Inne testy:

- TestVisitNegateExpression: Testuje przetwarzanie negacji dla różnych typów wyrażeń.

- TestCastExpression: Zbiór testów sprawdzający poprawność przetwarzania wyrażenia kastowania.

Testy innych elementów drzewa AST:

- TestVisitIdentifier: Sprawdza poprawność przetwarzania identyfikatorów, czy wartość ze Scope zostanie zwrócona.
- TestVisitVariable: Sprawdza poprawność przetwarzania zmiennej, czy zmienna zostaje poprawnie dodana do Scope.
- TestVisitWhileStatement: Sprawdza, czy wartość z wywołania pętli.
- TestVisitIfStatementConditionTrue: Sprawdza, czy wartość z wywołania bloku if jest przenoszona do poprzdniego Scope poprzez LastResult.
- TestVisitIfStatementElseBlock: Sprawdza czy wartość LastResult została poprawnie wyczyszczona wychodząc z bloków if.
- TestScopeVariables: Testuje czy zmienne są poprawnie przechowywane w Scope.
- TestVisitFunctionCall: Testuje Wywołanie funkcji.
- TestVisitSwitchCase: Testuje poprawność przetwarzania switch-case.
- TestVisitDefaultSwitchCase: Testuje poprawność przetwarzania default-case.
- TestVisitSwitchStatement: Testuje poprawność przetwarzania switch statement ze zwracaniem wartości.
- TestSwitchWithBlock: Testuje poprawność przetwarzania switch statement z blokiem.
- TestRecursion: Testuje poprawność przetwarzania funkcji rekurencyjnych i przekraczania limitu wywołań rekurencyjnych.

Testy złożonych funkcji (fragmentów kodu):

- TestReturningNestedBlocks: Sprawdza, funkcja jest poprawnie zamykana po osiągnięciu pierwszego **return**.
- TestVisitFunctionCallWithIdentifier: Testuje czy wywołanie funkcji odbywa się poprawnie z przekazaniem identyfikatorów jako argumenty.
- TestVisitAsignmentWithFunctionCall: Testuje czy przypisanie do zmiennej odbywa się poprawnie z wywołaniem funkcji.
- TestVisitNestedFunctionCallWithReturn: Testuje zwracanie poprawnej wartości z zagnieżdżonych Scope'ów.
- TestVisitWhileStatementWithReturn: Testuje poprawność przetwarzania pętli ze zwracaniem wartości.

Testy obsługi błędów:

- TestVariableNotInScope: Sprawdza czy poprawnie zostal wyrzucony błąd z niezdefiniowaną zmienną.
- TestSearchVariableInScope: Sprawdza czy wartośc z wywołania funkcji nie jest przenoszona do poprzdniego Scope poprzez LastResult. Wyłapuje błąd z odpowiednim komunikatem o niezdefiniowanej zmiennej.
- TestParametersAndArguments: Testuje poprawności wyłapywanych błędów podczas przetwarzania niepoprawnej ilości argumentów i parametrow lub podania niepoprawnych typów.
- TestFunctionWithSwitch: Sprawdza czy poprawnie zostanie wyłapany błąd, ze względu na brak return'a i switch_case'ów które zwróciły by wartości, gdy funkcja w której zdefiniowany jest switch powinna takową zwrócić.
- TestVariableDeclarationMissmatch: Testuje poprawnośc wyłapanych błędów i ich komunikatów.
- TestEmbededFunction: Testuje poprawność przetwarzania funkcji wbudowanych.
