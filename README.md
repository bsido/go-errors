# go-errors

This library provides a simple way to create errors in Go 
that look like they were generated by the Rust compiler.

The below error message is from Rust:

```text
error[E0277]: `Foo` is not an iterator
 --> src/main.rs:4:16
  |
4 |     for foo in Foo {}
  |                ^^^ `Foo` is not an iterator
  |
  = note: maybe try calling `.iter()` or a similar method
  = help: the trait `std::iter::Iterator` is not implemented for `Foo`
  = note: required by `std::iter::IntoIterator::into_iter`

error[E0277]: `&str` is not an iterator
 --> src/main.rs:5:16
  |
5 |     for foo in "" {}
  |                ^^ `&str` is not an iterator
  |
  = note: call `.chars()` or `.bytes() on `&str`
  = help: the trait `std::iter::Iterator` is not implemented for `&str`
  = note: required by `std::iter::IntoIterator::into_iter`
```

Without modifying any of the defaults in this library we can achieve a similar error message in Go:

```go
package main

import (
    "github.com/bsido/go-errors/errors"
)

func rustCompilationError() error {
    return errors.New("'Foo' is not an iterator").
        Code(277).
        Causef(`src/main.rs:4:16

     for foo in Foo {}
                ^^^ 'Foo' is not an iterator`).
        Note("maybe try calling '.iter()' or a similar method").
        Help("the trait 'std::iter::Iterator' is not implemented for 'Foo'").
        Note("required by 'std::iter::IntoIterator::into_iter'").
        Wrap(errors.New("'&str' is not an iterator").
            Code(277).
            Causef(`src/main.rs:5:16

     for foo in "" {}
                ^^ '&str' is not an iterator`).
            Help("call '.chars()' or '.bytes() on '&str'").
            Help("the trait 'std::iter::Iterator' is not implemented for '&str'").
            Note("required by 'std::iter::IntoIterator::into_iter'"))
}
```

```text
error[E0277]: 'Foo' is not an iterator
  --> src/main.rs:4:16
   | 
   |      for foo in Foo {}
   |                 ^^^ 'Foo' is not an iterator
   | 
   = note: maybe try calling '.iter()' or a similar method
   = note: required by 'std::iter::IntoIterator::into_iter'
   = help: the trait 'std::iter::Iterator' is not implemented for 'Foo'

error[E0277]: '&str' is not an iterator
  --> src/main.rs:5:16
   | 
   | 	 for foo in "" {}
   | 				^^ '&str' is not an iterator
   | 
   = note: required by 'std::iter::IntoIterator::into_iter'
   = help: call '.chars()' or '.bytes() on '&str'
   = help: the trait 'std::iter::Iterator' is not implemented for '&str'
```

This library is not aimed at showing compilation errors only so it does not support source lines or column numbers by default.

## Usage

The following code snippets are also available in the `examples` directory.

Import the library:

```go
package main

import (
    "github.com/bsido/go-errors/errors"
)

// ....
```

If you would like to use the `errors` package from the standard library, it is recommended that you import it with an alias:

```go
package main

import (
    goerrors "errors"
	
    "github.com/bsido/go-errors/errors"
)

// ...
```

In general, you would want to use this library as follows:

```go
package main

import (
    "fmt"
	
    "github.com/thirdParty"
    "github.com/bsido/go-errors/errors"
)

func main() {
    if err := callFunction(); err != nil {
        // print the error
        fmt.Print(err)
    }
}

func callFunction() error {
    err := callThirdParty()
    if err != nil {
        // 1. set a reasonable message
        // 2. optional: add additional details to the message as the cause
        // 3. wrap the error from coming from the function
        // 4. optional: add help and/or note messages
        return errors.New("failed to execute our own function").
            // optional: otherwise the cause will be empty if we wrap the error
            // you may leave it empty if the above message is sufficient
            Causef("some additional details about the error"). 
            // we wrap the error because this is the second level of the call stack
            Wrap(err).
            Help("do this to fix the error") // optional
    }
    
    return nil
}

func callThirdParty() error {
    // call a third party library
    err := thirdParty.Exec()
	if err != nil {
        // 1. set a reasonable message
        // 2. set the error from the third party library as the cause
        return errors.New("failed to execute third party library").
            // we set the error as the cause because we know that this is not our error
            Cause(err) 
    }
	
    return nil
}
```

The final error message will look like this:

```text
error: failed to execute our own function
  --> some additional details about the error
   = help: do this to fix the error

error: failed to execute third party library
  --> something went wrong
```

### Basic usage

It is the same as if you were using the `errors` package from the standard library:

```go
errors.New("this is the main message")
```

```text
error: this is the main message
```

### Error code

```go
errors.New("this is the main message").
    Code(404)
```

```text
error[E0404]: this is the main message
```

### Cause

```go
errors.New("this is the main message").
    // alternatively, use Causef
    Cause("this is on the first line\nthis is on the second")
```

```text
error: this is the main message
  --> this is on the first line<br>
   | this is on the second
```

### Notes

```go
errors.New("this is the main message").
    // alternatively use Notef
    Note("this is the 1st note\nthis is the second line of the first note").
    Note("this is the 2nd note\nthis is the second line of the second note")
```

```text
error: this is the main message
   = note: this is the 1st note
           this is the second line of the first note
   = note: this is the 2nd note
           this is the second line of the second note
```

### Help

```go
errors.New("this is the main message").
    // alternatively use Helpf
    Help("this is the 1st help\nthis is the second line of the first help").
    Help("this is the 2nd help\nthis is the second line of the second help")
```

```text
error: this is the main message
   = help: this is the 1st help
           this is the second line of the first help
   = help: this is the 2nd help
           this is the second line of the second help
```

### Conditional help messages

Sometimes help messages should appear under certain conditions. This can be achieved by using the `HelpIf` method:

```go
if err != nil {
    return errors.New("this is the main message").
        Cause(err).
        HelpIf("seems like the thing that you are looking for was not found, do ... instead", func () bool {
            return strings.Contains(err.Error(), "not found")
        })
}
```

```text
error: this is the main message
  --> droids were not found
   = help: seems like the thing that you are looking for was not found, do ... instead
```

### Wrapping errors

```go
other := errors.New("bottom").
    Note("this might be happening because ...")

errors.New("top").
    Help("do this!").
    Wrap(other))
```

```text
error: top
   = help: do this!

error: bottom
   = note: this might be happening because ...
```

## Customization

TODO
