package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"movieapi/internal/movie"
)

type MovieGetter interface {
	Get(context.Context, int) (*movie.Movie, error)
}

type Service struct {
	client MovieGetter
}

func New(client MovieGetter) *Service { return &Service{client: client} }

func (s *Service) LoadMovies(ctx context.Context, from, to int) ([]*movie.Movie, []movie.LoadError) {
	type result struct {
		movie *movie.Movie
		err   movie.LoadError
	}
	jobs := make(chan int)
	results := make(chan result)
	var wg sync.WaitGroup
	workers := to - from + 1
	if workers > 10 {
		workers = 10
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for id := range jobs {
				m, err := s.client.Get(ctx, id)
				r := result{movie: m}
				if err != nil {
					r.err = movie.LoadError{ID: id, Message: err.Error()}
				}
				select {
				case results <- r:
				case <-ctx.Done():
					return
				}
			}
		}()
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
	var movies []*movie.Movie
	var errs []movie.LoadError
	for r := range results {
		if r.movie != nil {
			movies = append(movies, r.movie)
		} else if r.err.ID != 0 {
			errs = append(errs, r.err)
		}
	}
	return movies, errs
}

func (s *Service) LastMovie(ctx context.Context) (*movie.Movie, error) {
	const maxID = 99999
	if m, err := s.client.Get(ctx, maxID); err == nil {
		return m, nil
	} else if !errors.Is(err, movie.ErrNotFound) {
		return nil, err
	}
	low, high := 1, maxID-1
	var last *movie.Movie
	for low <= high {
		mid := low + (high-low)/2
		m, err := s.client.Get(ctx, mid)
		if err == nil {
			last = m
			low = mid + 1
			continue
		}
		if !errors.Is(err, movie.ErrNotFound) {
			return nil, err
		}
		high = mid - 1
	}
	if last == nil {
		return nil, fmt.Errorf("%w: no movies available", movie.ErrNotFound)
	}
	return last, nil
}
