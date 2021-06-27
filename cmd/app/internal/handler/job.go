package handler

import (
	"encoding/json"
	core "github.com/calebikhuohon/hasty-test"
	external2 "github.com/calebikhuohon/hasty-test/cmd/app/external"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req external2.CreateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("could not decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := core.Job{
		ObjectID: req.ObjectId,
	}

	jobId, err := h.jobService.RunJob(ctx, job)
	if err != nil {
		log.Printf("failed to run job with error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(jobId)
}

func (h Handler) GetJobDetails(w http.ResponseWriter, r *http.Request)  {
	jobId := chi.URLParam(r, "job_id")

	unescapedjobId, err := url.PathUnescape(jobId)
	if err != nil {
		log.Println(err)
	}

	data, err := h.jobService.GetJobDetails(r.Context(), unescapedjobId)
	if err != nil {
		log.Printf("couldn't fetch job details for job with id %s: %v", jobId, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}


