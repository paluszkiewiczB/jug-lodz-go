# JUG Łódź - It's timeto Go @ Docaposte

## Prerequisites

- [Go 1.21+](https://go.dev/doc/install)
- [Task v3.33.1+](https://taskfile.dev/installation/#go-modules) (optional, you can write commands yourself)


## Usage

1. Make 'examples' your working directory
2. Install the dependencies
    ```shell
    go mod download
    ```
3. Run specific example
    ```shell
    go run ./01_hello
    go test ./01_hello
    ```
4. Look for more specific commands in comments within the example

