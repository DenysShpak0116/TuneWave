package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port/services"
	"github.com/google/uuid"
)

type CollectionService struct {
	*GenericService[models.Collection]
	CollectionSongRepository port.Repository[models.CollectionSong]
	FileStorage              port.FileStorage
}

func NewCollectionService(
	repo port.Repository[models.Collection],
	fileStorage port.FileStorage,
	collectionSongRepository port.Repository[models.CollectionSong],
	logger *slog.Logger,
) services.CollectionService {
	return &CollectionService{
		GenericService:           NewGenericService(repo, logger),
		FileStorage:              fileStorage,
		CollectionSongRepository: collectionSongRepository,
	}
}

func (cs *CollectionService) SaveCollection(ctx context.Context, collectionParams services.SaveCollectionParams) (*models.Collection, error) {
	const op = "core.service.CollectionService.SaveCollection"
	logger := cs.logger.With(
		slog.String("op", op),
	)

	coverURL, err := cs.saveCoverFile(ctx, SaveFileParams{
		UserID:   collectionParams.UserID,
		Filename: collectionParams.CoverHeader.Filename,
		File:     collectionParams.Cover,
	})
	if err != nil {
		logger.Error("Error while saving collection", "err", err.Error())
		return nil, err
	}

	collection := &models.Collection{
		Title:       collectionParams.Title,
		Description: collectionParams.Description,
		CoverURL:    coverURL,
		UserID:      collectionParams.UserID,
	}
	if err := cs.repository.Add(ctx, collection); err != nil {
		logger.Error("Error while saving collection", "err", err.Error())
		return nil, err
	}

	logger.Info("Collecting saved successfully")
	return collection, nil
}

func (cs *CollectionService) UpdateCollection(ctx context.Context, id uuid.UUID, collectionParams services.UpdateCollectionParams) (*models.Collection, error) {
	const op = "core.service.CollectionService.SaveCollection"
	logger := cs.logger.With(
		slog.String("op", op),
	)

	collections, err := cs.repository.NewQuery(ctx).
		Where("id = ?", id).
		Find()
	if err != nil {
		logger.Error("Error while updating collection", "err", err)
		return nil, err
	}
	if len(collections) == 0 {
		logger.Error("Error to find collection")
		return nil, fmt.Errorf("collection with id %s not found", id)
	}
	collection := &collections[0]

	if err := cs.FileStorage.Remove(ctx, helpers.ExtractS3Key(collection.CoverURL)); err != nil {
		logger.Error("Error removing old collection's cover", "err", err.Error())
		return nil, err
	}

	if collectionParams.CoverHeader != nil && collectionParams.Cover != nil {
		coverURL, err := cs.saveCoverFile(ctx, SaveFileParams{
			UserID:   collectionParams.UserID,
			Filename: collectionParams.CoverHeader.Filename,
			File:     collectionParams.Cover,
		})
		if err != nil {
			logger.Error("Error adding new cover to collection", "err", err.Error())
			return nil, err
		}
		collection.CoverURL = coverURL
	}

	collection.Title = collectionParams.Title
	collection.Description = collectionParams.Description

	if err := cs.repository.Update(ctx, collection); err != nil {
		logger.Error("Error updaing collection", "err", err.Error())
		return nil, err
	}

	logger.Info("Collectiong successfully updated", "id", id.String())
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

func (cs *CollectionService) GetMany(ctx context.Context, limit, page int, sort, order string, preloads ...string) ([]models.Collection, error) {
	query := cs.repository.NewQuery(ctx).
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
	const op = "core.service.CollectionService.SaveCollection"
	logger := cs.logger.With(
		slog.String("op", op),
	)

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
		logger.Error("Error while trying to find collection songs", "collectionId", collectionID.String(), "err", err.Error())
		return nil, err
	}

	songs := []models.Song{}
	for _, collectionSong := range collectionSongs {
		songs = append(songs, collectionSong.Song)
	}

	logger.Info("Collection songs found", "collectionId", collectionID.String())
	return songs, nil
}
