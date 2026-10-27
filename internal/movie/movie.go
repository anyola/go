package movie

import "fmt"

type Movie struct {
	ID       int    `json:"id" protobuf:"varint,1,opt,name=id,proto3"`
	Title    string `json:"title" protobuf:"bytes,2,opt,name=title,proto3"`
	Year     int    `json:"year" protobuf:"varint,3,opt,name=year,proto3"`
	Director string `json:"director" protobuf:"bytes,4,opt,name=director,proto3"`
}

type LoadError struct {
	ID      int    `json:"id" protobuf:"varint,1,opt,name=id,proto3"`
	Message string `json:"message" protobuf:"bytes,2,opt,name=message,proto3"`
}

func (e LoadError) Error() string { return fmt.Sprintf("movie %d: %s", e.ID, e.Message) }
