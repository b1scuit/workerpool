package workerpool_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/B1scuit/workerpool"
)

var mockError = errors.New("mock error")

// An example worker function
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

/*
func TestMain(m *testing.M) {
	client, _ = workerpool.New(&workerpool.ClientOptions{
		WorkerFunc: workerFunction,
	})

	m.Run()
}
*/

func BenchmarkWaitGroup(b *testing.B) {
	var workerFunction workerpool.TaskFunc = func(t *workerpool.Task) {
		t.Output <- Fib(40)
		close(t.Output)
	}

	var wg sync.WaitGroup

	client := workerpool.Must(workerpool.New(&workerpool.ClientOptions{
		Workers:    10,
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

	workerpool.Must(nil, mockError)
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