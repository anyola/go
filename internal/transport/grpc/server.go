package grpctransport

import (
	"context"
	"errors"

	"movieapi/api"
	"movieapi/internal/movie"
	"movieapi/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	api.UnimplementedMovieServiceServer
	service *service.Service
}
func NewServer(s *service.Service) *Server { return &Server{service: s} }
func (s *Server) LoadMovies(ctx context.Context, req *api.LoadMoviesRequest) (*api.LoadMoviesResponse, error) {
	if req == nil || req.From <= 0 || req.To <= 0 || req.From > req.To || req.To-req.From > 10000 {
		return nil, status.Error(codes.InvalidArgument, "from and to must be positive, from <= to, range <= 10000")
	}
	movies, errs := s.service.LoadMovies(ctx, int(req.From), int(req.To))
	out := &api.LoadMoviesResponse{}
	for _, m := range movies { out.Movies = append(out.Movies, &api.Movie{ID: int32(m.ID), Title: m.Title, Year: int32(m.Year), Director: m.Director}) }
	for _, e := range errs { out.Errors = append(out.Errors, &api.LoadError{ID: int32(e.ID), Message: e.Message}) }
	return out, nil
}
func (s *Server) GetLastMovie(ctx context.Context, req *api.GetLastMovieRequest) (*api.Movie, error) {
	m, err := s.service.LastMovie(ctx)
	if err != nil {
		if errors.Is(err, movie.ErrNotFound) { return nil, status.Error(codes.NotFound, err.Error()) }
		if errors.Is(err, context.DeadlineExceeded) { return nil, status.Error(codes.DeadlineExceeded, err.Error()) }
		if errors.Is(err, context.Canceled) { return nil, status.Error(codes.Canceled, err.Error()) }
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return &api.Movie{ID: int32(m.ID), Title: m.Title, Year: int32(m.Year), Director: m.Director}, nil
}
