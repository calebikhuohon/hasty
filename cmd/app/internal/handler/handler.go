package handler

import (
	"context"
	core "github.com/calebikhuohon/hasty-test"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
)

type JobService interface {
	RunJob(ctx context.Context, job core.Job) (uuid.UUID, error)
	GetJobDetails(ctx context.Context, job string) (core.Job, error)
}

type Handler struct {
	router     http.Handler
	jobService JobService
}

type Config struct {
	Timeout time.Duration
}

func New(
	service JobService,
	config Config,
) http.Handler {
	r := chi.NewRouter()

	h := &Handler{
		router:     r,
		jobService: service,
	}

	timeout := 10 * time.Second
	if config.Timeout > 0 {
		timeout = config.Timeout
	}

	r.Use(
		chimiddleware.Timeout(timeout),
		chimiddleware.SetHeader("Content-Type", "application/json"),
	)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateJob)
		r.Get("/{job_id}", h.GetJobDetails)
	})

	return h
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
