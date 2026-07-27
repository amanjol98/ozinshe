package episodes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type EpisodeHandler struct {
	service *EpisodeService
}

func NewEpisodeHandler(service *EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{service: service}
}

type episodeRequest struct {
	EpisodeNumber int    `json:"episode_number"`
	VideoID       string `json:"video_id"`
}

func (h *EpisodeHandler) AddEpisodeToSeason(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req episodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	seasonIDString := r.PathValue("id")
	seasonID, err := strconv.Atoi(seasonIDString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.AddEpisodeToSeason(ctx, seasonID, req.EpisodeNumber, req.VideoID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EpisodeHandler) DeleteEpisode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	episodeID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteEpisode(ctx, episodeID)
	if err != nil {
		if errors.Is(err, ErrEpisodeNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
