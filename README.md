# libgopy2

<!--
#field

#groups
Libraries

#languages
Go
Python

#frames and libs

-->

A golang library for simple python execution

## Notes

The library default compilation is done with `cgo pkg-config` using `python3-embed`, this should infer the python version for most platforms. If not, a `gen_flags.go` is provided

## Functions

Here a list of the exposed thread safe functions:
```go
Init() // initialize python interpreter
Load(module string) error // load a python module and all the contained functions (the functions starting wiht '_' will be considered private and then not loaded)
Call(fun string, args ...any) (any, error) // calls fun with args as args, returns nothing
Finalize() // finalize python interpreter
```

Those above are safe wrappers of the same functions which are faster but not suited for multithreaded applications
```go
InitUnsafe()
LoadUnsafe(module string) error
CallUnsafe(fun string, args ...any) (any, error)
FinalizeUnsafe()
```

## Types

List of the supported types for arguments:
```go
int64  int32  int16  int8  int
uint64 uint32 uint16 uint8 uint
float64 float32
bool
byte
rune
nil
[]byte []uint8 // as python bytes object
string
[]any
map[any]any
C.PyObject // a raw python object
```

List of the supported types for python return values:
```py
int
float
bool
str
bytes
list
tuple
dict
set
None
C.PyObject # unrecognized types are returned as raw python object
```

## Example

`test.py`
```py
def hello(x, y, z):
    print(f"Hello from test.py {x} {y} {z}")
    return b"Hello from test.py in bytes"
```

`main.go`
```go
package main

import (
    "fmt"

    py "github.com/tde-nico/libgopy2"
)

func main() {
	py.Init() // intitialize python

    err := py.Load("test") // loads the python module test (test.py)
    if err != nil {
        panic(err)
    }

    res := py.Call_byte("hello", "world", 5, 6.7) // calls fuction hello as: func hello(string, int, float) []byte
    fmt.Printf("Result: %s\n", res)

    py.Finalize() // finalize python
}
```

## Benchmarks

tested with: `go test -bench=. -benchtime=10s -count=10` on Windows 10 WSL 2
```bash
goos: linux
goarch: amd64
pkg: github.com/tde-nico/libgopy
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkLibgopy-8       4186684              2870 ns/op
BenchmarkLibgopy-8       4076690              2961 ns/op
BenchmarkLibgopy-8       4256122              2941 ns/op
BenchmarkLibgopy-8       4148332              2925 ns/op
BenchmarkLibgopy-8       4042467              2985 ns/op
BenchmarkLibgopy-8       4039228              2984 ns/op
BenchmarkLibgopy-8       4178586              2871 ns/op
BenchmarkLibgopy-8       4013100              2971 ns/op
BenchmarkLibgopy-8       3895118              2882 ns/op
BenchmarkLibgopy-8       4016323              2982 ns/op
PASS
ok      github.com/tde-nico/libgopy     139.135s
```

