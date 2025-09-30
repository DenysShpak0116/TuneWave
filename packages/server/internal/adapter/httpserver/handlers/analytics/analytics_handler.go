package analytics

import (
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/go-chi/render"
)

type AnalyticsHandler struct {
	eventService services.EventService
}

func NewAnalyticsHandler(eventService services.EventService) *AnalyticsHandler {
	return &AnalyticsHandler{eventService: eventService}
}

// ListensByDay godoc
// @Summary Get number of listens grouped by day
// @Tags analytics
// @Produce json
// @Router /analytics/listens-by-day [get]
func (ah *AnalyticsHandler) ListensByDay(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	stats := make(map[string]int64)
	for _, ev := range events {
		day := ev.CreatedAt.Format("2006-01-02")
		stats[day]++
	}

	render.JSON(w, r, stats)
	return nil
}

// PeakAndSilentDays godoc
// @Summary Get peak and silent day of listens
// @Tags analytics
// @Produce json
// @Router /analytics/peak-silent [get]
func (ah *AnalyticsHandler) PeakAndSilentDays(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	stats := make(map[string]int64)
	for _, ev := range events {
		day := ev.CreatedAt.Format("2006-01-02")
		stats[day]++
	}

	var peakDay, silentDay string
	var maxCount, minCount int64
	minCount = 1<<31 - 1
	for day, count := range stats {
		if count > maxCount {
			maxCount = count
			peakDay = day
		}
		if count < minCount {
			minCount = count
			silentDay = day
		}
	}

	render.JSON(w, r, map[string]any{
		"peak_day":     peakDay,
		"peak_count":   maxCount,
		"silent_day":   silentDay,
		"silent_count": minCount,
	})
	return nil
}

// AvgListensPerUser godoc
// @Summary Get average listens per user
// @Tags analytics
// @Produce json
// @Router /analytics/avg-listens [get]
func (ah *AnalyticsHandler) AvgListensPerUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	userCounts := make(map[string]int64)
	for _, ev := range events {
		userCounts[ev.UserID.String()]++
	}

	var total, users int64
	for _, count := range userCounts {
		total += count
		users++
	}
	avg := float64(total) / float64(users)

	render.JSON(w, r, map[string]any{"avg": avg})
	return nil
}

// MedianListensPerUser godoc
// @Summary Get median listens per user
// @Tags analytics
// @Produce json
// @Router /analytics/median-listens [get]
func (ah *AnalyticsHandler) MedianListensPerUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	userCounts := make(map[string]int64)
	for _, ev := range events {
		userCounts[ev.UserID.String()]++
	}

	values := make([]int64, 0, len(userCounts))
	for _, count := range userCounts {
		values = append(values, count)
	}
	slices.Sort(values)

	var median float64
	n := len(values)
	if n%2 == 0 {
		median = float64(values[n/2-1]+values[n/2]) / 2
	} else {
		median = float64(values[n/2])
	}

	render.JSON(w, r, map[string]any{"median": median})
	return nil
}

// MostPopularTrack godoc
// @Summary Get most popular track
// @Tags analytics
// @Produce json
// @Router /analytics/most-popular-track [get]
func (ah *AnalyticsHandler) MostPopularTrack(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	trackCounts := make(map[string]int64)
	for _, ev := range events {
		trackCounts[ev.TrackID.String()]++
	}

	var popularTrack string
	var maxCount int64
	for trackID, count := range trackCounts {
		if count > maxCount {
			maxCount = count
			popularTrack = trackID
		}
	}

	render.JSON(w, r, map[string]any{
		"track_id": popularTrack,
		"count":    maxCount,
	})
	return nil
}

// TracksWithPopular godoc
// @Summary Get tracks most often listened together with the most popular track
// @Tags analytics
// @Produce json
// @Router /analytics/tracks-with-popular [get]
func (ah *AnalyticsHandler) TracksWithPopular(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	trackCounts := make(map[string]int64)
	for _, ev := range events {
		trackCounts[ev.TrackID.String()]++
	}
	var popularTrack string
	var maxCount int64
	for trackID, count := range trackCounts {
		if count > maxCount {
			maxCount = count
			popularTrack = trackID
		}
	}

	coCounts := make(map[string]int64)
	userTracks := make(map[string]map[string]struct{})
	for _, ev := range events {
		if userTracks[ev.UserID.String()] == nil {
			userTracks[ev.UserID.String()] = make(map[string]struct{})
		}
		userTracks[ev.UserID.String()][ev.TrackID.String()] = struct{}{}
	}
	for _, tracks := range userTracks {
		if _, ok := tracks[popularTrack]; ok {
			for t := range tracks {
				if t != popularTrack {
					coCounts[t]++
				}
			}
		}
	}

	render.JSON(w, r, coCounts)
	return nil
}

// CommonTrackCombos godoc
// @Summary Get most common track combinations
// @Tags analytics
// @Produce json
// @Router /analytics/common-combos [get]
func (ah *AnalyticsHandler) CommonTrackCombos(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	userTracks := make(map[string][]string)
	for _, ev := range events {
		userTracks[ev.UserID.String()] = append(userTracks[ev.UserID.String()], ev.TrackID.String())
	}

	comboCounts := make(map[string]int64)
	for _, tracks := range userTracks {
		sort.Strings(tracks)
		key := strings.Join(tracks, ",")
		comboCounts[key]++
	}

	type kv struct {
		Key   string
		Value int64
	}
	var sorted []kv
	for k, v := range comboCounts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value > sorted[j].Value })

	render.JSON(w, r, sorted)
	return nil
}

// RareTrackCombos godoc
// @Summary Get rarest track combinations
// @Tags analytics
// @Produce json
// @Router /analytics/rare-combos [get]
func (ah *AnalyticsHandler) RareTrackCombos(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	events, err := ah.eventService.Where(ctx, &models.Event{EventType: models.PlayEvent})
	if err != nil {
		return helpers.InternalServerError("failed to fetch listen events")
	}

	userTracks := make(map[string][]string)
	for _, ev := range events {
		userTracks[ev.UserID.String()] = append(userTracks[ev.UserID.String()], ev.TrackID.String())
	}

	comboCounts := make(map[string]int64)
	for _, tracks := range userTracks {
		sort.Strings(tracks)
		key := strings.Join(tracks, ",")
		comboCounts[key]++
	}

	type kv struct {
		Key   string
		Value int64
	}
	var sorted []kv
	for k, v := range comboCounts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value < sorted[j].Value })

	render.JSON(w, r, sorted)
	return nil
}
