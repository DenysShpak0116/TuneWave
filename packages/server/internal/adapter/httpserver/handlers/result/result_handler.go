package result

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ResultHandler struct {
	ResultService         services.ResultService
	CollectionSongService services.CollectionSongService
}

func NewResultHandler(resultService services.ResultService, collectionSongService services.CollectionSongService) *ResultHandler {
	return &ResultHandler{
		ResultService:         resultService,
		CollectionSongService: collectionSongService,
	}
}

type UserResultsDTO struct {
	CollectionSongID uuid.UUID `json:"collectionSongId"`

	SongID   uuid.UUID `json:"songId"`
	SongName string    `json:"songName"`

	UserID   uuid.UUID `json:"userId"`
	UserName string    `json:"userName"`
	SongRang int       `json:"songRang"`
}

type SendResultRequest struct {
	Results []struct {
		Song1ID    string `json:"song1Id"`
		ComparedTo []struct {
			Song2ID string `json:"song2Id"`
			Result  int    `json:"result"`
		} `json:"comparedTo"`
	} `json:"results"`
}

// @Summary Send result
// @Description Send result
// @Tags result
// @Accept json
// @Produce json
// @Param request body SendResultRequest true "Send result"
// @Router /collections/{id}/send-results [post]
func (h *ResultHandler) SendResult(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	var request SendResultRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	matrix := make(map[uuid.UUID]map[uuid.UUID]int, 0)
	for _, result := range request.Results {
		song1ID, err := uuid.Parse(result.Song1ID)
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
			return
		}

		if _, ok := matrix[song1ID]; !ok {
			matrix[song1ID] = make(map[uuid.UUID]int, 0)
		}
		matrix[song1ID][song1ID] = 0
		for _, comparedTo := range result.ComparedTo {
			song2ID, err := uuid.Parse(comparedTo.Song2ID)
			if err != nil {
				handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid song ID", err)
				return
			}

			matrix[song1ID][song2ID] = comparedTo.Result
			matrix[song2ID][song1ID] = -comparedTo.Result
		}
	}

	badCounts := make(map[uuid.UUID]int)
	for song1, comparisons := range matrix {
		for song2, result := range comparisons {
			if song1 == song2 {
				continue
			}
			if result <= -1 {
				badCounts[song1]++
			}
		}
	}

	type songScore struct {
		SongID uuid.UUID
		Bad    int
	}
	var scores []songScore
	for songID, bad := range badCounts {
		scores = append(scores, songScore{SongID: songID, Bad: bad})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Bad < scores[j].Bad
	})

	songRanks := make(map[uuid.UUID]int)
	for i, score := range scores {
		songRanks[score.SongID] = i + 1
	}

	for songID, rank := range songRanks {
		collectionSongs, err := h.CollectionSongService.Where(context.Background(), &models.CollectionSong{
			CollectionID: collectionUUID,
			SongID:       songID,
		})
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "CollectionSong not found", err)
			return
		}
		if len(collectionSongs) == 0 {
			handlers.RespondWithError(w, r, http.StatusNotFound, "CollectionSong not found", nil)
			return
		}
		cs := collectionSongs[0]

		result := &models.Result{
			UserID:           userUUID,
			CollectionSongID: cs.ID,
			SongRang:         rank,
		}
		if err := h.ResultService.Create(context.Background(), result); err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to save result", err)
			return
		}
	}

	collectionSongs, err := h.CollectionSongService.Where(context.Background(), &models.CollectionSong{
		CollectionID: collectionUUID,
	})
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "CollectionSong not found", err)
		return
	}

	userResults := make([]UserResultsDTO, 0)
	for _, cs := range collectionSongs {
		results, err := h.ResultService.Where(context.Background(), &models.Result{
			CollectionSongID: cs.ID,
			UserID:           userUUID,
		}, "CollectionSong", "CollectionSong.Song", "User")
		if err != nil {
			handlers.RespondWithError(w, r, http.StatusInternalServerError, "Failed to get results", err)
			return
		}
		if len(results) == 0 {
			continue
		}

		result := results[0]
		userResults = append(userResults, UserResultsDTO{
			CollectionSongID: result.CollectionSongID,
			SongID:           result.CollectionSong.SongID,
			SongName:         result.CollectionSong.Song.Title,
			UserID:           userUUID,
			UserName:         result.User.Username,
			SongRang:         result.SongRang,
		})
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userResults)
}

// protected.Get("/{id}/get-user-results/", resultHandler.GetUserResults)
// protected.Get("/{id}/get-results/", resultHandler.GetCollectiveResults)

func (h *ResultHandler) GetUserResults(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (h *ResultHandler) GetCollectiveResults(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
