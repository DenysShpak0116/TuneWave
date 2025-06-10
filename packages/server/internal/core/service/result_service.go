package service

import (
	"context"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/httpserver/handlers/dto"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type ResultService struct {
	GenericService[models.Result]
	CollectionSongRepository port.Repository[models.CollectionSong]
}

func NewResultService(repo port.Repository[models.Result], collectionSongRepository port.Repository[models.CollectionSong]) services.ResultService {
	return &ResultService{
		GenericService: GenericService[models.Result]{
			Repository: repo,
		},
		CollectionSongRepository: collectionSongRepository,
	}
}

func (rs *ResultService) ProcessUserResults(ctx context.Context, userID, collectionID uuid.UUID, request dto.SendResultRequest) ([]models.Result, error) {
	matrix, err := buildComparisonMatrix(request.Results)
	if err != nil {
		return nil, err
	}

	badCounts := calculateBadCounts(matrix)
	ranks := rankSongsByBadCounts(badCounts)

	if err := rs.saveUserRanks(ctx, userID, collectionID, ranks); err != nil {
		return nil, err
	}

	return rs.buildUserResultsDTO(ctx, userID, collectionID)
}

func buildComparisonMatrix(results []dto.ResultRequest) (map[uuid.UUID]map[uuid.UUID]int, error) {
	matrix := make(map[uuid.UUID]map[uuid.UUID]int)
	for _, result := range results {
		song1ID, err := uuid.Parse(result.Song1ID)
		if err != nil {
			return nil, err
		}
		if _, ok := matrix[song1ID]; !ok {
			matrix[song1ID] = make(map[uuid.UUID]int)
		}
		matrix[song1ID][song1ID] = 0

		for _, comparedTo := range result.ComparedTo {
			song2ID, err := uuid.Parse(comparedTo.Song2ID)
			if err != nil {
				return nil, err
			}
			if _, ok := matrix[song2ID]; !ok {
				matrix[song2ID] = make(map[uuid.UUID]int)
			}
			matrix[song1ID][song2ID] = comparedTo.Result
			matrix[song2ID][song1ID] = -comparedTo.Result
		}
	}
	return matrix, nil
}

func calculateBadCounts(matrix map[uuid.UUID]map[uuid.UUID]int) map[uuid.UUID]int {
	badCounts := make(map[uuid.UUID]int)
	for song1, comparisons := range matrix {
		for song2, result := range comparisons {
			if song1 != song2 && result <= -1 {
				badCounts[song1]++
			}
		}
	}

	for songID := range matrix {
		if _, ok := badCounts[songID]; !ok {
			badCounts[songID] = 0
		}
	}
	return badCounts
}

func rankSongsByBadCounts(badCounts map[uuid.UUID]int) map[uuid.UUID]int {
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

	ranks := make(map[uuid.UUID]int)
	for i, score := range scores {
		ranks[score.SongID] = i + 1
	}
	return ranks
}

func (rs *ResultService) saveUserRanks(ctx context.Context, userID, collectionID uuid.UUID, ranks map[uuid.UUID]int) error {
	for songID, rank := range ranks {
		collectionSongs, err := rs.CollectionSongRepository.NewQuery(ctx).Where(&models.CollectionSong{
			CollectionID: collectionID,
			SongID:       songID,
		}).Find()
		if err != nil {
			return err
		}
		if len(collectionSongs) == 0 {
			continue
		}
		cs := collectionSongs[0]

		result := &models.Result{
			UserID:           userID,
			CollectionSongID: cs.ID,
			SongRang:         rank,
		}
		if err := rs.Create(ctx, result); err != nil {
			return err
		}
	}
	return nil
}

func (rs *ResultService) buildUserResultsDTO(ctx context.Context, userID, collectionID uuid.UUID) ([]models.Result, error) {
	collectionSongs, err := rs.CollectionSongRepository.NewQuery(ctx).Where(&models.CollectionSong{
		CollectionID: collectionID,
	}).Find()
	if err != nil {
		return nil, err
	}

	var userResults []models.Result
	for _, cs := range collectionSongs {
		results, err := rs.Repository.NewQuery(ctx).Where(&models.Result{
			CollectionSongID: cs.ID,
			UserID:           userID,
		}).Preload("CollectionSong").Preload("CollectionSong.Song").Preload("User").Find()
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			continue
		}

		result := results[0]
		userResults = append(userResults, models.Result{
			BaseModel: models.BaseModel{
				ID:        result.ID,
				CreatedAt: result.CreatedAt,
			},
			SongRang: result.SongRang,
			UserID:   userID,
			User: models.User{
				BaseModel: models.BaseModel{
					ID:        userID,
					CreatedAt: result.User.CreatedAt,
				},
				Username:       result.User.Username,
				Email:          result.User.Email,
				PasswordHash:   result.User.PasswordHash,
				Role:           result.User.Role,
				ProfileInfo:    result.User.ProfileInfo,
				ProfilePicture: result.User.ProfilePicture,
			},
			CollectionSongID: collectionID,
			CollectionSong:   cs,
		})
	}
	return userResults, nil
}

