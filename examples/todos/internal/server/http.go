package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"time"
	"todo/internal/todo"
)

type Http struct {
	srv *http.Server
	ui  *UI
	s   todo.Storage
	h   *todo.Handler
}

type HttpCfg struct {
	Address string
}

var (
	defaultCfg = HttpCfg{
		Address: ":3456",
	}
)

func NewHttp(cfg *HttpCfg, handler *todo.Handler, storage todo.Storage) (*Http, error) {
	c := defaultCfg
	if cfg != nil {
		c = *cfg
	}

	ui, err := NewUI()
	if err != nil {
		return nil, fmt.Errorf("creating the UI: %w", err)
	}

	srv := &Http{ui: ui, s: storage, h: handler}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static", srv.StaticHandler()))
	mux.Handle("/api/", http.StripPrefix("/api", srv.APIHandler()))
	mux.Handle("/", srv.UIHandler())

	srv.srv = &http.Server{
		Addr:    c.Address,
		Handler: logRequest(closeBody(mux)),
	}

	return srv, nil
}

// Start blocks until the server is stopped.
// It can be stopped by cancelling the context.
func (h *Http) Start(ctx context.Context) error {
	errC := make(chan error)
	go func() {
		<-ctx.Done()
		err := h.srv.Shutdown(context.Background())
		errC <- err
	}()

	open := h.srv.Addr
	if open[0] == ':' {
		open = "http://localhost" + open
	} else {
		open = "http://" + open
	}

	slog.InfoContext(ctx, "starting the server", slog.String("address", h.srv.Addr), slog.String("open", open))
	err := h.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("starting up the server: %w", err)
	}

	// wait for shutdown to complete
	return <-errC
}

func (h *Http) StaticHandler() http.Handler {
	fs := http.FileServer(http.Dir("./static"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := path.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		}

		fs.ServeHTTP(w, r)
	})
}

type IndexModel struct {
	Title string
	Items []ItemModel
}

type ItemModel struct {
	ID      string
	Title   string
	Checked bool
}

func (h *Http) UIHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tasks, err := h.s.List(r.Context(), nil)
		if err != nil {
			slog.Error("listing the tasks", slog.String("err", err.Error()))
			httpErr(w, http.StatusInternalServerError)
			return
		}

		models := make([]ItemModel, 0, len(tasks))
		for _, t := range tasks {
			models = append(models, ItemModel{
				ID:      string(t.ID),
				Title:   t.Title,
				Checked: t.Done,
			})
		}

		err = h.ui.Render(w, IndexUI, IndexModel{
			Title: "Sam's tasks",
			Items: models,
		})

		if err != nil {
			slog.Error("rendering the index", slog.String("err", err.Error()))
			httpErr(w, http.StatusInternalServerError)
		}
	})

	return mux
}

func (h *Http) APIHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", h.HandlePostTodo)
	mux.HandleFunc("PUT /todos/{id}/toggle", h.HandlePostTodoToggle)
	return mux
}

func (h *Http) HandlePostTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := r.ParseForm()
	if err != nil {
		slog.ErrorContext(ctx, "parsing the form", slog.String("err", err.Error()))
		httpErr(w, http.StatusBadRequest)
		return
	}

	title := r.Form.Get("todo")
	if len(title) == 0 {
		slog.InfoContext(ctx, "empty todo field")
		httpErr(w, http.StatusBadRequest)
		return
	}

	task := todo.Task{
		ID:    todo.RandomID(),
		Title: title,
		Done:  false,
	}

	_, err = h.s.Upsert(ctx, task)
	if err != nil {
		slog.ErrorContext(ctx, "upserting the task", slog.String("err", err.Error()))
		httpErr(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Http) HandlePostTodoToggle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		slog.InfoContext(r.Context(), "empty id")
		httpErr(w, http.StatusBadRequest)
		return
	}

	_, err := h.h.Toggle(r.Context(), todo.ID(id))
	if errors.Is(err, todo.ErrTaskNotFound) {
		httpErr(w, http.StatusNotFound)
		return
	}

	if err != nil {
		slog.ErrorContext(r.Context(), "upserting the task", slog.String("err", err.Error()))
		httpErr(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func httpErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

type UI struct {
	tpl *template.Template
}

const (
	IndexUI = "index"
)

//go:embed ui/*
var uiFS embed.FS

func NewUI() (*UI, error) {
	tpl, err := template.ParseFS(uiFS, "ui/*.gohtml")
	if err != nil {
		return nil, fmt.Errorf("parsing the template from embedded filesystem, %w", err)
	}

	return &UI{tpl: tpl}, nil
}

func (u *UI) Render(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return u.tpl.ExecuteTemplate(w, name, data)
}

var (
	seed = time.Now().UnixNano()
	src  = rand.NewSource(seed)
)

type ridKey string

const RequestID = ridKey("request-id")

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rid := strconv.Itoa(int(src.Int63()))
		ctx := context.WithValue(r.Context(), RequestID, rid)
		defer func() {
			slog.InfoContext(ctx, "finished request", slog.String("id", rid), slog.String("took", time.Since(start).String()))
		}()
		slog.InfoContext(ctx, "received request", slog.String("id", rid), slog.String("method", r.Method), slog.String("url", r.URL.String()))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func closeBody(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		h.ServeHTTP(w, r)
	})
}

type storage struct {
	m map[string]todo.Task
}

func (s *storage) Upsert(ctx context.Context, t todo.Task) (stored todo.Task, err error) {
	slog.InfoContext(ctx, "upserting the task", slog.String("id", string(t.ID)))
	s.m[string(t.ID)] = t
	return t, nil
}

func (s *storage) List(ctx context.Context) ([]todo.Task, error) {
	tasks := make([]todo.Task, 0, len(s.m))
	for _, t := range s.m {
		tasks = append(tasks, t)
	}
	return tasks, nil
}
