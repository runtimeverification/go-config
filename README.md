# Overview

This package allows the user to define a struct such as the following:
```go
type MyConfiguration struct {
  MyStringVar string `cli:"my_string_var" env:"MY_STRING_VAR"`
}
```
which allows the field `MyStringVar` to be populated either
by passing `--my_string_var SOME_VALUE` on the command line,
or from the `MY_STRING_VAR` environment variable.

For example, the following program
```go
import (
  "fmt"
  config "github.com/runtimeverification/go-config"
)

type MyConfiguration struct {
  MyStringVar string `cli:"my_string_var" env:"MY_STRING_VAR"`
}

func main() {
  var c MyConfiguration

  err := Init(&c)

  if err != nil {
    fmt.Println("Some error ", err)
    return
  }

  fmt.Println(c.MyStringVar)
}
```
will print `SOME_VALUE` with either
* `MY_STRING_VAR=SOME_VALUE go run .` or
* `go run . --my_string_var=SOME_VALUE`


# Field Types

At the moment, only `string`, `int`, and `bool` are supported.


# Tags

| Tag | Description | Priority |
| --- | ----------- | -------- |
| `cli:"cli_arg"` | You can pass `--cli_arg=value` on the command line. | 1 |
| `env:"ENV_VAR"` | You can set `ENV_VAR=value`. | 2 |
| `default:"VALUE"` | For when no command-line argument or environment variable is found. | 3 |
| `description:"DESC"` | For when `--help` is used. | N/A |

Notes:
* Command-line arguments override environment variables, which
  override default values.
*  The default values of the `default` tag are
   * `default:"0"` for `int`
   * `default:""` for `string`
   * `default:"false"` for `bool`.
* The default value of the `description` tag is the empty string.


# Implementation Details

Command-line argument parsing is done with the [flag](https://pkg.go.dev/flag)
package. For instance, the value of the `description` tag is passed to the
`flag.StringVar()` call (same for `IntVar()` and `BoolVar()`).

One implication is that command-line arguments are processed exactly as in the
`flag` [package specification](https://pkg.go.dev/flag).

Test coverage is currently at 90.6%.