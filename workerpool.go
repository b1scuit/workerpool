package workerpool

type TaskFunc func(*Task)

type Task struct {
	Input  map[string]any
	Output chan any
}

type ClientOptions struct {
	Workers    int
	WorkerFunc TaskFunc
}

type Client struct {
	taskStack chan *Task
}

func worker(name int, taskQueue <-chan *Task, workerFunc TaskFunc) {
	for {
		workerFunc(<-taskQueue)
	}
}

func New(opts *ClientOptions) (*Client, error) {

	// Default to 3 workers
	if opts.Workers == 0 {
		opts.Workers = 3
	}

	taskStack := make(chan *Task)

	for i := 0; i <= opts.Workers; i++ {
		go worker(i, taskStack, opts.WorkerFunc)
	}

	return &Client{
		taskStack: taskStack,
	}, nil
}

func Must(client *Client, err error) *Client {
	if err != nil {
		panic(err)
	}

	return client
}

func (c *Client) Add(in *Task) error {

	// This exists in a goroutine so it's non blocking
	go func() {
		c.taskStack <- in
	}()

	return nil
}
