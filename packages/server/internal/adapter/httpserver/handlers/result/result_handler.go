package result

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
// @Security BearerAuth
// @Tags result
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
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

			if _, ok := matrix[song2ID]; !ok {
				matrix[song2ID] = make(map[uuid.UUID]int)
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

	for songID := range matrix {
		if _, ok := badCounts[songID]; !ok {
			badCounts[songID] = 0
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

// @Summary Get user results
// @Description Get user results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-user-results [get]
func (h *ResultHandler) GetUserResults(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(string)
	userUUID, err := uuid.Parse(userId)
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

// @Summary Get collective results
// @Description Get collective results
// @Security BearerAuth
// @Tags result
// @Produce json
// @Param id path string true "Collection ID"
// @Router /collections/{id}/get-results [get]
func (h *ResultHandler) GetCollectiveResults(w http.ResponseWriter, r *http.Request) {
	collectionID := chi.URLParam(r, "id")
	collectionUUID, err := uuid.Parse(collectionID)
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusBadRequest, "Invalid collection ID", err)
		return
	}

	collectionSongs, err := h.CollectionSongService.Where(context.Background(), &models.CollectionSong{
		CollectionID: collectionUUID,
	}, "Song", "Results", "Results.User")
	if err != nil {
		handlers.RespondWithError(w, r, http.StatusInternalServerError, "CollectionSong not found", err)
		return
	}
	if len(collectionSongs) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{})
		return
	}
	type UserProfile map[int][]string
	type RankedUserProfile map[string]UserProfile
	profiles := make(RankedUserProfile)

	songIDToName := make(map[uuid.UUID]string)
	songIDSet := make(map[uuid.UUID]struct{})
	usernames := make(map[string]struct{})

	for _, cs := range collectionSongs {
		songIDToName[cs.SongID] = cs.Song.Title
		songIDSet[cs.SongID] = struct{}{}
		for _, res := range cs.Results {
			username := res.User.Username
			usernames[username] = struct{}{}
			if _, ok := profiles[username]; !ok {
				profiles[username] = make(UserProfile)
			}
			profiles[username][res.SongRang] = append(profiles[username][res.SongRang], cs.Song.Title)
		}
	}
	type groupedProfile struct {
		Users   []string
		Profile UserProfile
	}
	grouped := []groupedProfile{}
	seen := map[string]bool{}
	for u, p := range profiles {
		if seen[u] {
			continue
		}
		group := groupedProfile{Users: []string{u}, Profile: p}
		for u2, p2 := range profiles {
			if u != u2 && !seen[u2] && equalProfiles(p, p2) {
				group.Users = append(group.Users, u2)
				seen[u2] = true
			}
		}
		seen[u] = true
		grouped = append(grouped, group)
	}

	songIDs := make([]uuid.UUID, 0, len(songIDSet))
	for id := range songIDSet {
		songIDs = append(songIDs, id)
	}

	comparisons := make(map[uuid.UUID]map[uuid.UUID]int)
	for _, cs := range collectionSongs {
		for _, res := range cs.Results {
			for _, cs2 := range collectionSongs {
				if cs.SongID == cs2.SongID {
					continue
				}
				var r1, r2 int
				for _, r := range cs.Results {
					if r.UserID == res.UserID {
						r1 = r.SongRang
					}
				}
				for _, r := range cs2.Results {
					if r.UserID == res.UserID {
						r2 = r.SongRang
					}
				}
				if r1 < r2 {
					if comparisons[cs.SongID] == nil {
						comparisons[cs.SongID] = map[uuid.UUID]int{}
					}
					comparisons[cs.SongID][cs2.SongID]++
				}
			}
		}
	}

	rank := []uuid.UUID{}
	remaining := make(map[uuid.UUID]bool)
	for _, id := range songIDs {
		remaining[id] = true
	}

	for len(remaining) > 0 {
		scores := make(map[uuid.UUID]int)
		for a := range remaining {
			for b := range remaining {
				if a == b {
					continue
				}
				if comparisons[a][b] > comparisons[b][a] {
					scores[b]++
				} else {
					scores[a]++
				}
			}
		}
		var minSong uuid.UUID
		minScore := 1 << 30
		for s := range remaining {
			if scores[s] < minScore {
				minScore = scores[s]
				minSong = s
			}
		}
		delete(remaining, minSong)
		rank = append([]uuid.UUID{minSong}, rank...)
	}

	finalRank := make(map[int]map[string]interface{})
	for i, id := range rank {
		finalRank[i+1] = map[string]interface{}{
			"songId":   id.String(),
			"songName": songIDToName[id],
		}
	}

	profileTable := make(map[int]map[string][]string)
	for _, group := range grouped {
		usernameKey := strconv.Itoa(len(group.Users)) + ": " + joinUsernames(group.Users)
		for rank, songs := range group.Profile {
			if profileTable[rank] == nil {
				profileTable[rank] = map[string][]string{}
			}
			profileTable[rank][usernameKey] = songs
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"profileTable":   profileTable,
		"collectiveRank": finalRank,
	})
}

func equalProfiles(a, b map[int][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		v2, ok := b[k]
		if !ok || !equalStringSlices(v, v2) {
			return false
		}
	}
	return true
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	amap := make(map[string]int)
	for _, x := range a {
		amap[x]++
	}
	for _, x := range b {
		amap[x]--
		if amap[x] < 0 {
			return false
		}
	}
	return true
}

func joinUsernames(users []string) string {
	sort.Strings(users)
	return strings.Join(users, ", ")
}
