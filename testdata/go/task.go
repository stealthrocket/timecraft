//go:build wasip1

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stealthrocket/timecraft/sdk/go/timecraft"
)

func main() {
	var err error
	switch {
	case len(os.Args) == 2 && os.Args[1] == "worker":
		err = worker()
	case len(os.Args) == 1:
		err = supervisor(context.Background())
	default:
		err = fmt.Errorf("usage: task.wasm [worker]")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func supervisor(ctx context.Context) error {
	client, err := timecraft.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to timecraft: %w", err)
	}

	// Spawn the same WASM module, but with the "worker" arg.
	workerModule := timecraft.ModuleSpec{Args: []string{"worker"}}

	requests := []timecraft.TaskRequest{
		{
			Module: workerModule,
			Input: &timecraft.HTTPRequest{
				Method: "POST",
				Path:   "/foo",
				Headers: map[string][]string{
					"X-Foo": []string{"bar"},
				},
				Body: []byte("foo"),
			},
		},
		{
			Module: workerModule,
			Input: &timecraft.HTTPRequest{
				Method: "POST",
				Path:   "/bar",
				Headers: map[string][]string{
					"X-Foo": []string{"bar"},
				},
				Body: []byte("bar"),
			},
		},
	}

	taskIDs, err := client.SubmitTasks(ctx, requests)
	if err != nil {
		return fmt.Errorf("failed to submit tasks: %w", err)
	}

	taskRequests := map[timecraft.TaskID]*timecraft.HTTPRequest{}
	for i, taskID := range taskIDs {
		taskRequests[taskID] = requests[i].Input.(*timecraft.HTTPRequest)
	}

	tasks, err := client.PollTasks(ctx, len(requests), -1) // block until all tasks are complete
	if err != nil {
		return fmt.Errorf("failed to poll tasks: %w", err)
	}
	if len(tasks) != len(requests) {
		return fmt.Errorf("incorrect response from poll tasks: %#v", tasks)
	}

	for _, task := range tasks {
		if task.State != timecraft.Success {
			panic("task did not succeed")
		}
		res, ok := task.Output.(*timecraft.HTTPResponse)
		if !ok {
			panic("unexpected task output")
		}
		req, ok := taskRequests[task.ID]
		if !ok {
			panic("invalid task ID")
		}
		if res.StatusCode != 200 {
			panic("unexpected response code")
		} else if string(req.Body) != string(res.Body) {
			panic("unexpected response body")
		} else if res.Headers.Get("X-Timecraft") != "1" {
			panic("unexpected response headers")
		}
	}

	return client.DiscardTasks(ctx, taskIDs)
}

func worker() error {
	return timecraft.StartWorker(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("X-Foo") != "bar" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if string(body) != r.URL.Path[1:] {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("X-Timecraft", "1")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
}