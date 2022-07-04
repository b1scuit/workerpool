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
	benchmarkWaitGroup(5, b)
}
func BenchmarkWorker10(b *testing.B) {
	benchmarkWaitGroup(10, b)
}
func BenchmarkWorker20(b *testing.B) {
	benchmarkWaitGroup(20, b)
}
func BenchmarkWorker50(b *testing.B) {
	benchmarkWaitGroup(50, b)
}

func benchmarkWaitGroup(num int, b *testing.B) {
	var workerFunction workerpool.TaskFunc = func(t *workerpool.Task) {
		t.Output <- Fib(10)
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
