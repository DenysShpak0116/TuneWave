package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type CollectionService struct {
	*GenericService[models.Collection]
	CollectionSongRepository port.Repository[models.CollectionSong]
	ReactionRepository       port.Repository[models.UserReaction]
	FileStorage              port.FileStorage
}

func NewCollectionService(
	repo port.Repository[models.Collection],
	fileStorage port.FileStorage,
	collectionSongRepository port.Repository[models.CollectionSong],
	reactionRepository port.Repository[models.UserReaction],
) services.CollectionService {
	return &CollectionService{
		GenericService:           NewGenericService(repo),
		FileStorage:              fileStorage,
		CollectionSongRepository: collectionSongRepository,
		ReactionRepository:       reactionRepository,
	}
}

func (cs *CollectionService) SaveCollection(ctx context.Context, collectionParams services.SaveCollectionParams) (*models.Collection, error) {
	coverURL, err := cs.saveCoverFile(ctx, SaveFileParams{
		UserID:   collectionParams.UserID,
		Filename: collectionParams.CoverHeader.Filename,
		File:     collectionParams.Cover,
	})
	if err != nil {
		return nil, err
	}

	collection := &models.Collection{
		Title:       collectionParams.Title,
		Description: collectionParams.Description,
		CoverURL:    coverURL,
		UserID:      collectionParams.UserID,
	}

	if err := cs.Repository.Add(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

func (cs *CollectionService) UpdateCollection(ctx context.Context, id uuid.UUID, collectionParams services.UpdateCollectionParams) (*models.Collection, error) {
	collections, err := cs.Repository.NewQuery(ctx).
		Where("id = ?", id).
		Find()
	if err != nil {
		return nil, err
	}
	if len(collections) == 0 {
		return nil, fmt.Errorf("collection with id %s not found", id)
	}
	collection := &collections[0]

	if err := cs.FileStorage.Remove(ctx, extractS3Key(collection.CoverURL)); err != nil {
		return nil, err
	}

	if collectionParams.CoverHeader != nil && collectionParams.Cover != nil {
		coverURL, err := cs.saveCoverFile(ctx, SaveFileParams{
			UserID:   collectionParams.UserID,
			Filename: collectionParams.CoverHeader.Filename,
			File:     collectionParams.Cover,
		})
		if err != nil {
			return nil, err
		}
		collection.CoverURL = coverURL
	}

	collection.Title = collectionParams.Title
	collection.Description = collectionParams.Description

	if err := cs.Repository.Update(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

type SaveFileParams struct {
	UserID   uuid.UUID
	Filename string
	File     multipart.File
}

func (cs *CollectionService) saveCoverFile(ctx context.Context, SaveFileParams SaveFileParams) (string, error) {
	key := fmt.Sprintf("collectionCovers/%s/%d-%s", SaveFileParams.UserID, time.Now().Unix(), SaveFileParams.Filename)

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, SaveFileParams.File); err != nil {
		return "", err
	}

	url, err := cs.FileStorage.Save(ctx, key, buf)
	if err != nil {
		return "", err
	}

	return url, nil
}

func extractS3Key(fullURL string) string {
	const baseURL = "https://tunewavebucket.s3.eu-west-3.amazonaws.com/"
	return strings.TrimPrefix(fullURL, baseURL)
}

func (cs *CollectionService) GetMany(ctx context.Context, limit, page int, sort, order string, preloads ...string) ([]models.Collection, error) {
	query := cs.Repository.NewQuery(ctx).
		Take(limit).
		Skip(page).
		Order(fmt.Sprintf("%s %s", sort, order))

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	collections, err := query.Find()
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (cs *CollectionService) GetCollectionSongs(
	ctx context.Context,
	collectionID uuid.UUID,
	search, sortBy, order string,
	page, limit int,
) ([]models.Song, error) {
	offset := (page - 1) * limit

	query := cs.CollectionSongRepository.NewQuery(ctx).
		Where("collection_id = ?", collectionID).
		Preload("Song").
		Preload("Song.Authors.Author").
		Preload("Song.SongTags.Tag")

	query = query.Join("JOIN songs ON songs.id = collection_songs.song_id")

	if search != "" {
		query = query.Where("songs.title ILIKE ?", "%"+search+"%")
	}

	if sortBy != "" {
		if order != "asc" && order != "desc" {
			order = "asc"
		}

		var orderClause string
		switch sortBy {
		case "title":
			orderClause = fmt.Sprintf("songs.%s %s", sortBy, order)
		case "added_at":
			orderClause = fmt.Sprintf("collection_songs.created_at %s", order)
		default:
			orderClause = "collection_songs.created_at desc"
		}

		query = query.Order(orderClause)
	} else {
		query = query.Order("collection_songs.created_at desc")
	}

	query = query.Skip(offset).Take(limit)

	collectionSongs, err := query.Find()
	if err != nil {
		return nil, err
	}

	songs := []models.Song{}
	for _, collectionSong := range collectionSongs {
		songs = append(songs, collectionSong.Song)
	}

	return songs, nil
}
