
# <span style="color:#6495ED">**Go Lang Basic Concept**</span>


![alt text](/images/go_lang.png)

Below Note cover the Go Language Base Fundamentals

**Table of content:**
1. [Introduction](#1-introduction)
2. [Data Types](#2-data-type)
3. [Variable](#3-variable)
4. [Control Flow](#4-control-statement--flow)
5. [Array](#5-array)
6. [Map](#6-map)
7. [Functions](#7-functions)
8. [Pointers](#8-pointers)
9. [Structs and Interface](#9-structs--interfaces)
10. [Packages](#10-packages)
11. [Go Test](#11-go-test)
12. [Concurrency](#12-concurrency)
13. [CLI](#13-command-line-reference)
14. [Core Pkg](#14-core-package)
15. [Go Library](#15-library-example)
16. [Reference](#16-reference-links)
17. [QA](#17-qa)

<a id="1-introduction"></a>
## <span style="color:#6495ED">1. Introduction<span>

The Go toolset uses an environment variable called GOPATH to find Go source code.

In Go Lang when go run command executed with the subsequent files (separated by spaces), it compiles them into an 
executable & saved in a temporary directory, and then runs the program.

Sample Program

---

```go
package main
import "fmt"

func main() {
    fmt.Println("Hello World")
}
    
```

main(): The name main is special because it's the function that gets called when we execute the program

```sh
go run main.go
```

There are two types of Go Programs

1. Executable : Executable applications are the kinds of program that can be run directly from the terminal 
(On Windows they end with .exe).

2. Libraries : Collections of code or programs that are package together which can be used in other programs.

### Comments

* Go support two different type of comments // (Single Line Comment) & /**/ (Multi Line Comment)

<a id="2-data-type"></a>
## <span style="color:#6495ED">2. Data Type

### 2.1 Integer

* Go’s integer types are uint8, uint16, uint32, uint64, int8, int16, int32, and int64. uint means “unsigned integer” 
while int means “_signed integer._”
* In this 8, 16, 32 and 64 represent number of bits it can hold.

### 2.2 Floating Point

* Go has two floating-point types: float32 and float64 (also often referred to as single
precision and double precision, respectively). It also has two additional types for representing complex numbers 
(numbers with imaginary parts): complex64 and complex128.

### 2.3 String

* String is a sequence of characters with a definite length
  used to represent text.

```go
  len("Hello, World")
  Finds the length of a string
```

```go
  "Hello, World"[1]
  Accesses a particular character in the string (in this case, the second character)
  But it returns 101 instead of e when you run this program. This is because the character is represented by a byte
```

```go
  "Hello, " + "World"
  Concatenates two strings together
```

### 2.4 Boolean

* Three logical operators are used with boolean values:

```go
  && and
  || or
  ! not
```

> **_NOTE:_** <span style="color:orange"> 
To print the data type of variable %T flag can be use in fmt.Printf() functions.
  ```go
  var numA = 50
  fmt.Printf("The data type of numA is : %T", numA)
  ```
</span>

<a id="3-variable"></a>
## <span style="color:#6495ED">3. Variable

* Variable can be defined in three ways :

  ```go
  Normal Declaration without type 
  var normStr = "Hello, World"
  ```

  ```go
  Normal Declaration with type
  var normStrType string = "Hello, World"
  ```

  ```go
  Short Declaration
  shortStr := "Hello, World"
  ```

* Notice the : before the = and that no type was specified. The type is not necessary
because the Go compiler is able to infer the type based on the literal value you assign the variable (because you are 
assigning a string literal, x is given the type string).
  
* Defining Multiple Variable

  ```go
  var (
    numA = 5  
    numB = 10
    numC = 20
    )
  ```

* Variable name must start with letter and may contain letter, underscore or number & camelCase is used for naming 
a variable.

> **_NOTE:_** <span style="color:orange">The Go compiler won’t allow you to create variables that you never use.</span>

> **_NOTE:_** <span style="color:orange"> When we declare variable in Go lang without intiliaziation 
Go puts up default value for the type.
For string : "" , int : 0, float: 0, bool: false
</span>

* constants are essentially variables whose value cannot be changed later.

  ```go
  const x string = "Hello World"
  ```
<a id="4-control-statement--flow"></a>
## <span style="color:#6495ED"> 4. Control Statement & Flow

* Go only support for loop which can be used in a different ways. Following examples are mentioned :

```go
i := 1
for i <= 10 {
    fmt.Println(i)
    i = i + 1
}

// OR

for i := 1; i <= 10; i++ {
    fmt.Println(i)
}
```

* for infinite loop in Go language we can simply provide for loop without any condition as mentioned below

```go
for {
  fmt.Println("Infinite Loop")
}
```

* if, else ,else if and switch control statements are used in Go similar to other programming languages.
Following examples are mentioned :

```go
num := 10 
if (num % 2 == 0) {
    fmt.Println("Number is divisible by 2 and it is even") 
} else if (num % 3 == 0) {
    fmt.Println("Number is divisible by 3")
} else {
    fmt.Println("Number is not divisible by 2 or 3")
}
```

```go
num := 2

switch num {
    case 0: fmt.Println("Zero")
    case 1: fmt.Println("One")
    case 2: fmt.Println("Two")
    case 3: fmt.Println("Three")
    default: fmt.Println("Unknown Number")
}
```

<a id="5-array"></a>
## <span style="color:#6495ED"> 5. Array

* Array is fixed length list of things 
### 5.1 Declaration of Arrays

```go
var arr []float64 // array been created with length zero
var arr [5]int
var arr [5]float64
arr := [5]float64{81, 82, 83, 84, 85}
arr := [5]float64{
        81,
        82,
        83,
        84,
        85
      }
```

### 5.2 Iteration and Access over Array

```go
var arr[5]int
arr[4] = 100 
fmt.Println(arr)

for i:=0; i < len(arr); i++ {
    fmt.Println(arr[i])
    }

for _, value := range arr {
    fmt.Println(value)
    }
```

### 5.3 Array Slice, Append and Copy

* Slice: To create a slice we need to use built-in make function. Slices are always associated with some array, and 
although they can never be longer than the array, they can be smaller.
* Slice is an array that can grow or shrink

```go
arr_slice := make([]float64, 5)
creates a slice that is associated with an underlying float64 array of length 5

arr_slice := make([]float64, 5, 10)
creates a slice of length 5 with a capacity of 10

arr := [5]float64{1,2,3,4,5}
arr_slice := arr[0:3]
```

* Append: Adds elements onto the end of a slice. If there is not sufficient capacity, a new array is created, all of 
the existing elements are copied over, the new element is added onto the end, and the new slice is returned.

* Copy: copy takes two arguments dst and src. All the entries of src are copied into dst overwritng whatever is there. 
If the lengths of two slices are not the same, the smaller of the two will be used.

```go
slice1 := []int{1,2,3}
slice2 := make([]int, 2)
copy(slice2, slice1)
fmt.Println(slice1, slice2)
```
Reference Link: https://go.dev/blog/slices-intro

<a id="6-map"></a>
## <span style="color:#6495ED"> 6. Map

* A map is an unordered collection of key-value pairs. The map type is represented by the keyword map, followed by the 
key type in brackets and finally the value type. 

* Map returns the zero value for the value type if the key is not present in the map. 

* key are statically type so it should be of same data type when defined.

Following example are different way to declare and intialize map.

```go
elements := make(map[string]string)
elements["H"] = "Hydrogen"
elements["He"] = "Helium"
elements["Li"] = "Lithium"

elements := map[string]string{
    "H": "Hydrogen",
    "He": "Helium",
    "Li": "Lithium",
  }

elements := map[string]map[string]string{
    "H": map[string]string{
        "name":"Hydrogen",
        "state":"gas",
      },
    "He": map[string]string{
        "name":"Helium",
        "state":"gas",
      },
    "Li": map[string]string{
        "name":"Lithium",
        "state":"solid",
      }
  }
```

* length of map

```go
len(elements)
```

* delete items from map

```go
delete(elements, "He")
```

* Accessing an element of a map can return two values instead of just one. The first value is the result of the lookup, 
the second tells us whether or not the lookup was successful. Below example for the result of map access:

```go
name, ok := elements["num"]
fmt.Println(name,ok)
```

<a id="7-functions"></a>
## <span style="color:#6495ED"> 7. Functions

Functions start with keyword func, followed by the functions name. The parameters(input) of the function defined 
by the name of parameters and followed by its type. After parameter(input) we provide return type. Parameters and 
return type are known as functions signature.

* Some points related to function need to be noted :
  1. Parameters name can be different
  2. Variables must be passed to the function
  3. Functions form call stack
  4. Return types can have name
  5. Multiple values can be returned

```go
  func f1() (r int) {
    r = 1
    return r
    }

  func f2() (int, int) {
    return 5, 6
  }
  
  func main() {
    x = f1()
    y, z := f2()
  }
```

* Multiple values are often used to return an error value along with the result (x, err := f()) 
or a boolean indicate sucess (x, ok := f())

### 7.1 Variadic Functions

* The function that is called with the varying number of arguments is known as variadic function. 
Or in other words, a user is allowed to pass zero or more arguments in the variadic function.

* Functions is allowed to be called with multiple values which is known as variadic parameter. 
By using elipsis (...) before the type name of last parameter, which can indicate that function require 
zero or more parameter.

```go
func add(args ...int) int {
  total := 0.0
  for _, num := range args {
    total += num
  }
  return total
}

func main() {
  fmt.Println(add(1,2,3,4,5))
  xs := []int{1,2,3}
  fmt.Println(add(xs...))
}
```

### 7.2 Closure

* Create a functions inside a function, One way to use closure is by writing a function that returns another function, 
which when called can generate a sequence of numbers. 
* Closure and recursion are powerful programming techniques that form the basis of a paradigm known as 
functional programming

```go
num := 0
increment := func() int {
  num ++
  return num
  }
fmt.Println("Increment done first time: ", increment())
```

### 7.3 Recursion

Function is able to call itself, For example a factorial function 

```go
func factorial(num int) int {
    if num == 0 {
        return 1
    }
    return num * factorial(num-1)
}
```

### 7.3 defer, panic & recover

#### 7.3.1 defer: 

* Go has a special statement called defer that schedules a function call to be run after the function completes. 
* defer is often used when resource need to be freed.
* defer moves the call to the end of the function:

For example when we open a file we need to make sure to close it later
```go
f, _ := os.Open(filename)
defer f.Close()
```
* This has three advantages:
  * It keeps our Close call near our Open call so it’s easier to understand. 
  * If our function had multiple return statements (perhaps one in an if and one in
  an else), Close will happen before both of them. 
  * Deferred functions are run even if a runtime panic occurs.

#### 7.3.2 panic:

* The panic function is used to terminate a program abruptly.
* A panic generally indicates a programmer error (e.g., attempting to access an index
of an array that’s out of bounds, forgetting to initialize a map, etc.) or an exceptional
condition that there’s no easy way to recover from (hence the name panic).
* When a panic occurs, the normal flow of the program is interrupted, and the program terminates.

#### 7.3.3 recover:

* The recover function is used to regain control after a panic.
* It can only be used in a deferred function.
* It returns the value passed to the call to panic.
* If there was no panic, or if the recover function is not in a deferred function or a goroutine, it returns nil.


```go
import "fmt"

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    
    result := divide(10, 0)
    fmt.Println(result)
}

func divide(a, b int) int {
    if b == 0 {
        panic("Division by zero")
    }
    return a/b
}
```

<a id="8-pointers"></a>
## <span style="color:#6495ED"> 8. Pointers

> **_NOTE:_** <span style="color:orange"> Go language is pass by value. To overcome this problem pointer is used.
</span>

> **_NOTE:_** <span style="color:orange"> In Go language there are value type and reference type data types.
* Value types: "int, float, string, bool, structs". Use pointers to change these things in a function.
* Reference types: "slices, maps, channel, pointers, functions". No need of pointer to update or change these thing in a function.
</span>


* Pointer concept is same as in C language, Pointers reference a location in memory where a value is stored rather 
than the value itself. 
* In Go a pointer is represented using an asterisk (*) followed by the type of the stored value 
and & operator to find the address of a variable or pointer is pointing to that address.

```go
func zero(xPtr *int) {
  *xPtr = 0
}

func main() {
  num := 5
  zero(&num)
  fmt.Println(num) // xPtr is 0
}
```

* Another way to declare a pointer using built-in new function. new takes a type as an argument, allocates enough 
memory to fit a value of that type, and returns a pointer to it.

```go
func one(xPtr *int) {
  *xPtr = 1
}

func main() {
  xPtr := new(int)
  one(xPtr)
  fmt.Println(xPtr) // xPtr is 1
}
```

<a id="9-structs--interfaces"></a>
## <span style="color:#6495ED"> 9. Structs & Interfaces

* A struct is a type of data structure which is collection of same or different properties that are related together.
* A struct is a type that contains named fields. 


Following example for Declaration as well as Initialization of Structs.

```go
type Circle struct {
  x float64
  y float64
  r float64
}

func main() {
  var c Circle // Default value zero
  // Initialize with values
  var c1 Circle = Circle{1, 1, 4}

  // Decalration and Initialization of struct type 
  c2 := Circle{1, 2, 3} 
}
```

### 9.1 Methods and Embedded Types

### 9.2. Interface

An interface is created using the type keyword, followed by a name and the keyword interface.
But instead of defining fields, we define a method set. 
A method set is a list of methods that a type must have in order to implement the interface.

```go
type Shape interface {
  area() float64
}

```

<a id="10-packages"></a>
## <span style="color:#6495ED"> 10. Packages

* A package is a collection of source files.
* Go programs are organized into packages
* Go's standard library, provdides different core packages for us to use

Go contains two types of packages:

1. ```Executable``` :

* Generate a files that can be run. When we define package main in our main.go file then the word main is used to 
make executable file. 

* pacakge main => go build main.go => main.exe.
main.exe file is the executable file created after go build command.
It can be run automatically by simply './main.exe'

2. ```Reusable``` : Code used as helpers. Good place to put reusable logic.

<a id="11-go-test"></a>
## <span style="color:#6495ED"> 11. Go Test
* Before running "go test" go.mod file required for the package so execute "go mod init <module_name>"
* To make a test, create a new file ending in _test.go
* To run all test in package, run command "go test"
* If static file to be used for testing use name "_testing"

<a id="12-concurrency"></a>
## <span style="color:#6495ED"> 12. Concurrency

* Concurrency is a concept in computer science and software development that refers to the ability of a program to 
execute multiple tasks or processes simultaneously. 

* It's not necessarily about running tasks in parallel on multiple processors (though that is a form of concurrency 
known as parallelism), but rather about making progress on multiple tasks over a short period.

* In Go, the concept of concurrency is central to its design. Go has goroutines (lightweight threads) and channels, 
which make it easy to write concurrent programs. 

* Goroutines allow functions to be executed concurrently, and channels provide a way for goroutines to communicate 
and synchronize their actions. This design makes it straightforward to express concurrent algorithms and enables 
the development of highly concurrent and scalable systems.


![!alt text](/images/go_package.png)

<a id="13-command-line-reference"></a>
## <span style="color:#6495ED"> 13. Command Line Reference

| Command | Description                                                |
|:------- |:-----------------------------------------------------------|
| go version | Check Go Lang version Installed                            |
| go help | Help Command                                               |
| go doc \<package> <br> go doc \<package> \<module> | Package or Module Document                                 |
| go run <main.go> | Compiles and Execute one or two files |
| go build src\\<dir_main.go> | Build Package and provide executable file .exe             |
| go build | Compiles a bunch of go source code files                   |
| go mod init <module_path> | - To intialize project as as Go module  or create a new module.<br> - Module path can correspond to a repository you plan to publish a module (eg: github.com/amolpratap-singh/booking-app).<br> -It initailize a go.mod file.<br> -Desrcibe the module: with name/module_path and go version used in the program.<br> -The module_path is also an import path (eg: import github.com/amolpratap-singh/booking-app)                      |
| go fmt | Formats all the code in each file in the current directory |
| go install | Compiles and install a pacakage                            |
| go get | Downloads the raw source code of someone else's package    |
| go test | Runs any test associated with the current project          |

<a id="14-core-package"></a>
## <span style="color:#6495ED"> 14. Core Package

### 14.1 Creating Packages and Documentation

Go introduced the concept of modules to manage dependencies and versioning.
Below is the method to create a math package to find the average of float number in a array:

1. Initialize project as a Go module and command to be executed inside project directory.

```bash
go mod init myproject
```
After executing above command we will get go.mod file which will contain package detail.

2. Project Structure has to be correct directory structure.

```text
myproject/
|- main.go
|- test/
   |-math/
     math.go 
```

3. Import paths

* main.go code:
```go
package main

import (
	"fmt"
	"myproject/test/math"
)

func main() {
	fmt.Println("I am main program")
	xs := []float64{1, 2, 3, 5, 5, 6, 7, 7}

	fmt.Println("Average of a number array", math.Average(xs))
}
```

* math.go code:
```go
package math

// Finds the average of a series of numbers
func Average(xs []float64) float64 {
	total := float64(0)

	for _, x := range xs {
		total += x
	}

	return total / float64(len(xs))
}
```

4. Build and Run 

```bash
go build main.go
```
After executing above command we will get executable file with extension .exe

5. Documentation 

```bash
go doc myproject/test/math Average
```
Above command will provide the documentation related the func as to reflect more info 
we can put up in the comment area. As output of above command provided below:

```text
package math // import "myproject/test/math"

func Average(xs []float64) float64
    Finds the average of a series of numbers
```

> **_NOTE:_** <span style="color:orange">Identifiers that start with a capital letter are exported (meaning accessible 
from other packages), whereas identifiers that start with a lowercase letter are not. For example Average func as 
mentioned above will be accessible but with small average it will not be accessible</span>


### 14.2 Strings

1. ```strings.Contains(s, substr string) boolean```: To search sub string in string.
2. ```strings.Count(s, sep string) int```: To get the count for the sub string present in the string.
3. ```strings.HasPrefix(s, prefix string) boolean```: To check the string is started with provided prefix string.
4. ```strings.HasSuffix(s, suffix string) boolean```: To check the string is end with the provided suffix string.
5. ```strings.Index(s, sep string) int```: To get the position of the provided sub string if it doesn't found it returns -1.
6. ```strings.Join(a [] string, sep string) string```: To take a list of string and join them together in a single 
string separate by character or another string (eg. "m", "-", ";").
7. ```strings.Repeat(s string, count int) string```: To repeat a sting.
8. ```strings.Replace(s, old, new string, n int ) string```: To replace old string present in the string with 
new string and 'n' number indicates for number of interval to replace the string. (n = -1) to do as many time as possible.
9. ```strings.Split(s, sep string) []string```: To split a string by a separate string or separators.
10. ```strings.ToLower(s string) string```: To convert a string to all in to lower case.
11. ```strings.ToUpper(s string) string```: To convert a string to all in to upper case.
12. To convert a string into slice of bytes or list of bytes to a string.

```go
arr := []byte("goLang")
str := string([]byte{'g', 'o', 'L', 'a', 'n', 'g'})
```

### 14.3 I/O [Need more information]

The io package consists of a few functions, but mostly interface used in other package. 
The two main interfaces are Reader and Writer. Reader supports reading via Read method and Writer supports writing 
via Write method. For example, the io package has a Copy function that copies data from a Reader to a Writer:

```go
func Copy(dst Writer, src Reader) (written int64, err error)
```

To read or write to a []byte or a string, you can use the Buffer struct found in the bytes package:

```go
var buf bytes.Buffer
buf.Write([]byte("test"))
```

### 14.4 File



### 14.5 Containers (Container/list)

The container/list package implements doubly linked list. Due to double linked list pointer will be available for 
previous node as well. Values are appended in the list using PushBack(). Following example for linked list :

```go
import ("fmt"; "container/list")

var linkedList list.List

// Insertion
linkedList.PushBack(1)
linkedList.PushBack(2)

// Access

for e := linkedList.Front(); e != nil; e=linkedList.Next() {
  fmt.Println(e.Value)
}
```

### 14.6 Sort [Need more Information]

The sort package contains to sort the arbitrary data

### 14.7 Hashes and Cryptography

A hash function takes a set of data and reduces it to a smaller fixed size. Hash function in Go are broken into two 
categories: cryptographic and non-cryptographic.

The non-cryptographic hash function can be found underneath the hash package and include adler32, crc32, crc64 and fnv

```go

```

### 14.8 time

In Go, the time package provides functionality for measuring and displaying time.

Here are some of the key components and functionalities of the time package:

1. Time Types:
   * ```time.Time```: Represents a moment in time with nanosecond precision.
   * ```time.Duration```: Represents a duration of time.
2. Constants:
   * ```time.Second, time.Minute, time.Hour```: Constants representing common time durations.
   * ```time.RFC3339, time.RFC822```: Constants defining standard time formats. 
3. Functions:
   * ```time.Now()```: Returns the current local time.
   * ```time.Parse(layout, value string)```: Parses a formatted string and returns the time value it represents.
   * ```time.Format(layout string)```: Formats a time value according to the provided layout.
   * ```time.Sleep(d Duration)```: Pauses the program execution for the specified duration. 
4. Operations: Arithmetic operations are supported on time.Time values, allowing you to add or subtract durations.

```go
package main

import (
"fmt"
"time"
)

func main() {

  // Get the current time
  now := time.Now()
  fmt.Println("Current Time:", now)

  // Format time as a string
  formattedTime := now.Format("2006-01-02 15:04:05")
  fmt.Println("Formatted Time:", formattedTime)
  
  // Parse a string into a time.Time
  parsedTime, err := time.Parse("2006-01-02 15:04:05", "2023-12-27 12:30:00")
  if err != nil {
    fmt.Println("Error parsing time:", err)
  } else {
	fmt.Println("Parsed Time:", parsedTime)
  }
  
  // Sleep for 2 seconds
  time.Sleep(2 * time.Second)
  
  // Add and subtract durations
  newTime := now.Add(24 * time.Hour)
  fmt.Println("Time 24 hours from now:", newTime)
  
  // Calculate the duration between two times
  duration := newTime.Sub(now)
  fmt.Println("Duration between times:", duration)
}

```

### 14.9 reflect (To find data type)

To find the type of a variable or expression in Go, you can use the reflect package. The reflect package provides 
runtime reflection capabilities, allowing you to inspect the type and value of variables dynamically.

```go
import(
    "fmt"
    "reflect"
)

x := 10
fmt.Println("x is type of", reflect.TypeOf(x))
```

### 14.10 Servers

### 14.11 TCP

In Go we can create a TCP server using the net package's Listen function. Listen takes a network type(eg: tcp) and an 
address with port to bind & returns a net.Listener.

<a id="15-library-example"></a>
## <span style="color:#6495ED"> 15. Library Example

* ```fmt```: The fmt package (shorthand for format) implements formating for input and output

Print
Println
Printf
Scanf
Sprintln
Sprintf

```sh
go doc fmt Println
```

```text
package fmt // import "fmt"

func Println(a ...any) (n int, err error)
Println formats using the default formats for its operands and writes to standard output. 
Spaces are always added between operands and a newline is appended. 
It returns the number of bytes written and any write error encountered.
```

```sh
go doc fmt Scanf
```

```text
package fmt // import "fmt"

func Scanf(format string, a ...any) (n int, err error)
Scanf scans text read from standard input, storing successive
space-separated values into successive arguments as determined by the
format. It returns the number of items successfully scanned. If that is less
than the number of arguments, err will report why. Newlines in the input
must match newlines in the format. The one exception: the verb %c always
scans the next rune in the input, even if it is a space (or tab etc.) or
newline.
```

```sh
go doc fmt Sprintf
```

```text
package fmt // import "fmt"

func Sprintf(format string, a ...any) string
    Sprintf formats according to a format specifier and returns the resulting
    string.
```

* ```math```:
* ```io```:
* ```os```:
* ```encoding```:
* ```crypto```:
* ```debug```:

<a id="16-reference-links"></a>
## <span style="color:#6495ED"> 16. Reference Links

Download Link https://go.dev/dl/
Standarad Pkg Link https://pkg.go.dev/std

### 16.1  Reference Books

* The Go Programming Language by Alan A. A. Donovan and Brian W. Kernighan
* Go in Action by William Kennedy, Brian Ketelsen, and Erik St. Martin
* Go Web Programming by Sau Sheong Chang
* Go Programming Blueprints by Mat Ryer
* Learning Go by Miek Gieben

<a id="17-qa"></a>
## <span style="color:#6495ED"> 17. QA

#### 1. How do we run the code in our project ?
#### 2. What does 'package main' mean ?
#### 3. What does 'import' "fmt" mean ?
#### 4. Difference between println() and fmt.Println() ?
```text
println() is a built-in function in Go and does not support formatting options. 
It also important to note that print and println report to stderr, not stdout.
fmt.Println() is part of the fmt package and provides more advanced formatting options. 
It allows you to format output with placeholders, control spacing, and other formatting options.
```


#### 5. What is that 'func' thing ?
#### 6. How is the main.go file organized ?

