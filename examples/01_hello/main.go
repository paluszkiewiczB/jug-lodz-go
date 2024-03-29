package main

import "fmt"

// main is the entry point for the application.
// Build the application with: `go build .`
// Run the application with: `go run .`
// Be aware that the app by default **can require** a dynamic binary (depends on the dependencies) - on Linux it could be glibc.
// You can check it with `ldd` command.
// If it is, and you want to distribute the binary as a Docker image built `FROM scratch` (or Alpine Linux), you need to build a static binary.
// To build a static binary use: `CGO_ENABLED=0 go build .`
//
// More advanced:
// To reduce the binary size use: `go build -ldflags="-s -w" .`
// -w: This flag omits the DWARF symbol table, effectively removing debugging information.
// -s: This strips the symbol table and debug information from the binary.
// For further size reduction see: https://github.com/xaionaro/documentation/blob/master/golang/reduce-binary-size.md
func main() {
	fmt.Println("Hello, World!")
}
