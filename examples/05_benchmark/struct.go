package data

import "io"

type Empty struct{}

//go:noinline
func NewEmpty() io.Writer {
	return &Empty{}
}

var _ io.Writer = &Empty{}

func (e *Empty) Write(_ []byte) (int, error) {
	return 0, nil
}

type NotEmpty struct {
	Name   string
	Buffer [1024]byte
}

//go:noinline
func NewNotEmpty(name string) io.Writer {
	return &NotEmpty{Name: name}
}

var _ io.Writer = &NotEmpty{}

func (e *NotEmpty) Write(_ []byte) (int, error) {
	return len(e.Buffer), nil
}
