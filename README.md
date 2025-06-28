# goro - Safe Goroutine Wrapper

`goro` is a Go package that provides a safe wrapper for executing functions in goroutines. It handles panics and errors that might occur during goroutine execution, making concurrent programming in Go safer and more robust.

## Features

- Execute functions in goroutines with automatic panic recovery
- Handle errors returned by goroutines
- Provide custom error handlers for more control over error handling
- Object-oriented API with method chaining
- Support for result handling
- Simple and easy-to-use API with both functional and object-oriented approaches

## Installation

```bash
go get github.com/alex117de/goro
```

## Usage

### Object-Oriented API

The object-oriented API provides a flexible and chainable interface for working with goroutines.

#### Basic Usage

```text
// Create a new Goro instance and start it
goro.New(func() (any, error) {
    // Your code here
    return nil, nil
}).Start()
```

#### Error Handling

```text
// Define a custom error handler
errorHandler := func(err error) {
    log.Printf("Error occurred: %v", err)
    // You can also send the error to a monitoring service, etc.
}

// Create a new Goro instance with a custom error handler
goro.New(func() (any, error) {
    // Your code here
    return nil, errors.New("something went wrong")
}).WithErrorHandler(errorHandler).Start()
```

#### Result Handling

```text
// Define a result handler
resultHandler := func(result any) {
    // Type assertion to get the actual type
    if str, ok := result.(string); ok {
        log.Printf("Result: %s", str)
    } else {
        log.Printf("Result: %v", result)
    }
}

// Define an error handler
errorHandler := func(err error) {
    log.Printf("Error occurred: %v", err)
}

// Create a new Goro instance with result handling
goro.New(func() (any, error) {
    // Your code here
    return "operation completed", nil
}).WithResultHandler(resultHandler).WithErrorHandler(errorHandler).Start()
```

### Functional API

The functional API provides a simpler approach for working with goroutines.

#### Basic Usage

```text
// Execute a function in a goroutine with automatic panic recovery
goro.Go(func() (any, error) {
    // Your code here
    return nil, nil
})
```

#### Error Handling

```text
// The function can return an error
goro.Go(func() (any, error) {
    // Your code here
    return nil, errors.New("something went wrong")
})
```

#### Custom Error Handler

```text
// Define a custom error handler
errorHandler := func(err error) {
    log.Printf("Error occurred: %v", err)
    // You can also send the error to a monitoring service, etc.
}

// Use the custom error handler
goro.GoWithErrorHandler(func() (any, error) {
    // Your code here
    return nil, errors.New("something went wrong")
}, errorHandler)
```

## Examples

See the [functional example](./example/example.go) and [object-oriented example](./example/oo_example.go) for more detailed examples.

## License

MIT
