A library for parsing unstructured data such as json to go types. Parsing is done by implementing `Parse<type>`
interfaces on custom types. 

## Short example

```go
package main

import (
	"fmt"
	"errors"
	"github.com/branchgrove/kaeru"
)

type Person struct {
	Name Name `parse:"name"`
	CreatedErlang bool `parse:"createdErlang"`
}

type Name string

func (n *Name) ParseString(s string) error {
  if len(s) == 0 || len(s) > 64 {
      return errors.New("Invalid name must be at least one character and at most 64 characters")
  }

  *n = Name(s)

  return nil
}

func main() {
	data := []byte(`{ "name": "Joe Armstrong", "createdErlang": true }`)
	
	person := new(Person)
	
	if err := kaeru.ParseJsonBytes(data, person); err != nil {
		fmt.Printf("Failed to parse %s\n", err)
	}
}
```

## Why?

kaeru follows the spirit of [Parse, don't validate] as an alternative to packages like [go-playground/validator].
Problems such as accidentally passing the wrong value as an argument to a function or assuming that a `string` is
valid when accepting it as a parameter can be mitigated with domain types. A common saying for passing around values
as `string`s is that something is "stringly typed" which kaeru tries to resolve with easier mapping from
unstructured data to domain types without the need to parse `interface{}` directly.

[go-playground/validator]: https://github.com/go-playground/validator
[Parse, don't validate]: https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/
