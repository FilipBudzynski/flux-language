# TKOM - Dokumentacja

Filip Budzyński, numer albumu: 319021

## Opis zakładanej funkcjonalności

Projekt języka “Flux” zakłada stworzenie języka programowania, który umożliwia inicjalizację i przypisywanie zmiennych, wykonywanie operacji arytmetycznych, obsługę instrukcji warunkowych i pętli, definiowanie funkcji z argumentami, funkcje rekurencyjne oraz obsługę relacyjnych wzorców w instrukcjach switch.

Implementacja wszystkich elementów napisana w języku **Golang**.

## Standardowe operacje:

---

- obsługa inicjalizacji i przypisania zmiennej
- wykonywanie operacji arytmetycznych
- instrukcje warunkowe
- instrukcja pętli
- funkcje rekurencyjne
- definiowanie funkcji z lub bez argumentów
- wywołania zwykłe i rekurencyjne funkcji

## Charakterystyczne operacje:

---

- Dopasowanie wzorców relacyjnych
    - implementowane przy użyciu słowa kluczowego “switch”
    - argumentem “switch” mogą być wyrażenia w tym zmienne lub wywołania funkcji np:
        - switch 2*(12+3) {}
        - switch umUP(2, 2) {}
        - switch whatWord() {}
        - *przykłady podane dalej w dokumentacji*

## Dopuszczalne typy

---

- integer
- float
- string
- boolean

## Założenia języka

---

- statycznie typowany
- argumenty przekazywane przez wartość
- obsługa błędów na poziomie leksykalnym i semantyczny oraz na poziomie interpretera
- zmienna może być nadpisana, ale musi zgadzać się jej typ
- każdy program napisany w języku “flux” musi posiadać funkcję main(), która zawiera główne ciało programu i od tej funkcji rozpoczyna się jego działanie.
- Zmienna zewnętrzna może zostać przykryta poprzez zmienną o tej samej nazwie znajdującej się w bloku funkcji, pętli, instrukcji warunkowej lub instrukcji switch.

## Przykłady dopuszczalnych konstrukcji i semantyka

---

Inicjalizacja i przypisanie wartości

```go
a := 5
b := 2

a = 8
```

---

Operacje arytmetyczne

```go
a := 3
a = a + 3 * (2 - 1)
```

---

Komentarze

```go
# *To jest komentrz*
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

nazwa := "Ala ma psa"
if nazwa == "Ala ma kota" {
    print("Kot należy do Ani")
} else {
    print("Ani to Ala ani kot")
}
```

---

Instrukcja pętli for

```go
i := 0
while i < 10{
    print(i)
    i = i + 1
}

num := 10
while num > 0 {
    print(num)
    num = num - i
}
```

---

funkcja z argumentem

```go
circleArea(r int) float    {
    return 3.14 * (r * r)
}

r := 2
a := circleArea(r)        # a == 12.56636
```

---

funkcja rekurencyjna

```go
fibonacci(n) {
  if n <= 1 {
    return n
  } else {
    return fibonacci(n - 1) + fibonacci(n - 2)
  }
}
result := fibonacci(5)

2 + 4 as float
```

---

Konwersja typów

```go
a := 5
c := a as string
print(c)        # "5"

b := 0
b = b as bool   # "false"
```

Relational patterns - switch

```go
sumUp(a,b int) int {
    return a + b
}

whatWillGetMe() string {
    switch sumUp(2, 3) {
        >2 and <=4 => "A pint",
        5          => "Decent beverage",
        >5 and <15 => "A NICE bevrage",
        > 15       => "Whole bottle"
    }
}
```

```go
giveMeWord() string {
    return "word"
}

nameNumber() int {
    switch giveMeWord() {
        "Sammy"  => 0,
        "World"  => 1,
        "word"   => 2,
        default  => 3
    }
}

main(){
  print(nameNumber())
}

# ourput:
$ 2
```

```go
getUserRole(userId int) string {
    return "admin"
}

checkPermission(role string, permission string) bool {
    return role == "admin" && permission == "edit"
}

main() {
    userId := 123

    switch userRole := getUserRole(userId) {
        "admin" => {
            if checkPermission(userRole, "edit") {
                print("Użytkownik ma uprawnienia do edycji")
            } else {
                print("Użytkownik nie ma uprawnień do edycji")
            }
        },
        "user" => {
            print("Użytkownik ma ograniczone uprawnienia")
        },
        default => {
            print("Nieznana rola użytkownika")
        }
    }
}
```

