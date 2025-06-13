package dto

import (
	"time"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/domain/models"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/helpers"
	"github.com/google/uuid"
)

type SongDTO struct {
	ID         uuid.UUID    `json:"id"`
	CreatedAt  time.Time    `json:"createdAt"`
	Duration   string       `json:"duration"`
	Title      string       `json:"title"`
	Genre      string       `json:"genre"`
	SongURL    string       `json:"songUrl"`
	CoverURL   string       `json:"coverUrl"`
	Listenings int64        `json:"listenings"`
	Likes      int64        `json:"likes"`
	Dislikes   int64        `json:"dislikes"`
	User       UserDTO      `json:"user"`
	Authors    []AuthorDTO  `json:"authors"`
	SongTags   []TagDTO     `json:"songTags"`
	Comments   []CommentDTO `json:"comments"`
}

func (b *DTOBuilder) BuildSongDTO(song *models.Song) *SongDTO {
	songAuthors := []AuthorDTO{}
	for _, songAuthor := range song.Authors {
		songAuthors = append(songAuthors, *b.BuildAuthorDTO(&songAuthor))
	}

	songTags := []TagDTO{}
	for _, songTag := range song.SongTags {
		songTags = append(songTags, *b.BuildTagDTO(&songTag))
	}

	comments := []CommentDTO{}
	for _, comment := range song.Comments {
		comments = append(comments, b.BuildCommentDTO(&comment))
	}

	return &SongDTO{
		ID:         song.ID,
		CreatedAt:  song.CreatedAt,
		Duration:   helpers.FormatDuration(song.Duration),
		Title:      song.Title,
		Genre:      song.Genre,
		SongURL:    song.SongURL,
		CoverURL:   song.CoverURL,
		Listenings: song.Listenings,
		Likes:      b.CountSongLikes(song.ID),
		Dislikes:   b.CountSongDislikes(song.ID),
		User:       b.BuildUserDTO(&song.User),
		Authors:    songAuthors,
		SongTags:   songTags,
		Comments:   comments,
	}
}
