## Fuzz Tests 

### Background
To help us catch edge cases, we use a strategy to generate inputs called [fuzzing](https://go.dev/security/fuzz/)

Fuzzing is native in Go 1.18, but only for the following input types:
> string, []byte

> int, int8, int16, int32/rune, int64

> uint, uint8/byte, uint16, uint32, uint64

> float32, float64

> bool

To test struct inputs, we use a library called [`go-fuzz-headers`](https://github.com/AdaLogics/go-fuzz-headers#projects-that-use-go-fuzz-headers)


### Running Fuzz Tests 
Since fuzzing involves generating inputs at time intervals, you explicitly run the tests with the `-fuzz` option.

Watch the output, and view any errors that may arise. When you are satisfied the tests confirm your code working as expected, you can exit the test. 

eg. 
` go test -fuzz=FuzzNewVersionManager -v`
