package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/dtos"
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

func (cs *CollectionService) GetFullDTOByID(ctx context.Context, id uuid.UUID) (*dtos.CollectionExtendedDTO, error) {
	collections, err := cs.Repository.NewQuery(ctx).
		Where("id = ?", id).
		Preload("User").
		Preload("CollectionSongs").
		Preload("CollectionSongs.Song").
		Preload("CollectionSongs.Song.User").
		Preload("CollectionSongs.Song.Authors").
		Preload("CollectionSongs.Song.Authors.Author").
		Preload("CollectionSongs.Song.SongTags").
		Preload("CollectionSongs.Song.SongTags.Tag").
		Preload("CollectionSongs.Song.Comments").
		Preload("CollectionSongs.Song.Comments.User").
		Find()
	if err != nil {
		return nil, err
	}
	if len(collections) == 0 {
		return nil, fmt.Errorf("collection with id %s not found", id)
	}

	collection := collections[0]

	collectionSongs := make([]dtos.SongExtendedDTO, len(collection.CollectionSongs))
	for i, collectionSong := range collection.CollectionSongs {
		songLikes, err := cs.ReactionRepository.NewQuery(ctx).
			Where("song_id = ? AND type = ?", collectionSong.Song.ID, "like").
			Count()
		if err != nil {
			return nil, err
		}
		songDislikes, err := cs.ReactionRepository.NewQuery(ctx).
			Where("song_id = ? AND type = ?", collectionSong.Song.ID, "dislike").
			Count()
		if err != nil {
			return nil, err
		}

		songAuthorsDTO := make([]dtos.AuthorDTO, len(collectionSong.Song.Authors))
		for j, author := range collectionSong.Song.Authors {
			songAuthorsDTO[j] = dtos.AuthorDTO{
				ID:   author.AuthorID,
				Name: author.Author.Name,
				Role: author.Role,
			}
		}

		songTagsDTO := make([]dtos.TagDTO, len(collectionSong.Song.SongTags))
		for j, tag := range collectionSong.Song.SongTags {
			songTagsDTO[j] = dtos.TagDTO{
				Name: tag.Tag.Name,
			}
		}

		songCommentsDTO := make([]dtos.CommentDTO, len(collectionSong.Song.Comments))
		for j, comment := range collectionSong.Song.Comments {
			songCommentsDTO[j] = dtos.CommentDTO{
				ID: comment.ID,
				Author: dtos.UserDTO{
					ID:             comment.User.ID,
					Username:       comment.User.Username,
					Role:           comment.User.Role,
					ProfilePicture: comment.User.ProfilePicture,
					ProfileInfo:    comment.User.ProfileInfo,
				},
				Content:   "",
				CreatedAt: time.Time{},
			}
		}
		collectionSongs[i] = dtos.SongExtendedDTO{
			ID:         collectionSong.Song.ID,
			Title:      collectionSong.Song.Title,
			Duration:   formatDuration(time.Duration(collectionSong.Song.Duration)),
			CoverURL:   collectionSong.Song.CoverURL,
			Listenings: collectionSong.Song.Listenings,
			User: dtos.UserDTO{
				ID:             collectionSong.Song.User.ID,
				Username:       collectionSong.Song.User.Username,
				ProfilePicture: collectionSong.Song.User.ProfilePicture,
			},
			CreatedAt: collectionSong.Song.CreatedAt,
			Genre:     collectionSong.Song.Genre,
			SongURL:   collectionSong.Song.SongURL,
			Likes:     songLikes,
			Dislikes:  songDislikes,
			Authors:   songAuthorsDTO,
			SongTags:  songTagsDTO,
			Comments:  songCommentsDTO,
		}
	}

	collectionDTO := &dtos.CollectionExtendedDTO{
		ID:          collection.ID,
		Title:       collection.Title,
		Description: collection.Description,
		CoverURL:    collection.CoverURL,
		CreatedAt:   collection.CreatedAt,
		User: dtos.UserDTO{
			ID:             collection.User.ID,
			Username:       collection.User.Username,
			ProfilePicture: collection.User.ProfilePicture,
		},
		CollectionSongs: collectionSongs,
	}

	return collectionDTO, nil
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
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

	collection := collections[0]

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

	newCollection, err := cs.Repository.Update(ctx, &collection)
	if err != nil {
		return nil, err
	}
	return newCollection, nil
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

func (cs *CollectionService) GetMany(ctx context.Context, limit, page int, sort, order string) ([]models.Collection, error) {
	collections, err := cs.Repository.NewQuery(ctx).
		Take(limit).
		Skip(page).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Preload("User").
		Find()
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
) ([]dtos.SongExtendedDTO, error) {
	offset := (page - 1) * limit

	query := cs.CollectionSongRepository.NewQuery(ctx).
		Where("collection_id = ?", collectionID).
		Preload("Song").
		Preload("Song.User").
		Preload("Song.Authors.Author").
		Preload("Song.SongTags.Tag").
		Preload("Song.Comments").
		Preload("Song.Comments.User").
		Preload("Song.Reactions")

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

	result :=  make([]dtos.SongExtendedDTO, 0)
	for _, csong := range collectionSongs {
		song := csong.Song

		likes := int64(0)
		dislikes := int64(0)
		for _, reaction := range song.Reactions {
			if reaction.Type == "like" {
				likes++
			} else if reaction.Type == "dislike" {
				dislikes++
			}
		}

		var authors []dtos.AuthorDTO
		for _, sa := range song.Authors {
			authors = append(authors, dtos.AuthorDTO{
				ID:   sa.Author.ID,
				Name: sa.Author.Name,
				Role: sa.Role,
			})
		}

		var tags []dtos.TagDTO
		for _, st := range song.SongTags {
			tags = append(tags, dtos.TagDTO{
				Name: st.Tag.Name,
			})
		}

		var comments []dtos.CommentDTO
		for _, c := range song.Comments {
			comments = append(comments, dtos.CommentDTO{
				ID:        c.ID,
				Content:   c.Content,
				CreatedAt: c.CreatedAt,
				Author: dtos.UserDTO{
					ID:             c.User.ID,
					Username:       c.User.Username,
					Role:           c.User.Role,
					ProfilePicture: c.User.ProfilePicture,
					ProfileInfo:    c.User.ProfileInfo,
				},
			})
		}

		dto := dtos.SongExtendedDTO{
			ID:         song.ID,
			CreatedAt:  song.CreatedAt,
			Duration:   song.Duration.String(),
			Title:      song.Title,
			Genre:      song.Genre,
			SongURL:    song.SongURL,
			CoverURL:   song.CoverURL,
			Listenings: song.Listenings,
			Likes:      likes,
			Dislikes:   dislikes,
			User: dtos.UserDTO{
				ID:             song.User.ID,
				Username:       song.User.Username,
				Role:           song.User.Role,
				ProfilePicture: song.User.ProfilePicture,
				ProfileInfo:    song.User.ProfileInfo,
			},
			Authors:  authors,
			SongTags: tags,
			Comments: comments,
		}

		result = append(result, dto)
	}

	return result, nil
}
