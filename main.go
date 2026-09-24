package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type Result struct {
	Movie *Movie
	ID    int
	Err   error
}

func fetchMovie(ctx context.Context, client *http.Client, id int, timeout time.Duration) (*Movie, error) {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		requestCtx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	var movie Movie

	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	return &movie, nil
}

func worker(
	ctx context.Context,
	jobs <-chan int,
	results chan<- Result,
	client *http.Client,
	timeout time.Duration,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case id, ok := <-jobs:
			if !ok {
				return
			}

			movie, err := fetchMovie(ctx, client, id, timeout)

			result := Result{
				ID:  id,
				Err: err,
			}

			if err == nil {
				result.Movie = movie
			}

			select {
			case results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	var from int
	var to int
	var workers int
	var timeout time.Duration

	flag.IntVar(&from, "from", 0, "ID of the first movie")
	flag.IntVar(&to, "to", 0, "ID of the last movie")
	flag.IntVar(&workers, "workers", 10, "number of workers")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "HTTP request timeout")

	flag.Parse()

	if from <= 0 || to <= 0 {
		fmt.Fprintln(os.Stderr, "error: --from and --to are required and must be greater than 0")
		os.Exit(1)
	}

	if from > to {
		fmt.Fprintln(os.Stderr, "error: --from must be less than or equal to --to")
		os.Exit(1)
	}

	if workers <= 0 {
		fmt.Fprintln(os.Stderr, "error: --workers must be greater than 0")
		os.Exit(1)
	}

	if timeout <= 0 {
		fmt.Fprintln(os.Stderr, "error: --timeout must be greater than 0")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	client := &http.Client{}

	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker(ctx, jobs, results, client, timeout, &wg)
	}

	go func() {
		defer close(jobs)

		for id := from; id <= to; id++ {
			select {
			case jobs <- id:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.Err != nil {
			if errors.Is(result.Err, context.Canceled) {
				continue
			}

			fmt.Fprintf(os.Stderr, "movie %d: %v\n", result.ID, result.Err)
			continue
		}

		fmt.Printf(
			"%d — %s — %d — %s\n",
			result.Movie.ID,
			result.Movie.Title,
			result.Movie.Year,
			result.Movie.Director,
		)
	}
}