## Specyfikacja i składnia EBNF

---

- symbole terminalne wyróżnione znakiem *

```go
program               = { function_definition } ;

function_definition   = identifier , "(", [ parameters ], ")", [ type_annotation ] , block ;

parameters            = parameter_group , { "," , parameter_group } ;
parameter_group       = identifier , { ",", identifier }, type_annotation ;

type_annotation       = "int" | "float" | "bool" | "str" ;

block                 = "{" , { statement } , "}" ;

statement             = variable_declaration
                      | variable_assignment
                      | conditional_statement
                      | loop_statement
                      | function_call
                      | type_cast
                      | switch_statement
                      | return_statement
                      ;

variable_declaration  = identifier , ":=" , expression ;

variable_assignment   = identifier , "=" , expression ;

conditional_statement = "if" , expression , block , [ "else" , block ] ;

loop_statement        = "while" , expression, block ;

function_call         = identifier , "(", arguments, ")" ;

switch_statement      = "switch", expression, "{", switch_case, { ",", switch_case "}" ;
switch_case           = ( ( [relation_operator], expression ) | "default" ), "=>", ( expression | block ), } ; 

return_statement      = "return" , [ expression ] ;

arguments             = expression , { "," , expression } ;

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
multiplicative_term   = unary_operator, { ("*" | "/"), unary_operator } ; 
unary_operator        = [ ("-" | "!") ], type_cast ;
type_cast             = term, [ "as", type_annotation ] ;
term                  = integer
                      | float
                      | bool
                      | string
                      | identifier
                      | function_call
                      | type_cast
                      | switch_element
                      | "(" , expression , ")"
                      ;

identifier            = letter , { letter | digit | "_" } ;

float                 = integer , "." , digit , { digit } ;

integer               = "0" | positive_digit , { digit } ;

string                = '"', { literal }, "\n", '"' ;

literal               = letter
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

## Obsługa błędów i przykłady

---

Format błędów:

```go
<Error | Warning> [<line> : <column>]: <message>
```

---

Błąd przypisania:

```go
main(){
  a := 5
  a = "Ala ma kota"
}
```

```go
Error [3:4]: Invalid value assignment. The type of variable a is numeric, and you are trying to assign it as string. 
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
Error [2:10]: Cannot use int as return value, function should return float.
```

---

Błąd w konstrukcji switch:

```go
kelvinToCelcius(temp int) int {
  return temp - 273
}

howCold(kelvin int) string {
  switch kelvinToCelcius(kelvin) {
    <(-20)       => "Freezing",
    >0 and <10   => "Chilling",
    >=10 and <20 => "Warm",
    >=20         => "Hot"
  }
}

