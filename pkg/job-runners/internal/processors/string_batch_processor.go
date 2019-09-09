package processors

import (
	"errors"
	"github.com/mmacheerpuppy/GoPlayground/pkg/job-runners/internal/interfaces/runners"
	"gopkg.in/tomb.v1"
	"sync"
	"time"
)

// StringBatchProcessor is responsible for providing a runtime for concurrently batch processing
// ToStringTransformations and returning the results.
type StringBatchProcessor struct {
	jobs       []runners.ToStringJob
	primedJobs []func(<-chan struct{}) (string, error)
	mu sync.Mutex
	t  tomb.Tomb
	d time.Duration
}

func NewBatchStringProcessor(d time.Duration) *StringBatchProcessor {
	var mu sync.Mutex
	var t tomb.Tomb
	return &StringBatchProcessor{
		jobs: []runners.ToStringJob{},
		primedJobs: []func(<-chan struct{}) (string, error){},
		mu: mu,
		t:  t,
		d: d,
	}
}

// AddJob a new process.
func (p *StringBatchProcessor) AddJob(job runners.ToStringJob) []runners.ToStringJob {
	defer p.mu.Unlock()
	p.mu.Lock()

	p.jobs = append(p.jobs, job)
	return p.jobs
}

// AddJob a slice of processes.
func (p *StringBatchProcessor) AddJobs(transformations []runners.ToStringJob) []runners.ToStringJob {
	defer p.mu.Unlock()
	p.mu.Lock()

	p.jobs = append(p.jobs, transformations...)
	return p.jobs
}

// Start runs all the given functions concurrently until either they all complete or one returns an error, in which case it returns that error.
// The functions are passed a channel which will be closed when the function should stop.
func (p StringBatchProcessor) Run() ([]string, error) {
	p.prime()
	var wg sync.WaitGroup
	defer p.mu.Unlock()
	p.mu.Lock()
	if len(p.primedJobs) == 0 {
		return []string{}, nil
	}

	allDone := make(chan struct{})
	results := make([]string, len(p.primedJobs))
	for index, toExecute := range p.primedJobs {
		wg.Add(1)
		go func(index int, functionToExecute func(<-chan struct{}) (string, error)) {
			defer wg.Done()
			result, err := functionToExecute(p.t.Dying())
			if err != nil {
				p.t.Kill(errors.New(string(index) + ": " + err.Error()))
			}
			results[index] = result
		}(index, toExecute)
	}

	// Start a goroutine to wait for every process to finish.
	go func() {
		wg.Wait()
		close(allDone)
	}()

	// Wait for them all to finish, or one to
	select {
	case <-allDone:
	case <-p.t.Dying():
	}
	p.t.Done()
	return results, p.t.Err()
}

// Prime wraps the runtime jobs in a function which provides a channel that can be used to track the results.
func (p *StringBatchProcessor) prime() {
	defer p.mu.Unlock()
	p.mu.Lock()
	p.primedJobs = make([]func(<-chan struct{}) (string, error), len(p.jobs))
	for i := 0; i <= len(p.primedJobs)-1; i++ {
		job := p.jobs[i]
		primedJob := func(stop <-chan struct{}) (string, error) {
			result, err := job()
			if err != nil {
				return "", err
			} else {
				return result, nil
			}
		}
		p.primedJobs[i] = primedJob
	}
}