package main

import (
	"errors"
	"fmt"
	"time"

	"goro"
)

func main() {
	// Example 1: Simple goroutine with error handling
	goro.Go(func() (any, error) {
		fmt.Println("Running task 1...")
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		return nil, nil // No result, no error
	})

	// Example 2: Goroutine that returns an error
	goro.Go(func() (any, error) {
		fmt.Println("Running task 2...")
		// Simulate some work
		time.Sleep(200 * time.Millisecond)
		return nil, errors.New("task 2 failed")
	})

	// Example 3: Goroutine that panics
	goro.Go(func() (any, error) {
		fmt.Println("Running task 3...")
		// Simulate some work
		time.Sleep(300 * time.Millisecond)
		// This will panic but will be recovered
		panic("something went wrong in task 3")
	})

	// Example 4: Using custom error handler
	errorHandler := func(err error) {
		fmt.Printf("Custom error handler received: %v\n", err)
	}

	goro.GoWithErrorHandler(func() (any, error) {
		fmt.Println("Running task 4...")
		// Simulate some work
		time.Sleep(400 * time.Millisecond)
		return nil, errors.New("task 4 failed")
	}, errorHandler)

	// Example 5: Using custom error handler with panic
	goro.GoWithErrorHandler(func() (any, error) {
		fmt.Println("Running task 5...")
		// Simulate some work
		time.Sleep(500 * time.Millisecond)
		panic("something went wrong in task 5")
	}, errorHandler)

	// Example 6: Using result handler with successful result
	resultHandler := func(result any) {
		// Type assertion to get the actual type
		if str, ok := result.(string); ok {
			fmt.Printf("Result handler received: %s\n", str)
		} else {
			fmt.Printf("Result handler received unknown type: %v\n", result)
		}
	}

	goro.GoWithResultHandler(func() (any, error) {
		fmt.Println("Running task 6...")
		// Simulate some work
		time.Sleep(600 * time.Millisecond)
		return "task 6 completed successfully", nil
	}, resultHandler, errorHandler)

	// Example 7: Using result handler with error
	goro.GoWithResultHandler(func() (any, error) {
		fmt.Println("Running task 7...")
		// Simulate some work
		time.Sleep(700 * time.Millisecond)
		return "", errors.New("task 7 failed")
	}, resultHandler, errorHandler)

	// Example 8: Using result handler with panic
	goro.GoWithResultHandler(func() (any, error) {
		fmt.Println("Running task 8...")
		// Simulate some work
		time.Sleep(800 * time.Millisecond)
		panic("something went wrong in task 8")
	}, resultHandler, errorHandler)

	// Wait to see the results
	time.Sleep(2 * time.Second)
	fmt.Println("All tasks have been started. Some may still be running.")
}
