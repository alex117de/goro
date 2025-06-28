// Package goro provides a safe wrapper for executing functions in goroutines.
// It handles panics and errors that might occur during goroutine execution.
package goro

import (
	"fmt"
	"runtime/debug"
)

// Func is a function type that can be executed in a goroutine.
// It returns a result of any type and an error if the execution fails.
// If there is no result to return, nil should be returned.
type Func func() (any, error)

// ErrorHandler is a function type that handles errors.
type ErrorHandler func(error)

// ResultHandler is a function type that handles results of any type.
type ResultHandler func(any)

// defaultErrorHandler is the default error handler that logs errors to stdout.
func defaultErrorHandler(err error) {
	fmt.Printf("Error in goroutine: %v\n", err)
}

// Goro represents a goroutine execution with error and result handling.
type Goro struct {
	fn            Func
	errHandler    ErrorHandler
	resultHandler ResultHandler
}

// New creates a new Goro instance with the provided function.
func New(f Func) *Goro {
	return &Goro{
		fn:         f,
		errHandler: defaultErrorHandler,
	}
}

// WithErrorHandler sets a custom error handler for the Goro instance.
func (g *Goro) WithErrorHandler(handler ErrorHandler) *Goro {
	g.errHandler = handler
	return g
}

// WithResultHandler sets a custom result handler for the Goro instance.
func (g *Goro) WithResultHandler(handler ResultHandler) *Goro {
	g.resultHandler = handler
	return g
}

// Start executes the function in a goroutine with the configured handlers.
func (g *Goro) Start() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic recovered: %v\nStack trace: %s", r, debug.Stack())
				g.errHandler(err)
			}
		}()

		result, err := g.fn()
		if err != nil {
			g.errHandler(err)
			return
		}

		// Only call the result handler if it's set and there's a non-nil result
		if g.resultHandler != nil && result != nil {
			g.resultHandler(result)
		}
	}()
}

// Go executes the provided function in a goroutine and handles any panics.
// If a panic occurs, it will be recovered and logged.
func Go(f Func) {
	New(f).Start()
}

// GoWithErrorHandler executes the provided function in a goroutine and handles any panics.
// If an error occurs, it will be passed to the provided error handler function.
// If a panic occurs, it will be recovered, converted to an error, and passed to the error handler.
func GoWithErrorHandler(f Func, errHandler ErrorHandler) {
	New(f).WithErrorHandler(errHandler).Start()
}

// GoWithResultHandler executes the provided function in a goroutine and handles any panics.
// If the function executes successfully, the result will be passed to the provided result handler function.
// If an error occurs, it will be passed to the provided error handler function.
// If a panic occurs, it will be recovered, converted to an error, and passed to the error handler.
func GoWithResultHandler(f Func, resultHandler ResultHandler, errHandler ErrorHandler) {
	New(f).WithResultHandler(resultHandler).WithErrorHandler(errHandler).Start()
}