main(){
  weather := kelvinToCelcius(300)
  print(weather)
}
```

```go
Error [6:34]: Not all cases matched for the switch statement.
```

---

Niezainicjowana zmienna:

```go
main(){}
  print(a + 10)
}
```

```go
Error [2:9]: Undefind variable 'a'.
```

---

Błąd niezadeklarowanej funkcji:

```go
main() {
  print(unknownFunction())
}
```

```go
Error [2:9]: Undefined function 'unknownFunction'.
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
Error [5:12]: Function 'add' expects 2 arguemnts, but 1 provided.
```

---

Błąd niezgodności typów w instrukcji warunkowej:

```go
main() {
  a := 5
  if a == "test" {
    print("Equal")
  }
}
```

```go
Error [2:11]: Incompatible types in conditional statement.
```

---

Błąd niepoprawnego użycia operatora:

```go
main() {
  a := 10
  b := "5"
  print(a / b)
}
```

```go
Error [4:11]: Invalid operation. Division operator cannot be applied to operand of type int and string.
```

---

Błąd niepoprawnego użycia operatora relacyjnego:

```go
main() {
  a := 5
  if a < "test" {
    print("Less than")
  }
}
```

```go
Error [3:8]: Invalid comparison. Eelational operation cannot be applied to types int and string.
```

---

Błąd dzielenia przez zero:

```go
main() {
  a := 10
  b := 0
  result := a / b
}
```

```go
Error [3:12]: Division by zero is prohibited.
```

## Rozróżniane Tokeny

---

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
    - `**NEGATE**` ( “-” lub “!”)
5. Słowa kluczowe:
    - **`IF`**
    - **`ELSE`**
    - **`FOR`**
    - **`SWITCH`**
    - **`DEFAULT`**
    - `**AS**`
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
- **`INTEGER`**
- **`FLOAT`**
- **`STRING`**
- **`BOOL`**
1. Funkcje i procedury:
    - **`FUNCTION_DEFINITION`**
    - **`FUNCTION_CALL`**
    - **`PARAMETER`**
    - **`ARGUMENT`**
2. Ostrzeżenia i błędy:
    - **`ERROR`**
    - **`WARNING`**
3. Inne:
    - **`ERROR_MESSAGE`**
    - **`WARNING_MESSAGE`**
    - **`STX`** (start of text) ****
    - **`ETX`** (end of text) ****
4. Struktura tokenu
    - `**TOKEN_TYPE**`
    - `**UNDEFIND_TOKEN**`
    - `**POSSITION**`

## Specyfikacja danych wejściowych strumienie/pliki i danych konfiguracyjnych

---

Program napisany w języku Flux może być uruchamiany zarówno z pliku, jak i ze strumienia danych wejściowych.

Jeśli program jest przekazywany ze standardowego wejścia, wszystkie bajty aż do pierwszego znaku końca pliku (EOF) są traktowane jako program do interpretacji, a wszystkie bajty po pierwszym znaku EOF są traktowane jako standardowe dane wejściowe.

Język Flux nie wymaga żadnych specjalnych danych konfiguracyjnych do poprawnego działania.

Interpreter programu dostaje dostęp do standardowego wyjścia i wejścia, co pozwala na przechwytywanie wyników działania programu i pokazywanie błędów oraz na podanie danych wejściowych do programu.

## Sposób uruchomienia

---

Przykłady uruchomienia programu napisanego w języku flux:

Z Pliku

```go
$ flux program.fl
Hello Word from file! 
```

---

Ze strumienia

```go
$ echo 'main(){ print("Hello Word from steam!") }' | flux 
Hello Word from steam!
```

## Zwięzła analiza wymagań funkcjonalnych i niefunkcjonalnych

---

1. Wymagania Funkcjonalne:
    - Obsługa Inicjalizacji i Przypisania Zmiennych.
    - Wykonywanie Operacji Arytmetycznych np.dodawanie, odejmowanie, mnożenie i dzielenie.
    - Obsługa Instrukcji Warunkowych.
    - Instrukcje pętli**,** obsługa pętli while.
    - Funkcje Rekurencyjne: Zapewnienie wywoływanie funkcji rekurencyjnych.
    - Obsługa Relacyjnych Wzorców w Instrukcjach Switch**:** Implementacji instrukcji switch, która może dopasowywać się do relacyjnych wzorców w celu podejmowania decyzji.
2. Wymagania Niefunkcjonalne:
    - Statyczne i silne typowanie.
    - Obsługa błędów i ostrzeżeń na poziomie leksykalnym i semantycznym.
    - Komunikacja o błędach w omówiony wcześniej sposób.
    - Dopuszczalne typy oraz konwersja integer, float, string i boolean.
    - Możliwość zdefiniowania programu w pliku lub podanie programu poprzez strumień.

## Zwięzły opis sposobu realizacji modułów

---

1. **Analizator leksykalny**:
    - Przetwarza kod źródłowy i zgodnie z gramatyką produkuje tokeny do identyfikacji i grupowania leksemów, takich jak identyfikatory, liczby, operatory i słowa kluczowe.
    - Tokeny przechowują  informacje o swoim położeniu w kodzie źródłowym w postaci `(nr linii, nr kolumny)`.
    - Implementacja analizatora leksykalnego oparta na skanowaniu kodu programu zgodnie z gramatyką języka podaną powyżej w EBNF.
    - W przypadku natrafienia na niemożliwy do zdekodowania ciąg znaków, analizator skanuje ciąg aż do natrafienia na biały znak i zwraca token `UNDEFIND_TOKEN`
2. **Analizator składniowy**:
    - Analizator składniowy jako wejście przyjmuje strumien tokenów wyprodukowany przez analizator leksykalny.
    - Zadaniem parsera jest wyprodukowanie drzewa rozbioru składniowego programu.
    - Na podstawie dostarczonych tokenów, oczekuje na tokeny określonego typu. ściśle oczekuje na spodziewany token przy analizie wyrażenia.
    - Obsługa błędów składniowych realizowana poprzez wyjątki, zawierające informacje o położeniu błędnego wyrażenia w kodzie programu.
3. **Analizator semantyczny**:
    - Analizator semantyczny operuje na drzewie rozbioru składniowego wyprodukowanym przez analizator składniowy.
    - Zadaniem analizatora semantycznego jest kontrola poprawności analizowanego kodu.
    - Sprawdza zgodność typów, poprawność zasięgu zmiennych, odwołania do zmiennych, poprawność wywołań funkcji i innych właściwości programu. Jest również odpowiedzialny za kontrolę typów zmiennych.
    - W przypadku natrafienia na nieścisłość, analizator, rzuca wyjątkiem, zawierającym informacje o położeniu niepoprawnego semantycznie kawałka kodu.
4. **Interpreter**:
    - Operuje na drzewie rozbioru składniowego, dopiero po pomyślnym zakończeniu analizy semantycznej.
    - Zadaniem interpretera jest nadawanie wartości zmiennym oraz wykonywanie instrukcji zawartych w zdefiniowanych funkcjach a także funkcji wbudowanych np. `print()`.
    - Wykonuje operacje arytmetyczne, obsługuje instrukcje warunkowe, pętle, wywołania funkcji oraz inne konstrukcje językowe.

## Jak będzie realizowane przetwarzanie w poszczególnych modułach

---

1. Analizator leksykalny
    - Skanuje  wejście, znak po znaku, identyfikując leniwie kolejne leksemy (np. słowa kluczowe, identyfikatory, liczby, znaki specjalne) i przypisuje im odpowiednie tokeny.
    - Pomija znaki białe.
2. Analizator składniowy
    - Odbiera sekwencję tokenów wyprodukowanych przez lekser i tworzy na ich podstawie drzewo rozbioru składniowego (AST).
    - Sprawdza, czy struktura tokenów odpowiada zdefiniowanej gramatyce języka.
3. Analizator semantyczny
    - Operuje na drzewie AST stworzonym przez parser.
    - Przechodzi po węzłach drzewa i przeprowadza analizę relacji między węzłami, zgodności typów, odwołań, wywołań funkcji itd.
4. Interpreter
    - Odwiedza węzły drzewa AST w zależności od przebiegu działania programu.
    - Wykonuje odpowiednie akcje i kod w zależności od węzłów, które odwiedza.

## Konwersja typów w tabelce i jaka kombinacja typów jest akceptowalna dla operatorów wieloargumentowych i funkcji wbudowanych

---

Ponieważ język jest silnie i statycznie typowany, każda konwersja typu jest jawna a do jej dokonania dostępny jest operator “as”.

Konwersja typów dla typowania statycznego:

| Z | Do Integer | Do Float | Do String | Do Boolean |
| --- | --- | --- | --- | --- |
| Integer | - | Explicit | Explicit | Explicit |
| Float | Explicit | - | Explicit | Explicit |
| String | Explicit | Explicit | - | Explicit |
| Boolean | Explicit | Explicit | Explicit | - |

W przypadku int na boolean:

- int 0 oznacza                   `false`
- inne poza 0 oznaczają    `true`

W przypadku string na boolean:

- pusty string: “” oznacza `false`
- niepusty string oznacza `true`

W przypadku float na boolean:

- float 0.0 oznacza             `false`
- inne poza 0  oznacza      `true`

## Zasady przekazywania zmiennych do funkcji

---

Zmienne są przekazywane do funkcji przez wartość. Jako, że nie ma struktur to przekazywanie zmiennej przez referencje nie wydaje się być konieczne.

Przykład przykrywania zmiennych:

```go
main(){
  a := "Outside"
  if true {
    a = "Inside"
    print(a)
  }
  print(a)
}

# output
Inside
Outside
```

```go
main(){
  temp := 20
  switch temp = 30 {
    <=0          => { print("Cold") },
    >0 and <10   => { print("Chilling") },
    >=10 and <20 => { print("Warm") },
    >=20         => { print("Hot") }
  }
}

# temp zdefiniowane przed swtich zostaje przykryte przez zmienną temp w swtich
# output
Hot
```

## Czy dozwolone jest przeciążanie funkcji?

---

Przeciążanie funkcji nie jest dozwolone.
