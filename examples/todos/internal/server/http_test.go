package server_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"todo/internal/server"
	"todo/internal/todo"
)

// loggingTransport is a [http.RoundTripper] which logs every outgoing request
type loggingTransport struct {
	t    *testing.T
	next http.RoundTripper
}

func newLoggingTransport(t *testing.T) *loggingTransport {
	return &loggingTransport{
		t:    t,
		next: http.DefaultTransport,
	}
}

func (l *loggingTransport) RoundTrip(request *http.Request) (r *http.Response, err error) {
	log.Printf("request: %s %s %v", request.Method, request.URL, request.Header)
	defer func() {
		log.Printf("response: %s %v", r.Status, r.Header)
	}()
	return l.next.RoundTrip(request)
}

// compile-time guarantee, that *testStorage implements Storage interface
var _ todo.Storage = &testStorage{}

type testStorage struct {
	tasks map[todo.ID]todo.Task
}

func newTestStorage() *testStorage {
	return &testStorage{tasks: make(map[todo.ID]todo.Task)}
}

func (s *testStorage) Upsert(ctx context.Context, t todo.Task) (stored todo.Task, err error) {
	s.tasks[t.ID] = t
	return t, nil
}

func (s *testStorage) List(ctx context.Context, f *todo.TaskFilter) ([]todo.Task, error) {
	out := make([]todo.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		out = append(out, task)
	}

	return out, nil
}

// integration-like test for the 'backend' API, which spins-up an actual server
func Test_APIHandler(t *testing.T) {
	s := newTestStorage()
	h := todo.NewHandler(s)
	api := must(server.NewHttp(nil, h, s))

	srv := httptest.NewServer(api.APIHandler())
	client := srv.Client()
	client.Transport = newLoggingTransport(t)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	formData := url.Values{"todo": {"test-todo-1"}}
	resp, err := client.Post(srv.URL+"/todos", "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status: %d, actual: %d", http.StatusCreated, resp.StatusCode)
	}

	loc := must(resp.Location())

	if loc.String() != srv.URL+"/" {
		t.Errorf("expected to be redirected to root server url, actual location: %s", loc)
	}
}

func mustT[T any](t *testing.T) func(value T, err error) T {
	return func(value T, err error) T {
		if err != nil {
			t.Fatal(err)
		}

		return value
	}
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

type mustFunc[T any] func(T, error) T
