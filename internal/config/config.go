package config

import "os"

type Config struct {
	MovieBaseURL string
	HTTPAddr     string
}

func Load() Config {
	baseURL := os.Getenv("MOVIE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://homeworksite.site"
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return Config{MovieBaseURL: baseURL, HTTPAddr: addr}
}
