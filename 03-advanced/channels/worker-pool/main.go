// Topic: worker-pool | Concept: channels | Stage: 03-advanced
// Run:   go run ./03-advanced/channels/worker-pool
// Race:  go run -race ./03-advanced/channels/worker-pool
//
// Explore:
// * A fixed number of workers reading from one jobs channel
// * Results come back on a results channel
// * Close jobs to tell workers to stop; close results after all workers end
// * Limits concurrency (e.g. max 3 HTTP calls or DB queries at once)
// * Errors travel inside the result struct
// * Results arrive in completion order: sort or index them if order matters
package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Input int
}

type Result struct {
	JobID    int
	WorkerID int
	Output   int
	Err      error
}

func worker(id int, jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		time.Sleep(10 * time.Millisecond) // simulate I/O
		r := Result{JobID: j.ID, WorkerID: id}
		if j.Input < 0 {
			r.Err = fmt.Errorf("job %d: negative input %d", j.ID, j.Input)
		} else {
			r.Output = j.Input * 2
		}
		results <- r
	}
}

func main() {
	fmt.Println("Topic: worker pool")

	inputs := []int{3, 7, -1, 10, 4, 8, 2, 6, 5}
	const numWorkers = 3

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Go(func() { worker(w, jobs, results) })
	}

	// Feed jobs, then close so workers exit their range loop
	go func() {
		for i, in := range inputs {
			jobs <- Job{ID: i + 1, Input: in}
		}
		close(jobs)
	}()

	// Close results once every worker is done
	go func() {
		wg.Wait()
		close(results)
	}()

	start := time.Now()
	var all []Result
	for r := range results {
		all = append(all, r)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].JobID < all[j].JobID })

	for _, r := range all {
		if r.Err != nil {
			fmt.Printf("job %d: ERROR %v\n", r.JobID, r.Err)
			continue
		}
		fmt.Printf("job %d: output %2d\n", r.JobID, r.Output)
	}
	fmt.Printf("%d jobs x 10ms with %d workers took about %dms\n",
		len(inputs), numWorkers, time.Since(start).Round(10*time.Millisecond).Milliseconds())
}
