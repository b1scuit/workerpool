package workerpool_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/B1scuit/workerpool"
)

var errMock = errors.New("mock error")

// An example worker function
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

func BenchmarkWorker1(b *testing.B) {
	benchmarkWaitGroup(1, b)
}
func BenchmarkWorker2(b *testing.B) {
	benchmarkWaitGroup(2, b)
}
func BenchmarkWorker5(b *testing.B) {
	benchmarkWaitGroup(3, b)
}
func BenchmarkWorker10(b *testing.B) {
	benchmarkWaitGroup(4, b)
}
func BenchmarkWorker20(b *testing.B) {
	benchmarkWaitGroup(5, b)
}
func BenchmarkWorker50(b *testing.B) {
	benchmarkWaitGroup(10, b)
}

func benchmarkWaitGroup(num int, b *testing.B) {
	var workerFunction workerpool.TaskFunc = func(t *workerpool.Task) {
		t.Output <- Fib(20)
		close(t.Output)
	}

	var wg sync.WaitGroup

	client := workerpool.Must(workerpool.New(&workerpool.ClientOptions{
		Workers:    num,
		WorkerFunc: workerFunction,
	}))

	for i := 0; i <= b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			output := make(chan any)

			client.Add(&workerpool.Task{
				Output: output,
			})

			<-output
		}(&wg)
	}

	wg.Wait()

}

func TestClient(t *testing.T) {
	client, err := workerpool.New(&workerpool.ClientOptions{})

	if err != nil {
		t.Error(err)
		return
	}

	if client == nil {
		t.Error("client nil")
	}

}

func TestMust(t *testing.T) {
	client := workerpool.Must(&workerpool.Client{}, nil)
	if client == nil {
		t.Error("nil client, should have paniced")
	}
}

func TestMustPanic(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()

	workerpool.Must(nil, errMock)
}

func TestAddTask(t *testing.T) {
	var workerFunction workerpool.TaskFunc = func(t *workerpool.Task) {
		t.Output <- true // "successfully" complete any task
	}

	client := workerpool.Must(workerpool.New(&workerpool.ClientOptions{
		Workers:    1,
		WorkerFunc: workerFunction,
	}))

	client.Add(&workerpool.Task{})
}

func Example() {

	var basicWorkerFunction = func(t *workerpool.Task) {
		t.Output <- true // Just mark the task complete
	}

	client := workerpool.Must(workerpool.New(&workerpool.ClientOptions{
		Workers:    5,                   // The number of desired workers
		WorkerFunc: basicWorkerFunction, // The desired worker function
	}))

	// Add a task to the queue
	output := make(chan any, 1)
	client.Add(&workerpool.Task{
		Input:  nil,
		Output: output,
	})

	// Be careful here as there is nothing blocking the application exit
	// you will need to hold the main execution path open until all tasks
	// have completed, this is not an issue around a persistant service
	// however if you are using this in a non persistant application
	// you will need to use a sync.WaitGroup (See other example)
}

// In instances where you need to retrieve the responses, this can be used
// with a waitgroup
func Example_withWaitgroup() {
	var basicWorkerFunction = func(t *workerpool.Task) {
		t.Output <- true // Just mark the task complete
	}

	var taskList []string

	client := workerpool.Must(workerpool.New(&workerpool.ClientOptions{
		Workers:    5,                   // The number of desired workers
		WorkerFunc: basicWorkerFunction, // The desired worker function
	}))

	var wg sync.WaitGroup
	for _, task := range taskList {
		wg.Add(1)

		// For each task in the list, add to the pool and wait for a response
		go func(wg *sync.WaitGroup, data string) {
			defer wg.Done()

			output := make(chan any, 1)

			client.Add(&workerpool.Task{
				Input:  data,
				Output: output,
			})

			<-output // Do something with this value if needed
		}(&wg, task)
	}

	wg.Wait() // Wait for all tasks to complete
}
