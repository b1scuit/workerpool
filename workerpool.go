package workerpool

// A basic type for the package to refer to for a task function
type TaskFunc func(*Task)

// A very basic type for a task IO
type Task struct {
	Input  any      // The data provided to the worker task function
	Output chan any // The "return address" for when the worker func is completed
}

type ClientOptions struct {
	Workers    int      // The number of workers to use
	WorkerFunc TaskFunc // The "task" function the workers should process tasks with
}

// A base coordinating client to manage IO for this package
type Client struct {
	taskStack chan *Task
}

// The worker that will pick tasks up, complete them and rerun
//
// I've included a name var here so you can output which worker
// is listeneing, i've commented out the line for speed but if
// you want more visability on whats going on
func worker(name int, taskQueue <-chan *Task, workerFunc TaskFunc) {
	for {
		// log.Printf("Worker %v listening", name)

		// Instead of passing a worker function, you can just put
		// whatever needs to be done in this loop, I use a worker function
		// as it makes this library more portable
		workerFunc(<-taskQueue)
	}
}

func New(opts *ClientOptions) (*Client, error) {

	// Default to 3 workers
	if opts.Workers == 0 {
		opts.Workers = 3
	}

	// Create that paper stack of tasks to be done
	taskStack := make(chan *Task)

	// Spin up as many workers and tell them the task function they have to complete
	// and the task stack to listen to
	for i := 0; i <= opts.Workers; i++ {
		go worker(i, taskStack, opts.WorkerFunc)
	}

	// Return a client to have the caller use the package
	return &Client{
		taskStack: taskStack,
	}, nil
}

// Very simple Must pattern function
func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

// Add adds a task safely to the task stack for the workers to process
// if you want to listen to the response, listen to the in.Output external to this
// (see examples)
func (c *Client) Add(in *Task) error {

	// This exists in a goroutine so it's non blocking
	go func() {
		c.taskStack <- in
	}()

	return nil
}