func (rs *ResultService) GetUserResults(ctx context.Context, userID, collectionID uuid.UUID) ([]models.Result, error) {
	collectionSongs, err := rs.CollectionSongRepository.NewQuery(ctx).Where(&models.CollectionSong{
		CollectionID: collectionID,
	}).Find()
	if err != nil {
		return nil, err
	}

	userResults := make([]models.Result, 0)
	for _, cs := range collectionSongs {
		results, err := rs.Repository.NewQuery(ctx).Where(&models.Result{
			CollectionSongID: cs.ID,
			UserID:           userID,
		}).Preload("CollectionSong").Preload("CollectionSong.Song").Preload("User").Find()
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			continue
		}

		result := results[0]
		userResults = append(userResults, models.Result{
			BaseModel: models.BaseModel{
				ID:        result.ID,
				CreatedAt: result.CreatedAt,
			},
			SongRang: result.SongRang,
			UserID:   userID,
			User: models.User{
				BaseModel: models.BaseModel{
					ID:        userID,
					CreatedAt: result.User.CreatedAt,
				},
				Username:       result.User.Username,
				Email:          result.User.Email,
				PasswordHash:   result.User.PasswordHash,
				Role:           result.User.Role,
				ProfileInfo:    result.User.ProfileInfo,
				ProfilePicture: result.User.ProfilePicture,
			},
			CollectionSongID: collectionID,
			CollectionSong:   cs,
		})
	}

	return userResults, nil
}

func (rs *ResultService) GetCollectiveResults(ctx context.Context, collectionID uuid.UUID) (map[string]interface{}, error) {
	collectionSongs, err := rs.fetchCollectionSongsWithResults(ctx, collectionID)
	if err != nil || len(collectionSongs) == 0 {
		return map[string]interface{}{}, err
	}

	songIDToName, profiles := rs.buildUserProfiles(collectionSongs)
	groupedProfiles := groupSimilarUserProfiles(profiles)

	collectiveRank := calculateCollectiveRanking(collectionSongs, songIDToName)
	profileTable := buildProfileTable(groupedProfiles)

	return map[string]interface{}{
		"profileTable":   profileTable,
		"collectiveRank": collectiveRank,
	}, nil
}

func (rs *ResultService) fetchCollectionSongsWithResults(ctx context.Context, collectionID uuid.UUID) ([]models.CollectionSong, error) {
	return rs.CollectionSongRepository.NewQuery(ctx).Where(&models.CollectionSong{
		CollectionID: collectionID,
	}).Preload("Song").Preload("Results").Preload("Results.User").Find()
}

func (rs *ResultService) buildUserProfiles(collectionSongs []models.CollectionSong) (map[uuid.UUID]string, map[string]map[int][]string) {
	songIDToName := make(map[uuid.UUID]string)
	userProfiles := make(map[string]map[int][]string)

	for _, cs := range collectionSongs {
		songIDToName[cs.SongID] = cs.Song.Title

		for _, result := range cs.Results {
			username := result.User.Username
			if _, ok := userProfiles[username]; !ok {
				userProfiles[username] = make(map[int][]string)
			}
			userProfiles[username][result.SongRang] = append(userProfiles[username][result.SongRang], cs.Song.Title)
		}
	}
	return songIDToName, userProfiles
}

type groupedProfile struct {
	Users   []string
	Profile map[int][]string
}

func groupSimilarUserProfiles(profiles map[string]map[int][]string) []groupedProfile {
	seen := make(map[string]bool)
	var grouped []groupedProfile

	for user, profile := range profiles {
		if seen[user] {
			continue
		}
		group := groupedProfile{Users: []string{user}, Profile: profile}
		for otherUser, otherProfile := range profiles {
			if user != otherUser && !seen[otherUser] && equalProfiles(profile, otherProfile) {
				group.Users = append(group.Users, otherUser)
				seen[otherUser] = true
			}
		}
		seen[user] = true
		grouped = append(grouped, group)
	}
	return grouped
}

func calculateCollectiveRanking(collectionSongs []models.CollectionSong, songIDToName map[uuid.UUID]string) map[int]map[string]interface{} {
	comparisons := make(map[uuid.UUID]map[uuid.UUID]int)
	songIDs := make(map[uuid.UUID]bool)

	for _, cs := range collectionSongs {
		songIDs[cs.SongID] = true
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

	remaining := make(map[uuid.UUID]bool)
	for id := range songIDs {
		remaining[id] = true
	}

	rank := []uuid.UUID{}
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
		minScore := math.MaxInt
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
		finalRank[len(rank)-i] = map[string]interface{}{
			"songId":   id.String(),
			"songName": songIDToName[id],
		}
	}
	return finalRank
}

func buildProfileTable(groups []groupedProfile) map[int]map[string][]string {
	table := make(map[int]map[string][]string)
	for _, group := range groups {
		usernameKey := strconv.Itoa(len(group.Users)) + ": " + joinUsernames(group.Users)
		for rank, songs := range group.Profile {
			if table[rank] == nil {
				table[rank] = make(map[string][]string)
			}
			table[rank][usernameKey] = songs
		}
	}
	return table
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
