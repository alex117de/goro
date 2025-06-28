package goro

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	f := func() (any, error) { return nil, nil }
	g := New(f)

	if g == nil {
		t.Fatal("New returned nil")
	}
	if g.fn == nil {
		t.Error("New did not set fn")
	}
	if g.errHandler == nil {
		t.Error("New did not set default error handler")
	}
}

// TestWithErrorHandler tests the WithErrorHandler method
func TestWithErrorHandler(t *testing.T) {
	f := func() (any, error) { return nil, nil }
	customHandler := func(err error) {}
	g := New(f).WithErrorHandler(customHandler)

	if g.errHandler == nil {
		t.Error("WithErrorHandler did not set error handler")
	}
	// Verify it's the same instance (fluent interface)
	if g.fn == nil {
		t.Error("WithErrorHandler did not preserve fn")
	}
}

// TestWithResultHandler tests the WithResultHandler method
func TestWithResultHandler(t *testing.T) {
	f := func() (any, error) { return "result", nil }
	customHandler := func(result any) {}
	g := New(f).WithResultHandler(customHandler)

	if g.resultHandler == nil {
		t.Error("WithResultHandler did not set result handler")
	}
	// Verify it's the same instance (fluent interface)
	if g.fn == nil {
		t.Error("WithResultHandler did not preserve fn")
	}
}

// TestStartWithoutResult tests the Start method with functions that don't return results
func TestStartWithoutResult(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	f := func() (any, error) {
		defer wg.Done()
		executed = true
		return nil, nil
	}

	New(f).Start()

	// Wait for the goroutine to complete
	wg.Wait()

	if !executed {
		t.Error("Function was not executed")
	}
}

// TestStartWithoutResultError tests error handling in Start method
func TestStartWithoutResultError(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedErr := errors.New("test error")
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	f := func() (any, error) {
		return nil, expectedErr
	}

	New(f).WithErrorHandler(errHandler).Start()

	// Wait for the error handler to be called
	wg.Wait()

	if !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, receivedErr)
	}
}

// TestStartWithResult tests the Start method with functions that return results
func TestStartWithResult(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedResult := "test result"
	var receivedResult any

	resultHandler := func(result any) {
		receivedResult = result
		wg.Done()
	}

	f := func() (any, error) {
		return expectedResult, nil
	}

	New(f).WithResultHandler(resultHandler).Start()

	// Wait for the result handler to be called
	wg.Wait()

	if receivedResult != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, receivedResult)
	}
}

// TestStartWithResultError tests error handling in the Start method with result functions
func TestStartWithResultError(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedErr := errors.New("test error")
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	resultHandler := func(result any) {
		t.Error("Result handler should not be called when error occurs")
	}

	f := func() (any, error) {
		return nil, expectedErr
	}

	New(f).WithResultHandler(resultHandler).WithErrorHandler(errHandler).Start()

	// Wait for the error handler to be called
	wg.Wait()

	if !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, receivedErr)
	}
}

// TestStartWithResultPanic tests panic recovery in Start method with result functions
func TestStartWithResultPanic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	panicMsg := "test panic"
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	resultHandler := func(result any) {
		t.Error("Result handler should not be called when panic occurs")
	}

	f := func() (any, error) {
		panic(panicMsg)
	}

	New(f).WithResultHandler(resultHandler).WithErrorHandler(errHandler).Start()

	// Wait for the error handler to be called
	wg.Wait()

	if receivedErr == nil {
		t.Error("Error handler was not called with panic error")
	}

	if !strings.Contains(receivedErr.Error(), panicMsg) {
		t.Errorf("Error does not contain panic message. Got: %v", receivedErr)
	}
}

// TestStartWithoutResultPanic tests panic recovery in Start method without result functions
func TestStartWithoutResultPanic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	panicMsg := "test panic"
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	f := func() (any, error) {
		panic(panicMsg)
	}

	New(f).WithErrorHandler(errHandler).Start()

	// Wait for the error handler to be called
	wg.Wait()

	if receivedErr == nil {
		t.Error("Error handler was not called with panic error")
	}

	if !strings.Contains(receivedErr.Error(), panicMsg) {
		t.Errorf("Error does not contain panic message. Got: %v", receivedErr)
	}
}

// TestStartWithResultNoHandler tests that Start does not panic when no result handler is provided
// but still executes the function
func TestStartWithResultNoHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	f := func() (any, error) {
		defer wg.Done()
		executed = true
		return "result", nil
	}

	// This should not panic now
	New(f).Start()

	// Wait for the goroutine to complete
	wg.Wait()

	if !executed {
		t.Error("Function was not executed")
	}
}

// TestGo tests the Go function (backward compatibility)
func TestGo(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	f := func() (any, error) {
		defer wg.Done()
		executed = true
		return nil, nil
	}

	Go(f)

	// Wait for the goroutine to complete
	wg.Wait()

	if !executed {
		t.Error("Function was not executed")
	}
}

// TestGoWithErrorHandler tests the GoWithErrorHandler function
func TestGoWithErrorHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedErr := errors.New("test error")
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	f := func() (any, error) {
		return nil, expectedErr
	}

	GoWithErrorHandler(f, errHandler)

	// Wait for the error handler to be called
	wg.Wait()

	if !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, receivedErr)
	}
}

// TestGoWithResultHandler tests the GoWithResultHandler function
func TestGoWithResultHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedResult := "test result"
	var receivedResult any

	resultHandler := func(result any) {
		receivedResult = result
		wg.Done()
	}

	errHandler := func(err error) {
		t.Error("Error handler should not be called when no error occurs")
	}

	f := func() (any, error) {
		return expectedResult, nil
	}

	GoWithResultHandler(f, resultHandler, errHandler)

	// Wait for the result handler to be called
	wg.Wait()

	if receivedResult != expectedResult {
		t.Errorf("Expected result %v, got %v", expectedResult, receivedResult)
	}
}

// TestGoWithResultHandlerError tests error handling in GoWithResultHandler function
func TestGoWithResultHandlerError(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	expectedErr := errors.New("test error")
	var receivedErr error

	errHandler := func(err error) {
		receivedErr = err
		wg.Done()
	}

	resultHandler := func(result any) {
		t.Error("Result handler should not be called when error occurs")
	}

	f := func() (any, error) {
		return nil, expectedErr
	}

	GoWithResultHandler(f, resultHandler, errHandler)

	// Wait for the error handler to be called
	wg.Wait()

	if !errors.Is(receivedErr, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, receivedErr)
	}
}

// TestDefaultErrorHandler tests the default error handler
func TestDefaultErrorHandler(t *testing.T) {
	// This test is mainly for coverage, as the default handler just prints to stdout
	err := errors.New("test error")
	defaultErrorHandler(err) // Should not panic
}
